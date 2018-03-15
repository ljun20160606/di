package main

import (
	"fmt"
	"github.com/ljun20160606/di"
)

type Duck struct {
	Name string `di:"#.name"`
}

func main() {
	di.TomlLoad(`name = "duck"`)
	//di.TomlLoadFile("path")
	//di.TomlLoadReader(reader)
	duck := Duck{}
	di.Put(&duck)
	di.Start()

	fmt.Println(duck.Name == "duck")
}
