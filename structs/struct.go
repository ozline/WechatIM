package structs

type Conf struct {
	HttpPort string `yaml:"httpPort"`
	Mysql    struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Address  string `yaml:"ip"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"dbname"`
		Table    struct {
			Users string `yaml:"users"`
			Test  string `yaml:"test"`
		}
	}
}