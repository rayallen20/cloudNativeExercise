# class2:GO语言编译环境设置

## PART1. 下载GO

[GO安装文件以及源代码](https://golang.google.cn/dl/)

相关环境变量:

- GOROOT:GO的安装目录
- GOPATH:工作目录

	- src:存放源代码
	- pkg:存放依赖包
	- bin:存放可执行文件
- GOOS:指定编译所对应的操作系统
- GOARCH:指定编译所对应的CPU架构
- GOPROXY:代理
	- 国内建议设置:`export GOPROXY=https://goproxy.cn`

## PART2. IDE设置

[练习使用的github地址](https://github.com/cncamp/golang)

## PART3. 基本命令

### 3.1 常用命令介绍

- `go build`:从源代码构建成为可执行文件
- `go env`:打印GO相关的环境变量
- `go fmt`:格式化源代码
- `go get`:下载一个包
- `go install`:编译并安装一个包
- `go mod`:依赖管理
- `go test`:单元测试

### 3.2 常用命令详解

#### 3.2.1 go build

GO语言不支持动态链接,因此编译时会将所有依赖编译进同一个二进制文件.

优点:在OS和CPU架构等条件相同的前提下,因为所有依赖被编译进了同一个二进制文件,所以这个二进制文件被拷贝到任何地方都是可以执行的,是不依赖于第三方包的.

缺点:安装包会稍微大一点.

常用参数:

- `-o`:执行输出目录

	例:`go build -o bin/mybinary`

#### 3.2.2 环境变量GOOS/GOARCH

这两个变量常用于交叉编译.

- GOOS:指定编译所对应的操作系统
- GOARCH:指定编译所对应的CPU架构
	- 例:`GOOS=linux GOARCH=amd64 go build`

GOOS和GOARCH全支持列表:$GOROOT/src/go/build/syslist.go

```go
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

// List of past, present, and future known GOOS and GOARCH values.
// Do not remove from this list, as these are used for go/build filename matching.

const goosList = "aix android darwin dragonfly freebsd hurd illumos ios js linux nacl netbsd openbsd plan9 solaris windows zos "
const goarchList = "386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le loong64 mips mipsle mips64 mips64le mips64p32 mips64p32le ppc riscv riscv64 s390 s390x sparc sparc64 wasm "
```

#### 3.2.3 go test

GO语言原生自带测试.而单元测试往往是面向代码的.

现有一GO工作目录,其结构如下:

```
tree ./
./
├── go.mod
├── go.sum
└── module1
    └── class2
        └── unitTest
            ├── unitTestDemo.go
            └── unitTestDemo_test.go

3 directories, 4 files
```

其中`unitTestDemo.go`为源代码文件,`unitTestDemo_test.go`为该文件的单元测试文件.

源代码文件内容如下:

```go
cat ./module1/class2/unitTest/unitTestDemo.go 
package main

func main() {
        add(1, 2)
}

func add(a, b int) int {
        return a + b
}
```

单元测试文件内容如下:

```go
package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Log("start testing")
	result := add(1, 2)
	assert.Equal(t, result, 3)
}
```

使用`go test`命令进行单测,结果如下:

```
go test -count=1 ./module1/class2/unitTest
ok      cloudNativeExercise/module1/class2/unitTest     0.008s
```

Tip:此处由于已经做过单测了,所以结果被cache了,故使用参数`count=1`关闭缓存.

`go test`命令扫描传入路径下,所有以`_test.go`结尾的文件.惯例是将测试代码与正式代码放在同目录.这样做的好处在于,被测试的文件中如有一些私有函数/变量,测试用例可以覆盖到.

测试代码文件包含`TestXxx`函数.其中`Xxx`可以是任何字母,数字,字符串,但第1个字母不能是小写字母.通常`Test`后接被测试函数的名字.传递给测试函数的参数`t`是`*testing.T`类型.该类型用于管理测试状态并支持格式化测试日志.测试日志会在执行测试的过程中不断累积,并在测试完成时转储至标准输出.[本部分参考文章](https://blog.csdn.net/weixin_39897887/article/details/110256536)

#### 3.2.4 go vet

`go vet`:代码静态检查,发现可能的bug或疑似构造.

vet: v. 仔细检查;彻底审查(某人);诊疗;作兽医

有一些错误,不影响构建或编译,但有可能会有问题,例如:

##### a. 打印时占位符使用不当

现有一段代码如下:

```
package main

import "fmt"

func main() {
	fmt.Printf("%d\n", "hello world")
}
```

静态检查结果如下:

```
 go vet ./module1/class2/vetDemo/vetPrintf/             
# cloudNativeExercise/module1/class2/vetDemo/vetPrintf
module1/class2/vetDemo/vetPrintf/vetPrintf.go:6:2: Printf format %d has arg "hello world" of wrong type string
```

##### b. Boolean错误

对结果一直为true/false的表达式或冗余表达式进行检查

现有一段代码如下:

```go
package main

import "fmt"

func main() {
	var a int = 3
	fmt.Printf("%v\n", a != 3 || a != 2)
}
```

静态检查结果如下:

```
cloudNativeExercise/module1/class2/vetDemo/vetBoolean
module1/class2/vetDemo/vetBoolean/vetBoolean.go:7:21: suspect or: a != 3 || a != 2
```

##### c. 主线程未等待其他goroutine先行退出

现有一段代码如下:

```go
package main

import "fmt"

func main() {
	words := []string{"foo", "bar", "baz"}
	for _, word := range words {
		go func() {
			fmt.Println(word)
		}()
	}
}
```

静态检查结果如下:

```
go vet ./module1/class2/vetDemo/vetWaitGroup 
# cloudNativeExercise/module1/class2/vetDemo/vetWaitGroup
module1/class2/vetDemo/vetWaitGroup/vetWaitGroup.go:9:16: loop variable word captured by func literal
```

##### d. Unreachable的代码

现有一段代码如下:

```go
package main

import "fmt"

func main() {
	Hello()
}

func Hello() {
	fmt.Println("hello")
	return
	fmt.Println("world")
}
```

静态检查结果如下:

```
go vet ./module1/class2/vetDemo/vetUnreachable 
# cloudNativeExercise/module1/class2/vetDemo/vetUnreachable
module1/class2/vetDemo/vetUnreachable/vetUnreachable.go:12:2: unreachable code
```

##### e. error检查滞后

现有一段代码如下:

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://www.spreadsheetdb.io/")
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}
```

静态检查结果如下:

```
go vet ./module1/class2/vetDemo/vetInspectErrorLag 
# cloudNativeExercise/module1/class2/vetDemo/vetInspectErrorLag
module1/class2/vetDemo/vetInspectErrorLag/vetInspectErrorLag.go:10:8: using resp before checking for errors
```

Tip:此处若`err != nil`,`resp`是有可能为`nil`的,所以不能在检查错误之前使用`resp`

## PART4. 代码版本控制

基本操作懂的都懂;高级操作我也不会;故略.

## PART5. Golang Playground

[官方playground](https://go.dev/play/)

[国内可直接访问的playground](https://goplay.tools/)