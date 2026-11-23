package pro

import "fmt"
import "errors"

type Rectangle struct {
	Length int
	Width  int
}

func (r Rectangle) Area() int {
	return r.Length * r.Width
}

type flt int

func (f flt) Abs() int {
	if f < 0 {
		return int(-f)
	}
	return int(f)
}

func (r *Rectangle) increaseByOne(){
	r.Length += 1
	r.Width += 1
}

func increaseByTen(r *Rectangle){
	r.Length += 10
	r.Width += 10
}

type Animal interface{
	Sound() string
}

type Dog struct{
	name string
}

func (d Dog) Sound() string {
	return "Barks"
}

func Divide(a int, b int)(int, error){
	if b == 0 {
		return 0, errors.New("Cannot divide by zero")
	}

	return a / b, nil
}

func Multiply(a int, b int) int {
	defer panicHandler()
	if b == 0{
		panic("Teri Maa Kaa")
	}
	return a * b
}

func panicHandler(){
	r := recover()
	fmt.Println(r)
}

func Add[T int | float32](a T,b T) T {
	return a + b
}

func PrintEverything[A any](k A) A {
	return k
}

func IsEqual[T comparable](a T, b T) bool{
	return a == b
}


func Pro() {
	//panic("Panic")
	var A Animal
	dog := Dog{
		name:"Rocky",
	}
	A = dog
	fmt.Println(A)

	res :=Multiply(12,0)

	fmt.Println(res)
	

	
	
}