# class4:包引用与依赖管理 Makefile 项目编译

1. 包引用与依赖管理
2. Makefile
3. 项目编译

## PART1. 包引用与依赖管理

### 1.1 GO语言依赖管理的演进

12 factors中有一条:依赖.

显式声明依赖关系.一些传统语言(如C语言),需要有头文件做声明.但C语言的依赖关系无法声明依赖的文件与具体版本.这样发展下去,项目后期依赖管理会很难做.

#### 1.1.1 回顾GOPATH

GOPATH是通过环境变量设置OS级别的GO语言类库目录

Question:GOPATH的问题?

Answer:

- 不同项目对同一第三方类库,可能依赖的版本不同
- 代码被clone后,需要设置GOPATH才能编译

### 1.1.2 vendor目录

一个项目不可能孤立的存在,多多少少会需要依赖一些第三方类库.在早期,GO项目需要自己去github上把类库git clone下来,然后放到`$GOPATH/pkg/`的对应目录中去.这样才能通过编译.

如果多个项目对同一个第三方类库依赖的版本不同,就没办法了.

vendor目录:自GO1.6之后,支持vendor目录.在每个Go语言项目中,创建一个名叫vendor的目录,并将依赖拷贝至该目录.GO语言项目也会自动将vendor目录作为自身的项目依赖路径.(专门存放依赖的目录.该目录在项目的工作目录下.这样每个项目都有了自己的vendor目录.)

vendor目录的好处:

- 每个项目的vendor目录是独立的,可以灵活的选择版本
- vendor目录与源代码一起check in到 github,其他人 checkout以后可直接编译
- 无需在编译期间下载依赖包,所有依赖都已经与源代码保存在一起

有了vendor目录后,仍旧存在问题:没有依赖声明,仍旧需要手动拷贝类库的代码和文件

How to resolve it?

之后就有了一些依赖管理工具

### 1.2 vendor管理工具

通过声明式配置,实现vendor管理的自动化.

- 早期,GO语言无自带依赖管理工具,社区方案鱼龙混杂,比较出名的有:godep,glide
- 随后,GO语言发布了自带的依赖管理工具gopkg
- 后来,用新的工具gomod替换掉了gopkg
	- 切换mod开启模式:`export GO111MODULE=on/off/auto`(on即通过go mod方式管理;off即通过vendor目录管理;auto则根据GO语言的版本决定使用哪一种)
	- go mod相比之前的工具更灵活易用,已经基本统一了GO语言依赖管理

Question:为什么需要使用依赖管理工具?

1. 版本管理
	
	研发人员看了go.mod,就能够知道这个项目依赖了哪些类库,每个类库具体依赖了哪个版本

2. 提升拉取依赖的效率

	通过依赖管理工具拉取依赖时,并不会拉取依赖类库全部的文件,而是根据你的项目依赖了这个类库中的哪些文件,进而决定拉取依赖类库中的哪些文件.使得项目变的精简.
	
3. 防篡改

	有些疯了的研发,会去改vendor里的代码(可能是因为这样比改项目里的代码工作量要小吧,我没这么干过,怕丢饭碗).但是研发A把某个类库给hack了,研发B是不知道的.这样在离职交接时,接替研发A的人员会直接懵逼.最终结果要么是使用没有被魔改过的依赖类库并修改项目代码;要么是永远的依赖这个被hack过的类库.
	
	依赖工具可以防止出现这种情况.具体的做法是:每次编译时拉取依赖,确保依赖类库的来源一定是offical的.
	
### 1.3 go mod的使用

- step1. 创建项目
- step2. 初始化GO模块
	- `go mod init 模块名称`(初始化后在工作目录下能够看到`go.mod`和`go.sum`这两个文件)
- step3. 在你的代码里引用依赖(当然,此时你的工作目录下并没有vendor/,更没有你依赖的这个类库)
- step4. `go mod tidy`
	- `go mod tidy`的作用:拉取缺少的模块,移除不用的模块.这些模块是被拉取到`$GOPATH/pkg/`下的
- step5. `go mod vendor`(非必要,看需求)
	- `go mod vendor`的作用:将依赖从`$GOPATH/pkg/`复制到vendor目录下
	- 这样做的好处在于:`vendor/`随源码一起提交,别人`git clone`之后,直接就可以编译了;当然,不提交`vendor/`只提交`go.mod`和`go.sum`,让别人`git clone`之后自己`go mod tidy`也是可以的

### 1.4 go.mod sample

本示例出自`k8s.io/apiserver`

