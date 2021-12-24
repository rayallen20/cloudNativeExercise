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












































