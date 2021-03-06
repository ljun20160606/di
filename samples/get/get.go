package main

import (
	"fmt"
	"github.com/ljun20160606/di"
)

func main() {
	name := "duck"
	// 只支持指针类型
	di.Put(&name)
	withName := *(di.GetWithName("string").(*string))
	fmt.Println(name == withName)
}
