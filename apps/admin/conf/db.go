package conf

type DB struct {
	Mysql Mysql `yaml:"mysql"`
	Redis Redis `yaml:"redis"`
}
type Mysql struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Database    string `yaml:"database"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	LogMode     string `yaml:"logMode"`
	Config      string `yaml:"config"`
	TablePrefix string `yaml:"tablePrefix"`
}
type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database int    `yaml:"database"`
	Password string `yaml:"password"`
}
