package config

type Config struct {
	Server Server `yaml:"server"`
	Rules  []Rule `yaml:"rules"`
}

// ---- Server ----

type Server struct {
	Host  string `yaml:"host"`
	Ports Ports  `yaml:"ports"`
	TLS   TLS    `yaml:"tls"`
}

type Ports struct {
	HTTP  HTTP  `yaml:"http"`
	HTTPS HTTPS `yaml:"https"`
}

type TLS struct {
	Key  string `yaml:"key"`
	Cert string `yaml:"cert"`
}

type HTTP string

func (p HTTP) Enabled() bool { return p != "" }

type HTTPS string

func (p HTTPS) Enabled() bool { return p != "" }

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
	Forward  Forward  `yaml:"forward"`
	Redirect Redirect `yaml:"redirect"`
	Respond  Respond  `yaml:"respond"`
}

type Middleware struct {
	Cache Cache `yaml:"cache"`
}

// ---- Handlers ----

type Forward string

func (f Forward) Enabled() bool { return f != "" }

type Redirect string

func (r Redirect) Enabled() bool { return r != "" }

type Respond struct {
	Status  int               `yaml:"status"`
	Body    string            `yaml:"body"`
	Headers map[string]string `yaml:"headers"`
}

func (r Respond) Enabled() bool { return r.Status != 0 || r.Body != "" || len(r.Headers) > 0 }

// ---- Middlewares ----

type Cache struct {
	TTL string `yaml:"ttl"`
}

func (c Cache) Enabled() bool { return c.TTL != "" }
