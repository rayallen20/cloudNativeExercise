# class5:GO语言常用语法 多线程

## PART1. 常用语法

### 1.1 错误处理

GO语言无内置exception机制,只提供error接口供使用者定义错误.实际上这是一种简化.

创建一个error:

1. `errors.New()`
2. `fmt.Errorf()`

通常应用程序对error的处理大部分是判断error是否为nil.

kubernetes自定义了与apiserver交互的不同类型错误.

```go
type StatusError struct{
	ErrStatus metav1.Status
}

func (e *StatusError) Error() string {
	return e.ErrStatus.Message
}
```

### 1.2 defer

作用:**函数返回之前**执行某个语句或函数

使用场景:关闭你打开的资源

```go
defer file.Close()
defer mu.Unlock()
defer println("")
```

```go
package main

import "fmt"

func main() {
	defer fmt.Printf("1\n")
	defer fmt.Printf("2\n")
	defer fmt.Printf("3\n")
}
```

运行结果:

```
go run deferDemo.go
3
2
1
```

defer的本质是将指令**压栈**,所以调用顺序和代码的书写顺序相反,先进后出.

Question:以下代码出错的原因是什么?

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	defer fmt.Printf("1\n")
	defer fmt.Printf("2\n")
	defer fmt.Printf("3\n")
	loopFunc()
	time.Sleep(1 * time.Second)
}

func loopFunc() {
	lock := sync.Mutex{}
	for i := 0; i < 3; i++ {
		lock.Lock()
		defer lock.Unlock()
		fmt.Println("loopFunc:", i)
	}
}
```

Answer:第1次循环时加锁,第2次循环还没有解锁,所以再加锁就死锁了.defer是函数退出时执行,而非defer所在的代码块结束时执行.

Question:同样的逻辑,如何修改才能避免死锁?

Answer:Closure

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	defer fmt.Printf("1\n")
	defer fmt.Printf("2\n")
	defer fmt.Printf("3\n")
	loopFunc()
	time.Sleep(1 * time.Second)
}

func loopFunc() {
	lock := sync.Mutex{}
	for i := 0; i < 3; i++ {
		go func(i int) {
			lock.Lock()
			defer lock.Unlock()
			fmt.Println("loopFunc:", i)
		}(i)
	}
}
```

### 1.3 panic和recover

panic:当系统出现不可恢复的错误时,主动调用内置函数`panic()`,使当前线程直接crash

defer:保证执行并把控制权交还给接收到panic的函数调用者

recover:内置函数`recover()`用于从`panic`或错误场景中恢复,可以使线程不被crash

Tip:即使线程受`panic`而crash了,`defer`仍旧能够执行.因为`defer`的本质是将指令压栈(换言之,在`panic`之前`defer`的指令就已经压栈了).

```go
package main

import "fmt"

func main() {
	defer func() {
		fmt.Println("defer func is called")
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	panic("a panic is triggered")
}
```

运行结果:

```
go run recoverDemo.go 
defer func is called
a panic is triggered
```

## PART2. 多线程

### 2.1 并发与并行

并发(concurrency):两个或多个事件在同一**时间间隔**发生

并行(parallellism):两个或多个事件在同一**时刻**发生

