package config

// HTTP : HTTP configuration
type HTTP struct {
	Host string `yaml:"host" validate:"nonzero"`
	Port int    `yaml:"port" validate:"nonzero"`
}

// Init : init HTTP
func (h *HTTP) Init() {
	h.Host = "127.0.0.1"
	h.Port = 8080
}
