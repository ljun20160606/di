package main

import (
	"fmt"
	"github.com/ljun20160606/di"
)

type Duck struct {
	Name string `di:"#.name"`
}

func main() {
	di.ConfigLoad(`name = "duck"`, di.TOML)
	//di.ConfigLoadFile("path", di.TOML)
	//di.ConfigLoadReader(reader, di.TOML)
	duck := Duck{}
	di.Put(&duck)
	di.Start()

	fmt.Println(duck.Name == "duck")
}
