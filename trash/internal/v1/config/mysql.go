package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type MysqlConfig struct {
	Host     string
	Port     string
	DbName   string
	Username string
	Password string
}

func NewMysqlConfig() (dns string, err error) {
	viper.SetConfigFile("/Users/zz/GolandProjects/web-demo/vote/config/mysql.toml")
	viper.AddConfigPath("/Users/zz/GolandProjects/web-demo/vote/config/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore errno if desired
			log.Println("no such config file")
			return "", err
		} else {
			// Config file was found but another errno was produced
			log.Println("read config errno")
			return "", err
		}
	}

	host := viper.GetString(`mysql.host`)
	username := viper.GetString(`mysql.username`)
	password := viper.GetString(`mysql.password`)
	port := viper.GetString(`mysql.port`)
	dbname := viper.GetString(`mysql.db_name`)
	log.Println(host, username, password, port, dbname)

	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		username, password, host, port, dbname)
	log.Println(dns)
	return dns, nil
}
