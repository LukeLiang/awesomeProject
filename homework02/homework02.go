package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Task func()

func addTen(num *int) {
	*num = *num + 10
}

func multiply2(numSlice *[]int) {
	for index := range *numSlice {
		(*numSlice)[index] *= 2
	}
}

func printNumbers() {
	wg := sync.WaitGroup{}

	wg.Add(2)

	// 打印1-10奇数
	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 1 {
				fmt.Printf("奇数[%d]\r\n", i)
			}
		}
		wg.Done()
	}()

	// 打印2-10偶数

	go func() {
		for i := 2; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Printf("偶数[%d]\r\n", i)
			}
		}
		wg.Done()
	}()

	wg.Wait()
}

func homework03() {
	wg := sync.WaitGroup{}

	wg.Add(2)

	// 打印1-10奇数
	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 1 {
				fmt.Printf("奇数[%d]\r\n", i)
			}
		}
		wg.Done()
	}()

	// 打印2-10偶数

	go func() {
		for i := 2; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Printf("偶数[%d]\r\n", i)
			}
		}
		wg.Done()
	}()

	wg.Wait()
}

func homework04(task []Task) {
	wg := sync.WaitGroup{}
	wg.Add(len(task))

	for i := 0; i < len(task); i++ {
		go func(id int, t Task) {
			defer wg.Done()
			start := time.Now()
			t()
			duration := time.Since(start)
			fmt.Printf("任务 [%d] 执行完毕，耗时: %v\n", id, duration)
		}(i, task[i])
	}
	wg.Wait()
}

func homework07() {
	wg := sync.WaitGroup{}

	channel := make(chan int)

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			channel <- i
		}
		close(channel)
	}()

	go func() {
		defer wg.Done()
		for val := range channel {
			fmt.Println(val)
		}
	}()

	wg.Wait()
}

func homework08() {
	channel := make(chan int, 10)

	go func() {

		for i := 1; i <= 100; i++ {
			channel <- i
			fmt.Printf("channel写入【%d】\r\n", i)
		}
		close(channel)
	}()

	for val := range channel {
		fmt.Printf("接收到: %d\n", val)
	}
}

func homework09() {
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	num := 0

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				lock.Lock()
				num++
				lock.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("结果%d", num)
}

func homework10() {
	var wg sync.WaitGroup
	var num int32 = 0

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt32(&num, 1)
			}
		}()
	}

	wg.Wait()
	fmt.Printf("结果%d", atomic.LoadInt32(&num))
}

func main() {

	//编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值
	num := 10
	addTen(&num)
	fmt.Printf("相加之后的值:%d\n", num)

	// 实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
	numSlice := []int{1, 2, 3, 4, 5}
	multiply2(&numSlice)
	fmt.Printf("相乘之后的值:%d\n", numSlice)

	// 编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
	homework03()

	tasks := []Task{
		func() {
			fmt.Println("任务 0: 正在下载文件...")
			time.Sleep(2 * time.Second)
		},
		func() {
			fmt.Println("任务 1: 正在处理数据...")
			time.Sleep(1 * time.Second)
		},
		func() {
			fmt.Println("任务 2: 正在上传结果...")
			time.Sleep(500 * time.Millisecond)
		},
	}
	// 设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
	homework04(tasks)

	// 定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
	rectangle := Rectangle{
		Width: 10, Height: 5,
	}

	circle := Circle{Radius: 10}

	printInfo(rectangle)
	printInfo(circle)

	// 使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
	employee := Employee{
		EmployeeID: 10,
		Person: Person{
			Name: "Luke",
			Age:  18,
		},
	}
	PrintInfo2(employee)

	// 编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
	homework07()

	// 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
	homework08()

	// 编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	homework09()

	// 使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	homework10()
}

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width, Height float64
}

type Circle struct {
	Radius float64
}

func (rec Rectangle) Area() float64 {
	return rec.Width * rec.Height
}

func (rec Rectangle) Perimeter() float64 {
	return 2 * (rec.Width + rec.Height)
}

func (circle Circle) Area() float64 {
	return math.Pi * circle.Radius * circle.Radius
}

func (circle Circle) Perimeter() float64 {
	return 2 * math.Pi * circle.Radius
}

func printInfo(shape Shape) {
	fmt.Printf("图形信息 -> 面积: %.2f, 周长: %.2f\n", shape.Area(), shape.Perimeter())
}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person     Person
	EmployeeID int
}

func PrintInfo2(e Employee) {
	fmt.Printf("员工姓名:%s,员工年龄:%d,员工ID:%d", e.Person.Name, e.Person.Age, e.EmployeeID)
}
