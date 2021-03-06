package di

type (
	Plugin interface {
		Load(path string, ice Ice)
		Prefix() string
		Priority() int
	}
	Plugins []Plugin
)

func (pl Plugins) Len() int {
	return len(pl)
}

func (pl Plugins) Less(i, j int) bool {
	return pl[i].Priority() < pl[j].Priority()
}

func (pl Plugins) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}
