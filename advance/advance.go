package main

import (
	"fmt"
)

type Details struct{
	name string
	age int
	gender string
}

func main(){
	p := 10
	var addr *int = &p

	*addr = 100

	fmt.Println(*addr)

	d := Details{
		name: "Ayush",
		age: 23,
		gender: "M",
	}

	fmt.Println(d.name)

	var info [2]Details
	info[0] = d

	fmt.Println(info[1])

	k := info[0:2]
	fmt.Println(k)
	k[1].name = "Raj"
	fmt.Println(info)

	list1 := []string{"a","b","c","d"}
	fmt.Println(list1)

	list2 := []struct {
		id int
		name string
	}{
		{
			id:1,
			name:"Ayush",
		},
		{
			id:2,
			name:"Kanha",
		},
	}

	fmt.Println(list2)

}