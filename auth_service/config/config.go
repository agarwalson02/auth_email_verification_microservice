package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type ServerConfig struct {
	AppVersion        string        `mapstructure:"appversion"`
	Port              string        `mapstructure:"port"`
	PprofPort         string        `mapstructure:"pprofport"`
	Mode              string        `mapstructure:"mode"`
	JwtSecretKey      string        `mapstructure:"jwtsecretkey"`
	CookieName        string        `mapstructure:"cookiename"`
	ReadTimeOut       time.Duration `mapstructure:"readtimeout"`
	WriteTimeOut      time.Duration `mapstructure:"writetimeout"`
	SSL               bool          `mapstructure:"ssl"`
	CtxDefaultTimeout time.Duration `mapstructure:"ctxdefaulttimeout"`
	CSRF              bool          `mapstructure:"csrf"`
	Debug             bool          `mapstructure:"debug"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	Driver   string `mapstructure:"driver"`
}

type LoggerConfig struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disablecaller"`
	DisableStacktrace bool   `mapstructure:"disablestacktrace"`
	Level             string `mapstructure:"level"`
	Encoding          string `mapstructure:"encoding"`
}

// load config file from path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(filename)
	v.SetConfigType("yaml") // 👈 add this
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	if err := v.Unmarshal(&c); err != nil {
		log.Printf("Unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
