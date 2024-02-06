package configure

type IConfigure interface {
	ReSet()
	Load()
	Get()
}

var (
	Map = make(map[string]IConfigure)
)

func Init() {
	Register(Best, &BetConfigM{})
}

func Register(name string, config IConfigure) {
	if _, ok := Map[name]; !ok {
		Map[name] = config
	}
}

func Load() {
	Init()
	for _, config := range Map {
		config.Load()
	}
}

func Reset() {
	for _, config := range Map {
		config.ReSet()
	}
}

func Get(name string) IConfigure {
	return Map[name]
}
