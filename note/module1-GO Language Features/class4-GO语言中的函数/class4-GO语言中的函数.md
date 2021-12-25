# class4:GO语言中的函数

## PART1. main函数

每个GO语言程序都应该有个`main package`

`main package`中的`main`函数是GO语言程序入口

```go
package main

import "os"

func main() {
	args := os.Args
	if len(args) != 0 {
		println("Do not accept any argument")
		os.Exit(1)
	}
	println("Hello world")
}
```

## PART2. 参数解析

GO语言的`main`函数与其他语言不同,没有类似`java`的`[]string args`参数.

GO语言通过`os.Args`或`flag`包获取入参

`flag`包会解析入参.

`flag`包使用示例:

```go
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	name := flag.String("name", "world", "specify the name you want to say hi")
	flag.Parse()
	fmt.Println("os args is:", os.Args)
	fmt.Println("input parameter is:", *name)
	fullString := fmt.Sprintf("Hello %s from GO\n", *name)
	fmt.Println(fullString)
}
```

运行结果:

```
 go run flagDemo.go -name fuck
os args is: [/var/folders/9x/kkdcw3dx7js8frvbxh8gns580000gn/T/go-build3664199694/b001/exe/flagDemo -name fuck]
input parameter is: fuck
Hello fuck from GO
```

从运行结果中可以看到,`os.Args`中的第1个元素是二进制文件的绝对路径.第2个元素开始才是用户输入的参数.

## PART3. init函数

作用:在包初始化时运行

`init`函数会在`main`函数之前运行.即使同一个包被作为依赖引用了多次,这个包的`init`函数也只会执行1次.

// TODO:配图解释

```go
package main

import (
	"fmt"
	_ "github.com/cncamp/golang/examples/module1/init/a"
	_ "github.com/cncamp/golang/examples/module1/init/b"
)

func init() {
	fmt.Println("main init")
}

func main() {
	fmt.Println("main func")
}

```

运行结果:

```
go run initDemo.go
init from b
init from a
main init
main func
```

## PART4. 返回值

### 4.1 函数的多返回值

函数可以返回任意数量的返回值.

命名返回值:

- GO的返回值可以在函数签名中被命名,它们会被视作定义在函数顶部的变量
- 返回值的名称应当具有一定的意义,可作为文档使用
- 没有参数的`return`语句返回已命名的返回值.也就是直接返回.

调用者忽略部分返回值:

`result, _ := strconv.Atoi(origStr)`

### 4.2 传递变长参数

- 定义方式

	`func append(slice []Type, elems ...Type) []Type`
	
- 示例:

```go
package main

import "fmt"

func main() {
	b := 1
	c := 2
	d := 3
	sum := accumulate(b, c, d)
	fmt.Println(sum)
}

func accumulate(a ...int) (sum int) {
	for _, v := range a {
		sum += v
	}
	return
}
```

## PART5. 内置函数

// TODO: 写个表格

## PART6. Callback

函数作为参数传入其他函数,并在其它函数内部调用执行

```go
package main

func main() {
	DoOperation(1, increase)
	DoOperation(1, decrease)
}

func increase(a, b int) {
	println("increase result is:", a+b)
}

func decrease(a, b int) {
	println("decrease result is:", a-b)
}

func DoOperation(y int, f func(a, b int)) {
	f(y, 1)
}
```

- Kubernetes controller的leaderelection中有使用

## PART7. Closure

Closure其实就是匿名函数.

- 不能独立存在
- 可以赋值给其他变量
	- `x := func()`
- 可以直接调用
	- `func(x, y int){println(x+y)}(1, 2)`
- 可作为函数的返回值
	- `func Add() (func (b int) int)`

示例:

```go
package main

func main() {
	defer func() {
		if r := recover(); r != nil {
			println("recovered in Func Main")
		}
	}()
}
```

## PART8. Method

### 8.1 定义与使用场景

method:作用在接收者上的函数

`func (recv receiver_type) methodName(parameter_list) (return_value_list)`

没有接收者的话,需要通过`package.func`的方式调用;

有接收者,则是通过`recv.method`调用

- 使用场景

	很多场景下,函数需要的上下文可以保存在receiver属性中,通过定义receiver的方法,该方法可以直接访问receiver属性,减少参数传递

### 8.2 传值还是传指针?

GO语言只有1种规则:值传递

函数内修改参数的值,不会影响函数外原始变量的值

可以传递指针参数将变量地址传递给调用函数,GO语言会复制该指针作为函数内的地址,但指向同一地址

