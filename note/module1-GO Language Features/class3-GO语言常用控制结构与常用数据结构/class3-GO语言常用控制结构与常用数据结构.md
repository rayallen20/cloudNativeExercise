# class3:GO语言常用控制结构与常用数据结构

## PART1. 控制结构

### 1.1 分支结构

#### 1.1.1 if

##### a. 基本形式

```
if condition1 {
	// do sth
} else if condition2 {
	// do sth else
} else {
	// catch-all or default
}
```

##### b. 简短语句

和`for`一样,`if`与可以在条件表达式前执行一个简单语句.

```
if v := x - 100; v < 0 {
	return v
}
```

#### 1.1.2 switch

用途:分支较多时使用,相比于`if`,代码表现力更强.

```
switch var1 {
	case val1:
		// nil branch
	case val2:
		// 执行下一个case中的代码块
		fallthrough
	case val3:
		f()
	default:
		// default branch
}
```

Tip:关键字`fallthrough`表示执行下一个分支的代码块.

### 1.2 循环结构

#### 1.2.1 for

GO只有1种循环结构:`for`循环

##### a. 计入计数器的循环

```
for 初始化语句; 条件语句; 修饰语句 {}
```

```
for i := 0; i < 10; i++ {
	sum += i
}
```

##### b. 与while等价的for循环

GO语言不支持`while`.但`for`循环中,初始化语句和修饰语句是可选的,这种场景下与while等价.

```
for sum < 1000 {
	sum += sum
}
```

##### c. 无限循环

```
for {
	if condition1 {
		break
	}
}
```

#### 1.2.2 for range

作用:遍历数组,切片,字符串,map等集合型的数据结构.

```
for index, char := range myString {

}
```

```
for key, value := range myMap {

}
```

```
for key, value := range myArray {

}
```

**注意:如果for range遍历指针数组,则value取出的指针地址为原地址的拷贝**