```
// 声明本项目的module名称和本项目使用的GO语言版本
module k8s.io/apiserver go 1.13

// 声明本项目直接依赖或间接依赖的模块
require (
	// 间接依赖的模块
	github.com/evanphx/json-patch v4.9.0+incompatible 
	
	// 直接依赖的模块
	github.com/go-openapi/jsonreference v0.19.3 // indirect 
	github.com/go-openapi/spec v0.19.3
	github.com/gogo/protobuf v1.3.2
	golang.org/x/crypto master
	github.com/google/gofuzz v1.1.0
	k8s.io/apimachinery v0.0.0-20210518100737-44f1264f7b6b
)

// 做定制化的部分
replace (
	// replace版本号 当我们清楚地知道本项目依赖的某个第三方类库的具体版本时,可以将这个第三方类库的版本号写入go.mod文件
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	
	// replace第三方类库的源 一般这种replace有2种使用场景:
	// 1. 众所周知的科学上网问题
	// 2. 基于某个社区项目,自己进行了魔改并保存到了自己的github-repo上.但是在你的项目代码中,你不想使用自己的github-repo的路径,仍旧想使用社区的路径.此时就可以把社区的源换成你自己的github-repo的源
	golang.org/x/image => github.com/golang/image
	k8s.io/api => k8s.io/api v0.0.0-20210518101910-53468e23a787
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20210518100737-44f1264f7b6b
	k8s.io/client-go => k8s.io/client-go v0.0.0-20210518104342-fa3acefe68f3
	k8s.io/component-base => k8s.io/component-base v0.0.0-20210518111421-67c12a31a26a
)
```

### 1.5 GOPROXY和GOPRIVATE

- GOPROXY:代理
	- `export GOPROXY=https://goproxy.cn`
	- 在设置`GOPROXY`后,默认所有依赖拉取都需要经过proxy连接git repo拉取代码,并做`checksum`校验

- GOPRIVATE:用于设置让GOPROXY skip的一些仓库地址
	- 使用场景:当你公司里有一些私有的代码仓库时,如果所有的拉取都是通过`GOPROXY`拉取,那代理是连不到你公司内网的代码仓库的.此时就需要`GOPROXY`把这种仓库地址skip掉,让本机直连.
	
使用示例:
```
GOPRIVATE=*.corp.example.com
GOPROXY=proxy.example.com
GONOPROXY=myrepo.corp.example.com
```

## PART2. Makefile

### 2.1 Makefile的用途和语法

Makefile:和`make`命令一起配合使用的.`make`是一个开源的构建工具,很多大型项目的编译都是通过`Makefile`来组织的.

我个人的理解:Makefile用于在编译时,扩展`go build`和`go install`的功能,为编译的过程进行封装,将编译的一些场景(有点类似ORM的Scopes功能)进行模块化.

举例分析Makefile的构成:

```
# 定义变量
export tag=v1.0

# 定义tag 可以认为tag是make命令的参数
root:
	# 定义变量
	export ROOT=github.com/cncamp/golang

build:
	# 定义tag所封装的命令
	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

# 定义tag之间的依赖关系 此处release tag依赖build tag
# 即:调用release之前会先调用build (我没找到更合适的形容方式)
release: build
	echo "building httpserver container"
	docker build -t cncamp/httpserver:${tag} .

push: release
	echo "pushing cncamp/httpserver"
	docker push cncamp/httpserver:v1.0
```

注意:Makefile文件必须用制表符(Tab)做缩进,不能用4个空格代替

[Makefile详解](https://seisman.github.io/how-to-write-makefile/introduction.html)

[适用于GO语言的Makefile教程](https://www.cnblogs.com/FengZeng666/p/15750005.html)

### 2.2 举例

现有一项目,其目录结构如下:

```
 tree ./ -L 1
./
├── Makefile
└── makefileDemo.go

0 directories, 2 files
```

其中`makefileDemo.go`代码如下:

```go
package main

import "fmt"

func main() {
	fmt.Println("hello world")
}
```

`Makefile`文件内容如下:

```
export tag=v1.0

build:
	echo "building makefileDemo binary ${tag}"
	go build -o makefileExec makefileDemo.go

release: build
	echo "building makefileDemo container"
	echo "no docker, fake release"

push: release
	echo "pushing makefileDemo"
	# docker push的命令 但是我没有docker环境 所以省省吧
```

执行`make build`

```
 make build
echo "building makefileDemo binary v1.0"
building makefileDemo binary v1.0
go build -o makefileExec makefileDemo.go
```

查看编译结果:

```
tree ./ -L 1
./
├── Makefile
├── makefileDemo.go
└── makefileExec

0 directories, 3 files
```