```go
package main

import "fmt"

func main() {
	a := 1
	passByValue(a)
	fmt.Printf("after passByValue, a = %d\n", a)

	fmt.Printf("address of a in main:%p\n", &a)
	passByReference(&a)
	fmt.Printf("after passByReference, a = %d\n", a)
}

func passByValue(a int) {
	a += 1
	fmt.Printf("in passByValue, a = %d\n", a)
}

func passByReference(a *int) {
	fmt.Printf("address of a in func:%p\n", a)
	*a += 1
	fmt.Printf("in passByReference, a = %d\n", *a)
}
```

Question:当我们写代码时,函数的参数传递应该使用struct还是pointer?

Answer:看场景.如果使用pointer,则优点在于不涉及值拷贝,效率可能会高;如果传struct,则当函数退出时,临时变量就回收了,对GC友好.各有利弊,看场景.

## PART9. Interface

### 9.1 定义与使用场景

interface:一组方法的集合

在接口中只能定义方法的组合

```go
type IF interface{
	Method1(param_list) return_type
}
```

使用场景:Kubernetes中有大量的接口抽象和多种实现.比如定义标准接口,让厂商来适配接口即可.

- struct无需显式声明实现interface,只需直接实现方法即可
- struct除实现interface定义的接口外,还可以有额外的方法
- 1个类型可以实现多个接口(GO语言的多重继承)
- GO语言中接口不接受属性定义
- 接口可以嵌套其他接口

示例:

```go
package main

import "fmt"

type IF interface {
	getName() string
}

type Human struct {
	firstName, lastName string
}

func (h *Human) getName() string {
	return h.firstName + "," + h.lastName
}

type Car struct {
	factory, model string
}

func (c *Car) getName() string {
	return c.factory + "-" + c.model
}

func main() {
	interfaces := []IF{}

	h := new(Human)
	h.firstName = "first"
	h.lastName = "last"
	interfaces = append(interfaces, h)

	c := new(Car)
	c.factory = "benz"
	c.model = "z"
	interfaces = append(interfaces, c)

	for _, f := range interfaces {
		fmt.Println(f.getName())
	}
}
```

### 9.2 注意事项

- interface是可能为nil的,所以针对interface的使用一定要先判空,否则会引起程序crash(nil panic)
- struct初始化意味着空间分配,对struct的引用不会出现空指针


## PART10. reflect

### 10.1 定义与使用场景

reflect:让程序自己知道自己的类型或值等属性.

使用场景:Kubernetes中使用reflect来确认用户请求的类型

示例:

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	myMap := make(map[string]string, 10)
	myMap["a"] = "b"
	t := reflect.TypeOf(myMap)
	fmt.Printf("type:%s\n", t)
	v := reflect.ValueOf(myMap)
	fmt.Printf("value:%v\n", v)
}
```

### 10.2 基于struct的反射

```go
package main

import (
	"fmt"
	"reflect"
)

type T struct {
	A string
}

func (t *T) String() string {
	return t.A + "1"
}

func (t T) Change(p *T) string {
	fmt.Printf("p.A = %s\n", p.A)
	t.A = "fuck"
	return t.A
}

func main() {
	myStruct := T{A: "a"}
	v1 := reflect.ValueOf(myStruct)
	for i := 0; i < v1.NumField(); i++ {
		fmt.Printf("Field %d:%v\n", i, v1.Field(i))
	}

	for i := 0; i < v1.NumMethod(); i++ {
		// Tip: 接收者为*T的方法 不属于v1
		fmt.Printf("Method %d:%v\n", i, v1.Method(i))
	}

	// 反射调用传参示例
	param := &T{A: "c"}
	paramCall := reflect.ValueOf(param)
	result := v1.Method(0).Call([]reflect.Value{paramCall})
	fmt.Printf("%v\n", result)
}
```

## PART11. OOP

- 可见性控制
	- public:常量、变量、类型、接口、结构、函数等的名称首字母大写
	- private:非大写则只能在包内使用
- 继承
	- 通过组合实现,内嵌1个或多个struct
- 多态
	- 通过接口实现,通过接口定义方法集,编写多套实现

## PART12. JSON编解码

`json`包使用`map[string]interface{}`和`[]interface{}`类型保存任意对象

解析任意json:

```go
func decodeAny(humanStr string) {
	var obj interface{}
	err := json.Unmarshal([]byte(humanStr), &obj)
	if err != nil {
		return
	}
	objMap := obj.(map[string]interface{})
	for k, v := range objMap {
		switch value := v.(type) {
		case string:
			fmt.Printf("type of %s is string, value is %v\n", k, value)
		case interface{}:
			fmt.Printf("type of %s is interface{}, value is %v\n", k, value)
		default:
			fmt.Printf("type of %s is wrong, value is %v\n", k, value)
		}
	}
}
```