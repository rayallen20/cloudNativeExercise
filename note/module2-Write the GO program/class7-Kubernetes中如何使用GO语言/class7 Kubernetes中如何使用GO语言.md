# class7 Kubernetes中如何使用GO语言

## PART1. Kubernetes中常用代码解读

### 1.1 Rate Limit Queue

#### 老师讲解部分

```go
// ItemExponentialFailureRateLimiter:直白的翻译过来 失败指数级增速限速器 是一个速率限制队列(Rate Limit Queue)
func (r *ItemExponentialFailureRateLimiter) When(item interface{}) time.Duration {
 
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()
    
    exp := r.failures[item]
    r.failures[item] = r.failures[item] + 1
    
	// The backoff is capped such that ‘calculated’ value never overflows.
	backoff := float64(r.baseDelay.Nanoseconds()) * math.Pow(2, float64(exp)) 
	if backoff > math.MaxInt64 {
         return r.maxDelay
    }
    
    calculated := time.Duration(backoff)
    if calculated > r.maxDelay {
         return r.maxDelay
    }
    return calculated
}
```

K8S中有一个模式:控制器模式.这是K8S的核心组件,该组件使得整个集群动态运转起来.实际上该组件是通过监听K8S中对象的变化来实现的.当对象的状态发生变化时,控制器就会对配置做一些管理(这个听起来有点像OOP中的观察者模式?).但是对配置管理,有时候会成功,有时候会失败.

GO语言追求最终一致性,GO语言要求控制器在handle一件事情时,若handle失败,则把该对象重新in queue,以便能够重试(既然这么设计,就说明这次失败不意味着下次还失败,那为什么下次就有可能不失败呢?).

这里就有一个设计上的难点:如果控制器在handle失败时,立刻把该对象in queue,那么马上就会有新的worker thread,又把这个对象获取出来(所以此处的设计是1个对象1个worker thread吗?),worker thread获取到这个刚刚才handle失败的对象,又立刻retry,由于间隔时间过短,handle的结果大概率还是失败,又把这个对象in queue.

综上所述,就会出现一个handle失败的对象,会一直在同一个队列里进进出出,而且进出的频率非常高.最终的结果就是:所有的worker thread都被handle失败的对象占据了.如果在handle的过程中还有call other API的行为,则会给这些API造成极大的压力(1s内可能是千这个数量级的call API).

How to resolve it? -- Circuit breakers

Circuit breakers:在一个系统中,服务提供方(upstream)因访问压力过大而导致响应变慢或失败,服务发起方(downstream)为了保护系统整体的可用性,可以临时暂停对服务提供方的调用,这种牺牲局部保全整体的措施称为熔断(Circuit breakers).

当出现失败时,控制器不直接调用`Add()`操作,而是调用`AddRateLimit()`操作.`AddRateLimit()`操作会调用到`When()`方法.`When()`方法具有指数级backoff(n. 倒扣;补偿)的能力.

如果错误第1次出现,则等待2s后再in queue

如果错误第2次出现,则等待`2^失败次数`(即第2次4s,第3次8s,以此类推)的时长(前边乘的那个`float64(r.baseDelay.Nanoseconds())`不知道是干啥的玩意儿),等待时间上限为5min.过了这个等待时间后再in queue.

这样设计的目的在于:如果对象出现的错误是一个临时性的,短暂的错误,则在retry的早期,由于频繁的retry,所以很快就能被正常的handle;如果这个错误是一个持久性的错误,则随着retry次数的增加,最终到达5min/次的retry频率.使得对短暂性错误和持久性错误的handle有一个平衡.

TODO:`r.failures[item] = r.failures[item] + 1`不写成`r.failures[item] += 1`,是否有什么原因?

**当你要去设计一个controller的failure retry机制时,记得等待,不要给controller周边的system带来流量雪崩.**

#### 我自己猜测的部分

其实这个函数的handle逻辑是统一的,`retry频率 = 2 ^ 失败次数`.

那么由此推测第6行的`exp := r.failures[item]`在第1次失败的时候是等于1的.

因为第1次失败必然返回的是`calculated`;第1次失败时`calculated = 2s`,因此`backoff = 2s`.进而推测在计算`backoff`时,`math.Pow(2, float64(exp)) = 2`,换言之即`2 ^ 1 = 2`.至于那个`float64(r.baseDelay.Nanoseconds())`,可能真的就是表示1e9ns或者1s.

至于`r.maxDelay`,是确定的5min.

那么至此基本可以推断出,`r.failures`是一个记录失败次数的集合.我倾向于这个集合是一个类似于map的数据结构,`item`表示控制器要处理的对象,`r.failures[item]`表示这个对象被handle失败的次数.那么`r.failures`既然是记录对象与handle失败次数映射关系的集合,进而可以推测:这个`r`大概率就是用来限制速率的队列.从命名上看,r也很有可能是Rate Limit Queue的缩写.

