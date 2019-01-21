package main

import "github.com/ljun20160606/di"

func main() {
	di.Configure(di.WithLoggerVerbose(true))
	name := "duck"
	// 只支持指针类型
	di.Put(&name)
}
