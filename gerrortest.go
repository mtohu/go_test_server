package main

import (
	"errors"
	"fmt"
	"runtime/debug"
)

func funcA() error {
	var err error
	defer func() {

		if p := recover(); p != nil {
			fmt.Println("panic recover! p: %v", p)
			str, ok := p.(string)
			if ok {
				err = errors.New(str)
			}else{
				err = errors.New("panic")
			}
			debug.PrintStack()
			fmt.Println(err.Error())
			//return err
		}
	}()
	fmt.Println("=======1111======")
	err = funcB()
	fmt.Println("=======0000======")

	return err
}
func funcB() error {
	// simulation
	panic("foo")
	return errors.New("success")
}

func test() {
	err := funcA()
	if err == nil {
		fmt.Println("err is nil\\n")
	} else {
		fmt.Println("err is %v\\n", err)
	}
}
func main() {
	test()
		//var a int
		// fn1
		/*defer func() {
			a = 3
			if err := recover(); err != nil {
				a = 4
				fmt.Println("++++")
				///f := err.(func() string)
				//fmt.Println(err, f(), reflect.TypeOf(err).Kind().String())
			} else {
				fmt.Println("fatal")
			}
		}()*/
		// fn2
		//defer func() {
		//	a = 2
		//	if r := recover(); r != nil {  // 这里的recover()去掉感受一下效果
		//		panic(r)
		//	}
			/*panic(func() string {
				return "defer panic"
			})*/
		//	fmt.Println("=====")
		//	panic("panic3")
		//}()
		//a = 1
		//panic("panic1")  // 这里的panic去掉感受一下效果
		// panic("panic2")

}
