package config

import (
	"flag"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port   int          `mapstructure:"port"`
	Matrix MatrixConfig `mapstructure:"matrix"` // separate interface so other code can refer to it
}

type MatrixConfig struct {
	HomeserverURL string            `mapstructure:"homeserver_url"`
	AccessToken   string            `mapstructure:"access_token"`
	MxID          string            `mapstructure:"mx_id"`
	Rooms         map[string]string `mapstructure:"rooms"` // room ID to token for particular room (TODO: improve?)
}

func Load() (error, AppConfig) {
	var cfg AppConfig

	configPath := flag.String("c", "", "path to config file")
	flag.Parse()

	v := viper.New()

	v.SetDefault("port", 8080)

	if *configPath != "" {
		v.SetConfigFile(*configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
	}

	if err := v.ReadInConfig(); err != nil {
		// Only error if explicitly specified
		if *configPath != "" {
			return err, cfg
		}

		log.WithError(err).Warn("Could not load config file, silently failing")
	}

	_ = v.BindEnv("port")

	_ = v.BindEnv("matrix.homeserver_url")
	_ = v.BindEnv("matrix.access_token")
	_ = v.BindEnv("matrix.mx_id")
	_ = v.BindEnv("matrix.rooms")

	v.SetEnvPrefix("CPANEL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.AutomaticEnv()

	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unable to decode config: %w", err), cfg
	}

	return nil, cfg
}