## PART2. Kubernetes日常运维中的代码调试场景

### 2.1 案例1:空指针

#### 2.1.1 问题描述

Kuberbetes调度器在调度有外挂存储需求的pod的时候,在获取节点信息失败时,会异常退出.

```
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x105e2 83]
```

#### 2.1.2 根因分析

`nil pointer`是GO语言中最常出现的一类错误,也最容易判断,这类错误表示对内存的非法访问.通常在call stack中就会告诉你哪行代码出问题了.

在调度器`csi.go`中有如下代码:

```go
node := nodeInfo.Node()
if node == nil {
   return framework.NewStatus(framework.Error,
fmt.Sprintf("node not found: %s", node.Name))
}
```

当`node`为`nil`时,对`node`的应用`node.Name`就会引发空指针.

#### 2.1.3 解决方案

[PR:当指针为空时,不要继续引用](https://github.com/kubernetes/kubernetes/pull/102229)

### 2.2 Map的读写冲突

#### 2.2.1 问题描述

程序在遍历Kubernetes对象的`Annotation`时异常退出

#### 2.2.2 根因分析

Kubernetes对象中`Label`和`Annotation`的数据类型是`map[string]string`,经常有代码需要修改这2个map.但修改map的同时可能会有其他线程`for range`遍历.

#### 2.2.3 解决方案

- 使用`sync.RWMutex`加读写锁
- 使用线程安全的map,比如`sync.Map{}`

### 2.3 kube-proxy消耗10个CPU

#### 2.3.1 问题描述

客户汇报问题:kube-proxy消耗了主机10个CPU.

#### 2.3.2 根因分析

- 登录问题节点,执行`top`命令查看cpu消耗,可以看到kube-proxy的cpu消耗和pid信息
- 对kube-proxy进程运行system profiling tool,发现10个CPU中,超过60%的CPU都在做GC,这说明GC需要回收的对象太多了,说明程序创建大量可回收对象.

排查过程:

- step1. 查看kube-proxy的CPU开销

```
perf top –p <pid>
26.48% kube-proxy [.] runtime.gcDrain 
13.86% kube-proxy [.] runtime.greyobject 
10.71% kube-proxy [.] runtime.(*lfstack).pop 10.04% kube-proxy [.] runtime.scanobject
```

就看`runtime`后边那几个词,那大概率就是占着CPU在做GC了.

为什么做GC? 那只能是有太多的对象需要回收(经典废话文学).

- step2. 通过pprof查看堆内存占用情况(因为怀疑的是GC,所以要去查看堆内存)

```
curl 127.0.0.1:10249/debug/pprof/heap?debug=2
1: 245760 [301102: 73998827520] @ 0x11ddcda 0x11f306e 0x11f35f5 0x11fbdce 0x1204a8a 0x114ed76 0x114eacb 0x11
# 0x11ddcd9 
k8s.io/kubernetes/vendor/github.com/vishvananda/netlink.(*Handle).RouteListFiltered+0x679
# 0x11f306d k8s.io/kubernetes/pkg/proxy/ipvs.(*netlinkHandle).GetLocalAddresses+0xed
# 0x11f35f4 k8s.io/kubernetes/pkg/proxy/ipvs.(*realIPGetter).NodeIPs+0x64
# 0x11fbdcd k8s.io/kubernetes/pkg/proxy/ipvs.(*Proxier).syncProxyRules+0x47dd
```

301102个对象,占用了73998827520的内存.

谁创建的这些对象? 继续往后看这些对象的调用栈,可以发现,是`k8s.io/kubernetes/pkg/proxy/ipvs`这个文件创建的.

我们根据方法名来猜测一下,这个文件都干了什么事:

`syncProxyRules`:配置负载均衡规则

`NodeIPs`:获取Node的IP

`GetLocalAddresses`:获取本地地址(老实说是这一步时创造了30多W个对象,不知道他是咋确定在这一步才创建大量对象的,PPT中标红的部分,是pprof)

`RouteListFiltered`:路由列表过滤

根据我对[源码](https://github.com/kubernetes/kubernetes/blob/master/pkg/proxy/ipvs/netlink_linux.go#L143)的分析:

```go
// GetLocalAddresses return all local addresses for an interface.
// Only the addresses of the current family are returned.
// IPv6 link-local and loopback addresses are excluded.
func (h *netlinkHandle) GetLocalAddresses(dev string) (sets.String, error) {
	ifi, err := net.InterfaceByName(dev)
	if err != nil {
		return nil, fmt.Errorf("Could not get interface %s: %v", dev, err)
	}
	addr, err := ifi.Addrs()
	if err != nil {
		return nil, fmt.Errorf("Can't get addresses from %s: %v", ifi.Name, err)
	}
	return utilproxy.AddressSet(h.isValidForSet, addr), nil
}
```

这个`addr`可能就是"大量对象"(注意我说的是可能,**Maybe**).

- step3. heap dump分析
	- `GetLocalAddresses()`函数调用创建了301102个对象,占用内存73998827520
	- 如此多的对象被创建,显然会导致kube-proxy进程忙于 GC,占用大量CPU
	- 对照代码分析`GetLocalAddresses()`的实现,发现该函数的主要目的是获取节点本机IP地址,获取的方法时通过`ip route`命令获取当前节点所有local路由信息并转换成go struct,然后过滤掉ipvs0网口上的路由信息
	- `ip route show table local type local proto kernel`
	- 因为集群规模较大,该命令返回5000条左右记录,因此每次函数调用都会有数万个对象被生成
	- 而kube-proxy在处理每一个服务的时候都会调用该方法,因为集群有数千个服务,因此,kube- proxy在反复调用该函数创建大量临时对象

Question:如果让你去设计一个程序,该程序用于获取一个节点的IP,你会怎么做?

Answer:如果是你自己公司里,你熟悉的集群,可能你知道所有的节点的网卡都叫`eth0`或者`ens33`,那你可以直接获取指定网卡的IP.但K8S要做的是一个通用的解决方案,通用方案就不能对device name做假设.

K8S实际上是通过`ip route show table local`命令查本地路由表的信息

```
ip route show table local
broadcast 127.0.0.0 dev lo proto kernel scope link src 127.0.0.1 
local 127.0.0.0/8 dev lo proto kernel scope host src 127.0.0.1 
local 127.0.0.1 dev lo proto kernel scope host src 127.0.0.1 
broadcast 127.255.255.255 dev lo proto kernel scope link src 127.0.0.1 
broadcast 内网IP dev eth0 proto kernel scope link src 内网IP
local 内网IP dev eth0 proto kernel scope host src 内网IP 
broadcast 内网IP dev eth0 proto kernel scope link src 内网IP
```

```
[root@wdm ~]# ip route show table local type local
local 127.0.0.0/8 dev lo proto kernel scope host src 127.0.0.1 
local 127.0.0.1 dev lo proto kernel scope host src 127.0.0.1 
local 内网IP dev eth0 proto kernel scope host src 内网IP 
```

```
ip route show table local type local proto kernel
local 127.0.0.0/8 dev lo scope host src 127.0.0.1 
local 127.0.0.1 dev lo scope host src 127.0.0.1 
local 内网IP dev eth0 scope host src 内网IP 
```

本地路由表里有5000条记录.

那为什么1个节点上能有这么多记录呢?这就和kube-proxy的实现有关了.kube-proxy在处理这些服务时,给每一个服务设置了一个虚IP.kube-proxy会把这些虚IP绑定在本地一个叫`ipvs0`的单位device上.当你执行`ip address`或`ip route`时,看到的IP数量就非常多了.这5000条的IP列表就是这么来的.

然后kube-proxy基于这些IP,把它们都反序列化为GO的struct.很明显,IP都是虚IP,那基于这些虚IP反序列化出来的struct,也必然是临时对象.这个根据需IP反序列化生成struct的逻辑,是写在一个for循环里的,每分钟执行1次.可以认为每1分钟创建5000个临时对象.使得节点的内存迅速耗尽,GC也需要疯狂输出以便尽快回收内存.但GC的速度始终是赶不上碎片对象产生的速度,因此CPU忙于GC,一直消耗CPU.

#### 2.3.3 解决方案

- 把对`GetLocalAddresses`函数的调用写到循环外边去.
- [PR地址](https://github.com/kubernetes/kubernetes/pull/79444)

### 2.4 线程池耗尽

#### 2.4.1 问题描述

在K8S中有一个控制器,叫做endpoint controller,该控制器符合生产者消费者模式,默认有5个worker线程作为消费者.该消费者在处理请求时,可能调用LBaaS的API更新负载均衡配置.我们发现该控制器会时不时的不工作.其具体表现为:该做的配置变更没有发生,相关的日志也没有打印.

#### 2.4.2 根因分析

通过pprof打印出该进程所有goroutine信息,发现worker线程都卡在http请求调用处.

当worker线程调用LBaaS API时,底层是net/http包的调用,即网络通信.而客户端在发起连接请求时,没有设置超时时间.这就导致当出现某些网络异常时,客户端会永远处于等待状态(说白了就是hang在连接上了).

#### 2.4.3 解决方案

修改代码加入客户端超时控制