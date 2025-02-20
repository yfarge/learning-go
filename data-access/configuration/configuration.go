package configuration

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Settings struct {
	Database DatabaseSettings `mapstructure:"database"`
}

type DatabaseSettings struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DatabaseName string `mapstructure:"database_name"`
}

func (self DatabaseSettings) ConnectionString() string {
	connection_url := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(self.Username, self.Password),
		Host:   fmt.Sprintf("%s:%d", self.Host, self.Port),
		Path:   self.DatabaseName,
	}

	query := connection_url.Query()
	query.Set("sslmode", "disable")
	connection_url.RawQuery = query.Encode()

	return connection_url.String()
}

func GetConfiguration() (*Settings, error) {
	base_path, err := os.Getwd()
	if err != nil {
		slog.Error("Failed to determine the current directory", "err", err)
		return nil, err
	}

	configuration_directory := filepath.Join(base_path, "configuration")
	viper.SetConfigFile(filepath.Join(configuration_directory, "base.yaml"))

	err = viper.ReadInConfig()
	if err != nil {
		slog.Error("Failed to read in the configuration file", "err", err)
		return nil, err
	}

	var settings Settings
	err = viper.Unmarshal(&settings)
	if err != nil {
		slog.Error("Failed unmarshal configuration file.")
		return nil, err
	}

	return &settings, nil
}
