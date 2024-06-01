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
	MongoConfig
}

type MongoConfig struct {
	URI        string
	Host       string
	ReplicaSet string
	Username   string
	Password   string
	DbName     string
	Stage      string
	Tls        string
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

	viper.BindEnv("MONGO_HOST")
	viper.BindEnv("MONGO_REPLICASET")
	viper.BindEnv("MONGO_USERNAME")
	viper.BindEnv("MONGO_PASSWORD")
	viper.BindEnv("MONGO_DBNAME")
	viper.BindEnv("MONGO_URI")
	viper.BindEnv("STAGE")

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
			MongoConfig: MongoConfig{
				URI:        viper.GetString("MONGO_URI"),
				Host:       viper.GetString("MONGO_HOST"),
				ReplicaSet: viper.GetString("MONGO_REPLICASET"),
				Username:   viper.GetString("MONGO_USERNAME"),
				Password:   viper.GetString("MONGO_PASSWORD"),
				DbName:     viper.GetString("MONGO_DBNAME"),
				Stage:      viper.GetString("STAGE"),
			},
			PSQLConfig: PSQLConfig{
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
