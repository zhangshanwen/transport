package conf

var (
	C *Conf
)

type (
	Conf struct {
		Host        string   `yaml:"host"`
		Port        string   `yaml:"port"`
		Rpc         string   `yaml:"rpc"`
		Pid         int      `yaml:"pid"`
		MaxConnect  int32    `yaml:"maxConnect"`
		ConnectTime int      `yaml:"connectTime"`
		Module      []Module `yaml:"module"`
		Nodes       []Node   `json:"nodes"`
	}
	Module struct {
		Name    string `yaml:"name"`
		Cmd     string `yaml:"cmd"`
		Prefix  string `yaml:"prefix"`
		Scheme  string `yaml:"scheme"`
		Replica int    `yaml:"replica"`
	}
	Node struct {
		Ip   string `json:"ip"`
		Port string `json:"port"`
		Rpc  string `json:"rpc"`
	}
)
