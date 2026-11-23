package calculator

func Add(a, b int) int {
	return a + b
}

func Subtract(a, b int) int {
	return a - b
}

func Multiply(a, b int) int {
	return a * b
}

func Divide(a, b int) int {
	return a / b
}

func CalculatorAPI(a, b int, op string) int {
	if op == "Add" {
		return Add(a, b)
	} else if op == "Sub" {
		return Subtract(a, b)
	} else if op == "Mul" {
		return Multiply(a, b)
	} else if op == "Div" {
		return Divide(a, b)
	} else {
		return -1
	}
}
