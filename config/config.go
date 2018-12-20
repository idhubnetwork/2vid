package config

import (
	"fmt"
	"io/ioutil"

	"github.com/olebedev/config"
)

type Redis struct {
	Url      string
	Password string
}

type Mysql struct {
	Username string
	Password string
	Host     string
	Port     string
	Dbname   string
}

type Config struct {
	Redis
	Mysql
	LogName     string
	LogLevel    string
	DESTINATION string
	Username    string
	Password    string
}

var V Config

func init() {
	file, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		panic(err)
	}

	env, err := cfg.String("environment")
	if err != nil {
		panic(err)
	}

	cfg, err = cfg.Get(env)
	if err != nil {
		panic(err)
	}

	V.LogName, err = cfg.String("logger.logName")
	if err != nil {
		panic(err)
	}
	V.LogLevel, err = cfg.String("logger.logLevel")
	if err != nil {
		panic(err)
	}

	V.DESTINATION, err = cfg.String("DESTINATION")
	if err != nil {
		panic(err)
	}

	V.Mysql.Username, err = cfg.String("mysql.username")
	if err != nil {
		panic(err)
	}
	V.Mysql.Password, err = cfg.String("mysql.password")
	if err != nil {
		panic(err)
	}
	V.Mysql.Host, err = cfg.String("mysql.host")
	if err != nil {
		panic(err)
	}
	V.Mysql.Port, err = cfg.String("mysql.port")
	if err != nil {
		panic(err)
	}
	V.Mysql.Dbname, err = cfg.String("mysql.dbname")
	if err != nil {
		panic(err)
	}

	V.Redis.Url, err = cfg.String("redis.url")
	if err != nil {
		panic(err)
	}
	V.Redis.Password, err = cfg.String("redis.password")
	if err != nil {
		panic(err)
	}

	fmt.Println(V)
}
