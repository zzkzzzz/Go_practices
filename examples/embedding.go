package main

import "fmt"

type base struct {
	num int
}

func (b base) describe() string {
	return fmt.Sprintf("base with num=%v", b.num)
}

// A container embeds a base. An embedding looks like a field without a name.
type container struct {
	base
	str string
}

func main() {

	co := container{
		base: base{
			num: 1,
		},
		str: "some name",
	}

	fmt.Printf("co={num: %v, str: %v}\n", co.num, co.str)

	fmt.Println("also num:", co.base.num)

	fmt.Println("describe:", co.describe())

	type describer interface {
		describe() string
	}
	package main

	import (
		"errors"
		"fmt"
	)
	
	func f1(arg int) (int, error) {
		if arg == 42 {
	
			return -1, errors.New("can't work with 42")
	
		}
	
		return arg + 3, nil
	}
	
	type argError struct {
		arg  int
		prob string
	}
	
	func (e *argError) Error() string {
		return fmt.Sprintf("%d - %s", e.arg, e.prob)
	}
	
	func f2(arg int) (int, error) {
		if arg == 42 {
	
			return -1, &argError{arg, "can't work with it"}
		}
		return arg + 3, nil
	}
	
	func main() {
	
		for _, i := range []int{7, 42} {
			if r, e := f1(i); e != nil {
				fmt.Println("f1 failed:", e)
			} else {
				fmt.Println("f1 worked:", r)
			}
		}
		for _, i := range []int{7, 42} {
			if r, e := f2(i); e != nil {
				fmt.Println("f2 failed:", e)
			} else {
				fmt.Println("f2 worked:", r)
			}
		}
	
		_, e := f2(42)
		if ae, ok := e.(*argError); ok {
			fmt.Println(ae.arg)
			fmt.Println(ae.prob)
		}
	}
	var d describer = co
	fmt.Println("describer:", d.describe())
}
