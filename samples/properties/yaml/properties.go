package main

import (
	"fmt"
	"github.com/ljun20160606/di"
)

type Duck struct {
	Name string `di:"#.name"`
}

func main() {
	di.ConfigLoad(`name: duck`, di.YAML)
	//di.ConfigLoadFile("path", di.YAML)
	//di.ConfigLoadReader(reader, di.YAML)
	duck := Duck{}
	di.Put(&duck)
	di.Start()

	fmt.Println(duck.Name == "duck")
}
