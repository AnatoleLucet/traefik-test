package config

type Config struct {
	Server Server `yaml:"server"`
	Rules  []Rule `yaml:"rules"`
}

// ---- Server ----

type Server struct {
	Host  string `yaml:"host"`
	Ports Ports  `yaml:"ports"`
}

type Ports struct {
	HTTP  string `yaml:"http"`
	HTTPS string `yaml:"https"`
}

// ---- Rule ----

type Rule struct {
	If         If         `yaml:"if"`
	Then       Then       `yaml:"then"`
	Middleware Middleware `yaml:"middleware"`
}

type If struct {
	Host   string `yaml:"host"`
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

type Then struct {
	Forward  string  `yaml:"forward"`
	Redirect string  `yaml:"redirect"`
	Respond  Respond `yaml:"respond"`
}

type Middleware struct {
	Cache Cache `yaml:"cache"`
}

// ---- Actions ----

type Respond struct {
	Status int               `yaml:"status"`
	Body   string            `yaml:"body"`
	Header map[string]string `yaml:"header"`
}

// ---- Middlewares ----

type Cache struct {
	TTL string `yaml:"ttl"`
}
