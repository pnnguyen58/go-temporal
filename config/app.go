package config

import (
	"github.com/spf13/viper"
	"log"
	"path"
	"path/filepath"
	"runtime"
)

type config struct {
	Secret struct {
		Region   string `mapstructure:"REGION"`
		Database string `mapstructure:"DATABASE"`
		Logger   string `mapstructure:"LOGGER"`
	} `mapstructure:"SECRET"`

	Server struct {
		Name     string `mapstructure:"NAME"`
		GRPCPort int    `mapstructure:"GRPC_PORT"`
		HTTPPort int    `mapstructure:"HTTP_PORT"`
		TempoPort int    `mapstructure:"TEMPO_PORT"`
		TempoNameSpace string    `mapstructure:"TEMPO_NAMESPACE"`
	} `mapstructure:"SERVER"`
}

var C config

func ReadConfig() {
	Config := &C

	viper.SetConfigName(".env")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join(rootDir(), "configs"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalln(err)
	}
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
