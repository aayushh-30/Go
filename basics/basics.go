package basics

func Add(a int, b int) int {
	return a + b
}

func Swap(a string, b string) (string, string) {
	return b, a
}

func GetDiscount(price int) (offeredPrice int) {
	offeredPrice = price - 20
	return
}

const pi = 3.14
