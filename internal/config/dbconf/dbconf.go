package dbconf

import "fmt"

type Database struct {
	DbModel  string `yaml:"dbmodel"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DbName   string `yaml:"dbname"`
}

func (d *Database) DbUrl() string {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s", d.DbModel, d.Username, d.Password, d.Host, d.Port, d.DbName)
}
