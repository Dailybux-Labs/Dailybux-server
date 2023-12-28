package model

type Config struct {
	Project Project `toml:"project"`
	MySQL   MySQL   `toml:"mysql"`
	Redis   Redis   `toml:"redis"`
}

type Project struct {
	Name string `toml:"name"`
}

type MySQL struct {
	Host     string `toml:"host"`
	Port     int64  `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
}

type Redis struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}
