<p align="center">
    <a href="https://github.com/ljun20160606/di"><img src="doc/di.jpeg" width="325"/></a>
</p>

<p align="center"> <code>DI(依赖注入)</code> 的 <code>Golang</code> 实现 😋</p>

<p align="center">
    🔥 <a href="#快速开始">快速入门</a>
</p>

<p align="center">
    <a href="https://github.com/ljun20160606/di/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg"></a>
    <a href="https://travis-ci.org/ljun20160606/di"><img src="https://travis-ci.org/ljun20160606/di.svg?branch=master"></a>
    <a href="https://codecov.io/gh/ljun20160606/di"><img src="https://codecov.io/gh/ljun20160606/di/branch/master/graph/badge.svg"></a>
    <a href="http://commitizen.github.io/cz-cli"><img src="https://img.shields.io/badge/commitizen-friendly-brightgreen.svg"></a>
</p>

***

## 功能特性

* [x] 根据注解依赖注入
* [x] 根据名称获取实体

## 快速开始

创建一个main文件

```go
package main

import (
	"github.com/ljun20160606/di"
)

func init() {
	ducks := new(Ducks)
	di.Put(ducks)
}

type Duck interface {
	Quack()
}

type Ducks struct {
	Duck Duck `di:"*"`
}

type DarkDuck struct{}

func (*DarkDuck) Quack() {
	panic("implement me")
}

func main() {
	duck := new(DarkDuck)
	di.Put(duck)

	di.Start()

	ducks := di.GetWithName("main.Ducks").(*Ducks)
	if ducks.Duck.(*DarkDuck) != duck {
		panic("error")
	}
}
```