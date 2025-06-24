package h_config

var (
	Cfg Config
)

type Config struct {
	User        User      `ymal:"User"`
	Log         LogConfig `ymal:"Log"`
	ServiceName string    `ymal:"ServiceName"`
	Env         string    `ymal:"Env"`
}

type User struct {
	Name string `yaml:"Name"`
	Sex  bool   `yaml:"Sex"`
	Age  int    `yaml:"Age"`
}

type LogConfig struct {
	Level    string `yaml:"Level"`
	Path     string `yaml:"Path"`
	Interval int    `yaml:"Interval"`
}