[并发概述(《GO语言并发之道》 第1章)](https://github.com/rayallen20/ConcurrencyInGo/blob/master/note/%E7%AC%AC1%E7%AB%A0%20%E5%B9%B6%E5%8F%91%E6%A6%82%E8%BF%B0/%E7%AC%AC1%E7%AB%A0%20%E5%B9%B6%E5%8F%91%E6%A6%82%E8%BF%B0.md)

### 2.2 协程

#### 2.2.1 进程

- 系统分配资源(CPU时间片、内存等)的基本单位
- 有独立内存空间,切换开销大.进程都有自己的虚拟内存空间,有可能多个进程的虚拟内存空间最终映射到相同的一个物理内存空间地址上.

#### 2.2.2 线程

- 线程:进程的一个执行流,这个执行流和进程共享同样的内存地址空间.是CPU调度并能够独立运行的基本单位
- 同一进程中的多线程共享内存空间,线程切换代价小(切换线程时不需要切换进程空间)
- 多线程通信方便(地址空间相同,所以所有线程可以访问同样的内存地址,多线程通信就可以通过共享内存实现)
- 从内核层面来看线程其实也是一种特殊的进程,它跟父进程共享了打开的文件和文件系统信息,共享了地址空间和信号处理函数

#### 2.2.3 协程

- 协程:GO语言中的轻量级线程实现
- golang在runtime、sysCall等多方面对goroutine调度进行了封装和处理,当遇到长时间执行或进行系统调用时,会主动把当前goroutine的CPU(P)转让出去,让其他goroutine能够被调度并执行.也就是说,golang从语言层面上支持了协程

#### 2.2.4 线程和协程的差异

- 每个goroutine默认占用内存远远小于Java、C的线程
	- goroutine:2KB
	- 线程:8MB

- goroutine切换开销方面,远远小于线程
	- 线程:设计模式切换(从用户态切换到内核态)、16个寄存器(PC、SP、AX等)的刷新
	- goroutine:3个寄存器的值修改:PC、SP、DX

- `runtime.GOMAXPROCS()`:控制并行线程的数量

[进程、线程、协程](https://github.com/rayallen20/ConcurrencyInGo/blob/master/extend/%E8%BF%9B%E7%A8%8B%E3%80%81%E7%BA%BF%E7%A8%8B%E3%80%81%E5%8D%8F%E7%A8%8B.md)

### 2.3 CSP

#### 2.3.1 定义

CSP:Communicating Sequential Process.通信顺序进程.描述两个独立的并发实体通过共享的通讯channel进行通信的并发模型.**不要通过共享内存进行通信.相反,通过通信来共享内存.**

#### 2.3.2 goroutine

goroutine:是一种轻量线程,它不是OS的线程,而是将一个OS线程分段使用,通过调度器实现协作式调度.

goroutine是一种绿色线程,微线程.它与Coroutine协程也有区别,能够在发现阻塞后启动新的微线程.

#### 2.3.3 channel

channel:类似Unix的Pipe,用于协程之间通讯和同步.协程之间虽然解耦,但它们和channel有耦合(需要依靠channel来传递数据,协同步调)

[对你的代码建模-通信顺序进程(《GO语言并发之道》 第2章)](https://github.com/rayallen20/ConcurrencyInGo/blob/master/note/%E7%AC%AC2%E7%AB%A0%20%E5%AF%B9%E4%BD%A0%E7%9A%84%E4%BB%A3%E7%A0%81%E5%BB%BA%E6%A8%A1-%E9%80%9A%E4%BF%A1%E9%A1%BA%E5%BA%8F%E8%BF%9B%E7%A8%8B/%E7%AC%AC2%E7%AB%A0%20%E5%AF%B9%E4%BD%A0%E7%9A%84%E4%BB%A3%E7%A0%81%E5%BB%BA%E6%A8%A1-%E9%80%9A%E4%BF%A1%E9%A1%BA%E5%BA%8F%E8%BF%9B%E7%A8%8B.md)

### 2.4 协程示例

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		go fmt.Println(i)
	}
	time.Sleep(time.Second)
}
```

### 2.5 使用channel实现多线程通信

#### 2.5.1 channel的定义与使用场景

channel:多个协程之间通讯的管道

- 一端发送数据,一端接收数据(**不要通过共享内存进行通信.相反,通过通信来共享内存.**)
- 同一时间只有1个协程可以访问数据,不会出现共享内存模式的内存竞争问题(当然就更不涉及锁、临界区等概念了)
- 协调协程的执行顺序("内存访问同步"中的同步并不是"同时进行"的意思,而是"协同步调"--协调指令运行的先后顺序)

#### 2.5.2 声明方式与操作符

声明方式:`var identifier chan datatype`

操作符:`<-`

#### 2.5.3 示例

```go
package main

import "fmt"

func main() {
	ch := make(chan int)
	go func() {
		fmt.Println("hello from goroutine")
		ch <- 0
	}()
	<-ch
}
```

#### 2.5.4 通道缓冲

- 基于channel的通信是同步的(线程安全)
- 当缓冲区满时,数据的发送是阻塞的
- 通过`make`关键字创建通道时可定义缓冲区容量,默认缓冲区容量为0

channel缓冲区长度为0,就意味着发送方和接收方都就绪,发送方才能把数据送入管道;否则发送方阻塞.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ch := make(chan int, 10)
	go func() {
		for i := 0; i < 10; i++ {
			// generate a random int n and n will be between 0 and 10
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(10)
			fmt.Println("putting:", n)
			ch <- n
		}
		close(ch)
	}()

	fmt.Println("hello from main")
	for v := range ch {
		fmt.Println("receiving:", v)
	}
}
```

#### 2.5.5 单向通道

只写通道:`var sendOnly chan<- int` 

只读通道:`var readOnly <-chan int`

单向通道基本上不能独立存在.

使用场景:双向通道的转换

```go
package main

func main() {
	c := make(chan int)
	go prod(c)
	go consume(c)
}

func prod(ch chan<- int) {
	for {
		ch <- 1
	}
}

func consume(ch <-chan int) {
	for {
		<-ch
	}
}
```

#### 2.5.6 关闭通道

- 通道无需每次关闭
- 关闭的作用:告知接收者,该通道再无新数据发送
- 只有发送方需要关闭通道

```go
package main

import "fmt"

func main() {
	ch := make(chan int)
	defer close(ch)
	if v, notClosed := <- ch; notClosed {
		fmt.Println(v)
	}
}
```

#### 2.5.7 select

作用:当多个协程同时运行时,可通过关键字`select`轮询多个通道.

- 若所有通道都阻塞时,如定义了`default`则执行`default`;如没定义`default`则等待
- 如多个通道就绪则随机选择

#### 2.5.8 time.Ticker

#### a. 定义

```go
// A Ticker holds a channel that delivers ``ticks'' of a clock
// at intervals.
type Ticker struct {
	C <-chan Time // The channel on which the ticks are delivered.
	r runtimeTimer
}
```

#### b. 使用场景

为协程设定超时时间.

机制:当到达设定的超时时间时,向`C`中发送一个`Time`对象.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	timer := time.NewTicker(time.Second)
	ch := make(chan int)
	select {
	case <-ch:
		fmt.Println("received from ch")
	case <-timer.C:
		fmt.Println("timeout waiting from channel ch")
	}
}
```

#### 2.5.9 使用channel停止子协程

机制:关闭管道时,管道的消费者一方会收到消息

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("done channel is triggered, exit child goroutine")
				return
			default:
				fmt.Println("handle biz")
			}
		}
	}()

	time.Sleep(time.Millisecond)
	close(done)
}
```

