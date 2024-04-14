package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	address          = "localhost"
	port             = 8080
	loggerLevel      = "debug"
	translateMode    = false
	envFile          = ".env"
	postgresPort     = 5432
	defaultLimit     = 10
	defaultOffset    = 0
	expireCookieDays = 7
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
		Port    uint64 `yaml:"port"`
	} `yaml:"server"`
	Postgresql struct {
		User    string `yaml:"user"`
		DbName  string `yaml:"dbname"`
		Host    string `yaml:"host"`
		Port    uint64 `yaml:"port"`
		SslMode string `yaml:"sslmode"`
		Limit   uint64 `yaml:"limit"`
		Offset  uint64 `yaml:"offset"`
	} `yaml:"postgresql"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	LoggerLvl      string `yaml:"logger_level"`
	CookieSettings CookieSettings
}

type CookieSettings struct {
	Secure     bool `yaml:"secure"`
	HttpOnly   bool `yaml:"http_only"`
	ExpireDate struct {
		Years  uint64 `yaml:"years"`
		Months uint64 `yaml:"months"`
		Days   uint64 `yaml:"days"`
	} `yaml:"expire_date"`
}

func New() *Config {
	return &Config{
		Server: struct {
			Address string `yaml:"address"`
			Port    uint64 `yaml:"port"`
		}(struct {
			Address string
			Port    uint64
		}{
			Address: address,
			Port:    port,
		}),
		Postgresql: struct {
			User    string `yaml:"user"`
			DbName  string `yaml:"dbname"`
			Host    string `yaml:"host"`
			Port    uint64 `yaml:"port"`
			SslMode string `yaml:"sslmode"`
			Limit   uint64 `yaml:"limit"`
			Offset  uint64 `yaml:"offset"`
		}(struct {
			User    string
			DbName  string
			Host    string
			Port    uint64
			SslMode string
			Limit   uint64
			Offset  uint64
		}{
			User:    "postgres",
			DbName:  "cyber_garden",
			Host:    "localhost",
			Port:    postgresPort,
			SslMode: "disable",
			Limit:   defaultLimit,
			Offset:  defaultOffset,
		}),
		Redis: struct {
			Addr     string `yaml:"addr"`
			Password string `yaml:"password"`
		}(struct {
			Addr     string
			Password string
		}{
			Addr:     "localhost:6379",
			Password: "qewrty",
		}),
		CookieSettings: struct {
			Secure     bool `yaml:"secure"`
			HttpOnly   bool `yaml:"http_only"`
			ExpireDate struct {
				Years  uint64 `yaml:"years"`
				Months uint64 `yaml:"months"`
				Days   uint64 `yaml:"days"`
			} `yaml:"expire_date"`
		}(struct {
			Secure     bool
			HttpOnly   bool
			ExpireDate struct {
				Years  uint64
				Months uint64
				Days   uint64
			}
		}{
			Secure:   true,
			HttpOnly: true,
			ExpireDate: struct {
				Years  uint64
				Months uint64
				Days   uint64
			}{
				Years:  0,
				Months: 0,
				Days:   expireCookieDays,
			},
		}),
	}
}

func (c *Config) Open(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(c); err != nil {
		return err
	}

	return nil
}

func (c *Config) FormatDbAddr() string {
	return fmt.Sprintf(
		"host=%s user=%s password=admin dbname=%s port=%d sslmode=%s",
		c.Postgresql.Host,
		c.Postgresql.User,
		c.Postgresql.DbName,
		c.Postgresql.Port,
		c.Postgresql.SslMode,
	)
}

func ParseFlag(path *string) {
	flag.StringVar(path, "ConfigPath", getPwd()+"/configs/app/local.yml", "Path to Config")
}

func getPwd() string {
	pwd, _ := os.Getwd()
	return pwd
}
