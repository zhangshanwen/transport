package conf

var (
	C *Conf
)

type Conf struct {
	Host          string        `yaml:"host"`
	Port          string        `yaml:"port"`
	DB            DB            `yaml:"db"`
	Authorization Authorization `yaml:"authorization"`
	Level         string        `yaml:"level"`
	ResetPassword string        `yaml:"resetPassword"`
}