[原因详解](https://github.com/rayallen20/cloudNativeExercise/blob/master/note/module1-GO%20Language%20Features/class3-GO%E8%AF%AD%E8%A8%80%E5%B8%B8%E7%94%A8%E6%8E%A7%E5%88%B6%E7%BB%93%E6%9E%84%E4%B8%8E%E5%B8%B8%E7%94%A8%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84/for%20range%E8%AF%A6%E8%A7%A3.md)

```go
package main

import "fmt"

const MAX = 3

func main() {
	ints := []int{10, 100, 200}
	// 指针数组
	var ptrs [MAX]*int

	// 对指针数组赋值
	for i := 0; i < MAX; i++ {
		ptrs[i] = &ints[i]
	}

	// 修改指针数组
	for _, v := range ptrs {
		*v += 100
	}

	// 在for range中对元素的修改 会影响到数组
	for _, v := range ptrs {
		fmt.Printf("%d\n", *v)
	}
}
```

### 1.3 示例代码

```go
package main

import "fmt"

func main() {
	for i := 0; i < 3; i++ {
		fmt.Println(i)
	}

	fullString := "hello world"
	fmt.Println(fullString)
	for i, c := range fullString {
		// c的数据类型为rune 需转换为string 打印时才是这个字符的字面量
		// 否则打印的是ASCII码
		fmt.Println(i, string(c))
	}
}
```

## PART2. 数据结构

### 2.1 变量与常量

#### 2.1.1 常量

常量:不可修改

`const identifier type`

#### 2.1.2 变量

变量:可修改

`var identifier type`

Tip:GO语言不允许出现不被使用的变量

##### a. 关键字var

`var`语句用于声明一个**变量列表**,和函数的参数列表一样,类型在最后

例:

`var c, python, java bool`

##### b. 变量的初始化

- 变量声明可以包含初始值,每个变量对应一个
- 如果初始化值已存在,则可省略类型;变量会从初始值中获得类型
- 初始化同时赋值多个变量:`var i, j int = 1, 2`

##### c. 短变量声明

如果变量的值在初始化时是清晰且确定的,那么可以使用短变量声明的方式.这种场景下变量的类型声明可以省略,因为所赋的值的类型已经决定了这个变量的类型.

- 在函数中,简洁赋值语句`:=`可在类型明确的场景下代替`var`声明
- 函数外的每个语句都必须以关键字开始(`var`,`const`,`func`等),因此`:=`结构不能在函数外使用
- 例:`c, python, java := true, false, "no"`

### 2.2 类型转换与推导

#### 2.2.1 类型转换

表达式`T(v)`将值`v`转换为类型`T`.GO语言不支持隐式转换,只能使用显式转换.

- 一些关于数值的转换:

```go
var i int = 42
var f float64 = float64(i)
var u uint = uint(f)
```

- 或者使用更简单的形式:

```go
i := 42
f := float64(i)
u := uint(f)
```

#### 2.2.2 类型推导

在声明一个变量而不指定其类型时(即使用不带类型的`:=`语法或`var = `表达式语法),变量的类型由右值推导得出.

```go
var i int
j := i		// j也是一个int
```

### 2.3 数组

- 相同类型且**长度固定**、连续的内存片段
- 以索引访问每个元素

定义方法:

`var identifier [len]type`

例:

```go
myArray := [3]int{1,2,3}
```

### 2.4 切片(slice)

**切片是对数组一个连续片段的引用.可以认为切片是数组中某一个连续的子集.**

- 数组定义中不指定长度即为切片:`var identifier []type`
- 切片在未初始化之前默认为nil,长度为0

```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	myArray := [5]int{1, 2, 3, 4, 5}
	mySlice := myArray[1:3] // 取左闭右开区间
	fmt.Printf("%+v\n", mySlice)
	fullSlice := myArray[:] // 取数组全部元素
	remove3rdItem, _ := deleteItem(fullSlice, 2)
	fmt.Printf("remove3rdItem %+v\n", remove3rdItem)
}

// deleteItem 根据索引从数组中删除1个元素
func deleteItem(slice []int, index int) ([]int, error) {
	if index >= len(slice) {
		return nil, errors.New("illegal index")
	}
	return append(slice[:index], slice[index+1:]...), nil
}
```

### 2.5 内置函数make和new

`new`:返回切片的指针地址

`make`:返回第1个元素,可预设内存空间,避免未来的内存拷贝

```go
package main

import "fmt"

func main() {
	// 该指针中的地址 指向的切片长度和容量均为0
	mySlice1 := new([]int)
	fmt.Printf("len:%d, cap:%d\n", len(*mySlice1), cap(*mySlice1))

	// 以下是比较奇怪的用法 但是结果正确
	*mySlice1 = append(*mySlice1, 1)
	fmt.Printf("%#v\n", mySlice1)
	(*mySlice1)[0] = 2
	fmt.Printf("%#v\n", mySlice1)

	mySlice2 := make([]int, 0)
	mySlice3 := make([]int, 10)
	mySlice4 := make([]int, 10, 20)

	fmt.Printf("mySlice1: %p\n", mySlice1)
	fmt.Printf("mySlice2: len:%d cap:%d\n", len(mySlice2), cap(mySlice2))
	fmt.Printf("mySlice3: len:%d cap:%d\n", len(mySlice3), cap(mySlice3))
	fmt.Printf("mySlice4: len:%d cap:%d\n", len(mySlice4), cap(mySlice4))
}
```

### 2.6 切片的常见问题

#### 2.6.1 扩容导致的内存地址变化问题

```go
package main

import "fmt"

func main() {
	a := []int{}
	b := []int{1, 2, 3}
	c := a
	// a受扩容影响 内存地址发生变化 但c仍旧指向变化前的地址
	a = append(b, 1)
	fmt.Printf("address of a:%p\n", a)
	fmt.Printf("address of c:%p\n", c)
}
```

运行结果:

```
go run sliceProblem1.go 
address of a:0xc00001e120
address of c:0x116be80
```

#### 2.6.2 范围遍历修改切片中元素的值

范围遍历时,`k ,v := range mySlice`中,`v`是一个临时变量而非切片中的元素(二者仅仅是值相同),对`v`的修改不会影响到切片.

```go
package main

import "fmt"

func main() {
	mySlice := []int{10, 20, 30, 40, 50}
	for key, value := range mySlice {
		fmt.Printf("address of value: %p\n", &value)
		fmt.Printf("address of mySlice[%d]: %p\n", key, &mySlice[key])
	}
}
```

运行结果:

```
go run sliceProblem2.go 
address of value: 0xc0000b2008
address of mySlice[0]: 0xc0000aa030
address of value: 0xc0000b2008
address of mySlice[1]: 0xc0000aa038
address of value: 0xc0000b2008
address of mySlice[2]: 0xc0000aa040
address of value: 0xc0000b2008
address of mySlice[3]: 0xc0000aa048
address of value: 0xc0000b2008
address of mySlice[4]: 0xc0000aa050
```

`value`的地址相同这个情况,也可以侧面说明范围遍历指针数组时,`value`取出的指针地址是原指针地址的拷贝这个现象.

### 2.7 Map

#### 2.7.1 声明方法

示例:

```go
package main

import "fmt"

func main() {
	myMap := make(map[string]string, 10)
	myMap["a"] = "b"
	myFuncMap := map[string]func() int{
		"funcA": func() int { return 1 },
	}
	fmt.Println(myFuncMap)
	f := myFuncMap["funcA"]
	fmt.Println(f())
}
```

运行结果:

```
go run mapDemo.go 
map[funcA:0x108a500]
1
```

#### 2.7.2 访问Map元素

```go
package main

func main() {
	myMap := map[string]string{
		"a": "a1",
		"c": "c3",
		"d": "d4",
		"e": "e5",
	}
	element, exist := myMap["b"]
	if exist {
		println(element)
	}

	// 遍历Map
	for key, value := range myMap {
		println(key, value)
	}
}
```

运行结果:

```
go run elementAccessFromMap.go
e e5
a a1
c c3
d d4
```

### 2.8 结构体和指针

#### 2.8.1 结构体

用于定义对象

```go
package main

type MyType struct {
	Name string
}

func printMyType(t *MyType) {
	println(t.Name)
}

func main() {
	t := MyType{Name: "test"}
	printMyType(&t)
}
```

#### 2.8.2 指针

GO语言支持指针,但不支持指针运算

- 指针变量的值是内存地址
- 未赋值的指针**值**为`nil`

### 2.9 结构体标签

结构体中的字段,除了字段名和类型之外,还有一个可选的标签

使用场景:Kubernetes APIServer对所有资源的定义都用Json tag和protoBuff tag.

```go
package main

import (
	"fmt"
	"reflect"
)

type MyType struct {
	Name string `json:"name"`
}

func main() {
	mt := MyType{Name: "test"}
	myType := reflect.TypeOf(mt)
	name := myType.Field(0)
	tag := name.Tag.Get("json")
	fmt.Println(tag)
}
```

### 2.10 类型重命名

作用:对某个另外的类型进行重命名.

使用场景:可以当做枚举使用

```go
type ServiceType string

const (
	ServiceTypeA ServiceType = "serviceA"
	ServiceTypeB ServiceType = "serviceB"
)
```

## PART3. 课后练习

### 3.1 题目要求

给定一个字符串数组`[]string{"I", "am", "stupid", "and", "week"}`,将其修改为`[]string{"I", "am", "smart", "and", "strong"}`

### 3.2 实现

```go
package main

import "fmt"

func main() {
	origin := []string{"I", "am", "stupid", "and", "week"}
	fmt.Println(origin)
	changeElementsByRange(origin)
	fmt.Println(origin)
}

func changeElementsByRange(slice []string) {
	for key, value := range slice {
		if value == "stupid" {
			slice[key] = "smart"
		}

		if value == "weak" {
			slice[key] = "strong"
		}
	}
}
```