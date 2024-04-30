package env

import (
	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort    string
	HttpPort    string
	UserService service
	OrgService  service
	Database    *DatabaseConfig
}

type service struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	PSQLConfig
}

type PSQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
	SslMode  string
}

var Settings *Config

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	viper.BindEnv("GRPC_PORT")
	viper.BindEnv("HTTP_PORT")

	viper.BindEnv("USERS_SERVICE_HOST")
	viper.BindEnv("USERS_SERVICE_PORT")

	viper.BindEnv("ORGS_SERVICE_HOST")
	viper.BindEnv("ORGS_SERVICE_PORT")

	viper.BindEnv("POSTGRES_HOST")
	viper.BindEnv("POSTGRES_PORT")
	viper.BindEnv("POSTGRES_USERNAME")
	viper.BindEnv("POSTGRES_PASSWORD")
	viper.BindEnv("POSTGRES_DBNAME")
	viper.BindEnv("POSTGRES_SSL_MODE")

	Settings = &Config{
		GrpcPort: viper.GetString("GRPC_PORT"),
		HttpPort: viper.GetString("HTTP_PORT"),
		UserService: service{
			Host: viper.GetString("USERS_SERVICE_HOST"),
			Port: viper.GetString("USERS_SERVICE_PORT"),
		},
		OrgService: service{
			Host: viper.GetString("ORGS_SERVICE_HOST"),
			Port: viper.GetString("ORGS_SERVICE_PORT"),
		},
		Database: &DatabaseConfig{
			PSQLConfig{
				Host:     viper.GetString("POSTGRES_HOST"),
				Port:     viper.GetInt("POSTGRES_PORT"),
				Username: viper.GetString("POSTGRES_USERNAME"),
				Password: viper.GetString("POSTGRES_PASSWORD"),
				DbName:   viper.GetString("POSTGRES_DBNAME"),
				SslMode:  viper.GetString("POSTGRES_SSL_MODE"),
			},
		},
	}
}
