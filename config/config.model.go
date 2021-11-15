package config

// Config App config model.
type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Port    uint16 `json:"port"`
	CPU     int    `json:"cpu"`
	Jwt     string `json:"jwt"`
	Secret  string `json:"secret"`
	Mongo   struct {
		Address        string `json:"address"`
		Database       string `json:"database"`
		User           string `json:"user"`
		Password       string `json:"password"`
		MaxConnections int    `json:"maxConnections" yaml:"maxConnections"`
		Timeout        int    `json:"timeout"`
		Mechanism      string `json:"mechanism"`
		AuthSource     string `json:"authSource"`
		Debug          bool   `json:"debug"`
	} `json:"mongo"`
	Log struct {
		Filename   string `json:"filename"`
		MaxSize    int    `json:"maxSize" yaml:"maxSize"`
		MaxBackups int    `json:"maxBackups" yaml:"maxBackups"`
		MaxAge     int    `json:"maxAge" yaml:"maxAge"`
	} `json:"log"`
}
