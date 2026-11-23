package main

import (
	"basics/basics"
	"basics/calculator"
	"basics/moderate"
	"basics/pro"
	"fmt"
)

func fun1() {
	m1 := make(map[int][]int)
	m1[0] = []int{1, 2, 3}
	m1[1] = []int{4, 5, 6}
	m1[2] = []int{7, 8, 9}

	for key, values := range m1 {
		fmt.Printf("key: %d, values:", key)
		for _, value := range values {
			fmt.Printf(" %d", value)
		}
		fmt.Println()
	}
}

func calculatingFrequency(numbers []int) (mp map[int]int) {
	mp = make(map[int]int)
	for _, val := range numbers {
		mp[val]++
	}
	return
}

func main() {
	mp1 := calculatingFrequency([]int{1, 2, 2, 1, 3, 4, 5, 5, 6, 4})
	for key, values := range mp1 {
		fmt.Printf("key: %d, frequency: %d\n", key, values)
	}
	fmt.Println("Hello, World!")
	fun1()
	fmt.Println(basics.Add(1, 2))
	fmt.Println(basics.Swap("Yo", "Shak"))
	fmt.Println(basics.GetDiscount(100))

	// var name string, age int --> Invalid
	var (
		name string     = "Ayush"
		age  int        = 23
		iota complex128 = 2 + 5i
	)
	// var name, age string --> Valid
	var class, sex = "UG", "M"

	fmt.Println(name, age, class, sex, iota)

	school := "BPS"

	fmt.Println(school)

	mathsMarks := 26

	fmt.Println(mathsMarks)

	updatedMathMarks := float64(mathsMarks) + 0.000001
	fmt.Println(updatedMathMarks)

	var a = 10
	var b = 20

	fmt.Println(calculator.CalculatorAPI(a, b, "Mul"))

	fmt.Println(moderate.OddEven(10))
	moderate.PrintBlock(3)
	moderate.PrintPyramid(5)
	moderate.DefTest(3)

	pro.Pro()

}
