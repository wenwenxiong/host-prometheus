package apiserver

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/wenwenxiong/host-prometheus/pkg/client/cache"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring/prometheus"
	"github.com/wenwenxiong/host-prometheus/pkg/client/mysql"

)

const (
	// DefaultConfigurationName is the default name of configuration
	defaultConfigurationName = "conf"
	// DefaultConfigurationPath the default location of the configuration file
	defaultConfigurationPath = "/home/xww/conf"
)

// Config defines everything needed for apiserver to deal with external services
type Config struct {
	MonitoringOptions     *prometheus.Options `mapstructure:"monitoring"`
	RedisOptions          *cache.Options `mapstructure:"redis"`
	MysqlOptions		  *mysql.Options `mapstructure:"mysql"`
}

// newConfig creates a default non-empty Config
func New() *Config {
	return &Config{
		RedisOptions:          cache.NewRedisOptions(),
		MonitoringOptions:     prometheus.NewPrometheusOptions(),
		MysqlOptions:	       mysql.NewMySQLOptions(),
	}
}

// TryLoadFromDisk loads configuration from default location after server startup
// return nil error if configuration file not exists
func TryLoadFromDisk() (*Config, error) {
	viper.SetConfigName("conf")
	viper.SetConfigFile("/home/xww/conf/conf.yaml")
	viper.AddConfigPath("/home/xww/conf")

	// Load from current working directory, only used for debugging
	viper.AddConfigPath(".")

	if  err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		} else {
			return nil, fmt.Errorf("error parsing configuration file %s", err)
		}
	}

	conf := New()

	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}