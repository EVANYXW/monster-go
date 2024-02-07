package configure

type IConfigure interface {
	ReSet()
	Load()
	Get()
}

type globalConfig struct {
	Map map[string]IConfigure
}

var (
	Global = New()
)

func New() *globalConfig {
	return &globalConfig{
		Map: make(map[string]IConfigure),
	}
}

func init() {
	Global.Register(Best, &BetConfigM{})
}

func (c *globalConfig) Register(name string, config IConfigure) {
	if _, ok := c.Map[name]; !ok {
		c.Map[name] = config
	}
}

func (c *globalConfig) Load() {
	for _, config := range c.Map {
		config.Load()
	}
}

func (c *globalConfig) Reset() {
	for _, config := range c.Map {
		config.ReSet()
	}
}

func (c *globalConfig) Get(name string) IConfigure {
	return c.Map[name]
}
