package main

import (
	"fmt"
	"github.com/ljun20160606/di"
)

type Duck struct {
	Name string `di:"#.name"`
}

func main() {
	di.ConfigLoad(`{"name":"duck"}`, di.JSON)
	//di.ConfigLoadFile("path", di.JSON)
	//di.ConfigLoadReader(reader, di.JSON)
	duck := Duck{}
	di.Put(&duck)
	di.Start()

	fmt.Println(duck.Name == "duck")
}
