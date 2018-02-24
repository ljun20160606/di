package di

import (
	"reflect"
	"sort"

	"fmt"
	"github.com/ljun20160606/go-lib/reflectl"
)

func NewContainer() Container {
	return &container{
		cupMap:  make(map[string][]*Cup),
		iceMap:  make(map[string][]Ice),
		plugins: make(map[Lifecycle]Plugins),
	}
}

type (
	Container interface {
		// 注册插件 根据lifecycle决定在哪一层被初始化
		RegisterPlugin(lifecycle Lifecycle, p Plugin)

		// 放入iceMap, 根据插件来注入
		PutIce(ice Ice)

		// 放入cupMap，water名字由容器的默认规则来决定
		Put(water Water)

		// 放入cupMap，water名字自定义
		PutWithName(water Water, name string)

		// 获取除了excludedNames之外 dropType类型的cup
		GetCup(dropType reflect.Type, excludedNames ...string) (h *Cup)

		// 根据name获取 dropType类型的cup
		GetCupWithName(name string, dropType reflect.Type) (h *Cup)

		// 根据name获取 water
		GetWaterWithName(name string) Water

		// 深度遍历cupMap
		EachCup(cupFunc CupFunc)

		// 开始根据生命周期初始化
		Start()
	}

	container struct {
		// 存放 water
		cupMap map[string][]*Cup

		// 存放 ice
		iceMap map[string][]Ice

		// 存放 插件
		plugins map[Lifecycle]Plugins
	}
)

func (c *container) RegisterPlugin(lifecycle Lifecycle, p Plugin) {
	if _, ok := c.plugins[lifecycle]; !ok {
		c.plugins[lifecycle] = []Plugin{}
	}
	c.plugins[lifecycle] = append(c.plugins[lifecycle], p)
}

func (c *container) PutIce(ice Ice) {
	prefix := ice.Prefix()
	if _, has := c.iceMap[prefix]; !has {
		c.iceMap[prefix] = []Ice{ice}
		return
	}
	c.iceMap[prefix] = append(c.iceMap[prefix], ice)
}

func (c *container) putWater(water Water, name string) {
	v := reflect.ValueOf(water)
	t := v.Type()
	switch kind := t.Kind(); kind {
	case reflect.Ptr:
	default:
		panic(ErrorType.Panic(kind))
	}
	if name == "" {
		name = reflectl.GetValueDefaultName(v)
	}
	logger.Output(4, fmt.Sprintf("放入 %v", name))
	// 额，没想到并发的场景所以没加锁
	if _, ok := c.cupMap[name]; !ok {
		cup := newCup(water, t, v, c)
		c.cupMap[name] = append(c.cupMap[name], cup)
	}
}

func (c *container) Put(water Water) {
	c.putWater(water, "")
}

func (c *container) PutWithName(water Water, name string) {
	c.putWater(water, name)
}

func (c *container) GetWaterWithName(name string) Water {
	if hs, ok := c.cupMap[name]; ok {
		switch len(hs) {
		case 0:
			return nil
		default:
			return hs[0].Water
		}
	}
	return nil
}

func (c *container) GetCupWithName(name string, t reflect.Type) (h *Cup) {
	if cups, found := c.cupMap[name]; found {
		for _, cup := range cups {
			if reflectl.TypeEqual(t, cup.Class) {
				return cup
			}
		}
	}
	return nil
}

func (c *container) GetCup(t reflect.Type, excludedNames ...string) (h *Cup) {
	switch len(excludedNames) {
	case 0:
		c.EachCup(func(name string, cup *Cup) bool {
			if reflectl.TypeEqual(t, cup.Class) {
				h = cup
				return true
			}
			return false
		})
	case 1:
		c.EachCup(func(name string, cup *Cup) bool {
			if name != excludedNames[0] {
				if reflectl.TypeEqual(t, cup.Class) {
					h = cup
					return true
				}
			}
			return false
		})
	default:
		s := make(map[string]struct{})
		for _, v := range excludedNames {
			s[v] = struct{}{}
		}
		c.EachCup(func(name string, cup *Cup) bool {
			if _, has := s[name]; !has {
				if reflectl.TypeEqual(t, cup.Class) {
					h = cup
					return true
				}
			}
			return false
		})
	}
	return
}

func (c *container) EachCup(cupFunc CupFunc) {
	for k, v := range c.cupMap {
		for _, vv := range v {
			if cupFunc(k, vv) {
				return
			}
		}
	}
}

func (c *container) injectDependency() {
	c.EachCup(func(name string, cup *Cup) bool {
		cup.injectDependency()
		return false
	})
}

func (c *container) loadPlugins(lifecycle Lifecycle) {
	plugins, ok := c.plugins[lifecycle]
	if !ok || len(plugins) == 0 {
		return
	}
	sort.Sort(plugins)
	for _, plugin := range plugins {
		ices, ok := c.iceMap[plugin.Prefix()]
		if !ok {
			continue
		}
		for _, ice := range ices {
			ice.LoadPlugin(plugin)
		}
	}
}

func (c *container) Start() {
	c.injectDependency()
	c.loadPlugins(BeforeInit)
	c.init()
	c.loadPlugins(BeforeReady)
	c.ready()
}

func (c *container) init() {
	set := make(CupSet)
	c.EachCup(func(name string, cup *Cup) bool {
		cup.init(set)
		return false
	})
}

func (c *container) ready() {
	set := make(CupSet)
	c.EachCup(func(name string, cup *Cup) bool {
		cup.ready(set)
		return false
	})
}
