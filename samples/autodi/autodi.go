package main

import (
	"github.com/ljun20160606/di"
)

type Duck struct {
	Name *string `di:"*"`
}

func main() {
	name := "duck"
	duck := Duck{}
	di.Put(&name)
	di.Put(&duck)
	// 开始依赖注入
	di.Start()
}
