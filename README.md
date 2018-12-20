<p align="center">
    <a href="https://github.com/ljun20160606/di"><img src="doc/di.jpeg" width="325"/></a>
</p>

<p align="center"> <code>DI(依赖注入)</code> 的 <code>Golang</code> 实现 😋</p>

<p align="center">
    🔥 <a href="#快速开始">快速入门</a>
	✨ <a href="#内容">内容目录</a>
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
* [x] 自定义生命周期

## 快速开始

下载

```sh
$ go get github.com/ljun20160606/di
```

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

## 内容

* [放入容器](#放入容器)
* [自动注入](#自动注入)
* [主动从容器中获取](#主动从容器中获取)
* [加载配置](#加载配置)

## DI

### 放入容器

```go
package main

import "github.com/ljun20160606/di"

func main() {
	name := "duck"
	// 只支持指针类型
	di.Put(&name)
}

```

### 自动注入

`di:"*"`代表根据类型自动注入

```go
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

```

### 主动从容器中获取

`name`可以根据日志获得，规则为package+typeName

```go
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

```

### 加载配置

`config`插件的关键字为`#`，可以自定义插件详细看文件[plugin_config.go](plugin_config.go),共支持`di.JSON` | `di.YAML` | `di.TOML`

```go
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

```