#### 2.5.10 使用channel模拟生产者/消费者模型

```go
package main

import (
	"fmt"
	"time"
)

// producer与consumer(worker)通信用
var messages = make(chan int, 10)

// 用于标识worker是否退出用
var done = make(chan bool)

func main() {
	go consumer()
	producer()
	defer close(messages)

	time.Sleep(5 * time.Second)
	close(done)
	time.Sleep(time.Second)
	fmt.Println("main process exit!")
}

func consumer() {
	ticker := time.NewTicker(time.Second)
	// 每秒消费1个
	for _ = range ticker.C {
		select {
		// close(done)会触发该分支
		case <-done:
			fmt.Println("consumer process interrupt...")
			return
		default:
			// Tip:此处原文写的是fmt.Printf("send message: %d\n", <-messages)
			// 个人推测可能写错了?
			fmt.Printf("received message:%d\n", <-messages)
		}
	}
}

func producer() {
	for i := 0; i < 10; i++ {
		messages <- i
	}
}
```

### 2.6 Context

#### 2.6.1 定义

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}
```

#### 2.6.2 作用

- 超时控制、取消操作或一些异常情况,通常是需要抢占操作或中断后续操作时使用
- `Context`可以设置截止日期(超时时长)、同步信号、传参

#### 2.6.3 具体用法

- `context.Background`
	- `Background`通常被用于`main`函数、初始化以及测试中,作为顶层`context`出现.也就是说一般我们创建的`context`都是基于`Background`

- `context.TODO`
	- 不确定上下文时使用`TODO`

- `context.WithDeadline()`
	- 设置超时时间.和`time.Ticker`机制相同,到了超时时间会向管道中发送一个空结构体,标识到达超时时间,供调用方进行超时控制

- `context.WithValue`
	- 向`context`添加键值对

- `context.WithCancel`
	- 创建一个可取消的`context`

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	baseCtx := context.Background()
	// 向上下文中写入一个键值对
	ctx := context.WithValue(baseCtx, "key", "value")
	go func(c context.Context) {
		// 被主线程创建的goroutine 通过上下文读取传递的变量
		fmt.Println(c.Value("key"))
	}(ctx)

	// 超时控制
	timeoutCtx, cancel := context.WithTimeout(baseCtx, time.Second)
	defer cancel()
	go func(ctx context.Context) {
		// 定时器 每秒向ticker.C中放入一个Time对象
		ticker := time.NewTicker(time.Second)
		// 每秒循环1次
		for _ = range ticker.C {
			select {
			case <-ctx.Done():
				// 超时处理
				fmt.Println("child process interrupt...")
				return
			default:
				// 正常逻辑
				fmt.Println("enter default")
			}
		}
	}(timeoutCtx)

	select {
	// 等待子线程完成
	case <-timeoutCtx.Done():
		time.Sleep(time.Second)
		fmt.Println("main process exit!")
	}
}
```

#### 2.6.4 基于context停止goroutine

使用`context.WithTimeout()`即可

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("child process interrupt...")
				return
			default:
				fmt.Println("handle biz")
			}
		}
	}(timeoutCtx)

	time.Sleep(time.Microsecond)
	fmt.Println("err:", ctx.Err())
}
```

## PART3. 课后练习

要求:基于`chan`编写一个简单的单线程生产者消费者模型

1. channel长度为10,元素类型为int
2. 生产者每1s向channel中放入1个类型为int的元素,当队列满时可阻塞
3. 消费者每1s从channel中获取1个元素并打印,channel为空时消费者阻塞

```go
package main

import (
	"fmt"
	"time"
)

var message = make(chan int, 10)

func main() {
	go consumer()
	producer()
}

func consumer() {
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		select {
		case v := <-message:
			fmt.Println(v)
		}
	}
}

func producer() {
	for i := 0; i < 10; i++ {
		message <- i
		time.Sleep(time.Second)
	}
}
```