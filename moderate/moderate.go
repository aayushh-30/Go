package moderate

import (
	"fmt"
)

func OddEven(num int) bool {
	if isEven := num%2 == 0; isEven {
		return true
	}
	return false
}

func DefTest(b int) {
	for i := 0; i < b; i++ {
		defer fmt.Println(i)
	}
}

func PrintBlock(n int) {
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			fmt.Print("*")
		}
		fmt.Println()
	}
}

func PrintPyramid(n int) {
	for i := 1; i <= n; i++ {
		for j := 1; j <= i; j++ {
			fmt.Print("*")
		}
		fmt.Println()
	}
}
