package datasource

type DBconfigs struct {
	Host     string `json:"Host" yaml:"Host"`
	Port     string `json:"Port" yaml:"Port"`
	UserName string `json:"UserName" yaml:"UserName"`
	PassWord string `json:"PassWord" yaml:"PassWord"`
	DBName   string `json:"DBName" yaml:"DBName"`
	OpenNum  int    `json:"OpenNum" yaml:"OpenNum"`
	IdleNum  int    `json:"IdleNum" yaml:"IdleNum"`
	Charset  string `json:"Charset" yaml:"Charset"`
	ShowSQL  bool   `json:"ShowSQL" yaml:"ShowSQL"`
}

func DefaultDbconfig() DBconfigs {
	return DBconfigs{
		Host: "127.0.0.1",
		Port: "",
		UserName: "",
		PassWord: "",
		DBName: "taigu_dms",
		OpenNum: 500,
		IdleNum: 100,
		Charset: "utf8",
		ShowSQL: true,
	}
}
