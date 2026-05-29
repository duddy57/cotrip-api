package config

import "os"

type ENV struct {
	DatabaseConfig DbConfig
	ServerConfig   ServerConfig
	MailerConfig   MailerConfig
}

type DbConfig struct {
	CotripDatabaseUser     string
	CotripDatabasePassword string
	CotripDatabaseHost     string
	CotripDatabasePort     string
	CotripDatabaseName     string
}

type ServerConfig struct {
	Port        string
	Environment string
}

type MailerConfig struct {
	ApiKey    string
	EmailFrom string
}

func LoadEnv() ENV {
	return ENV{
		DatabaseConfig: DbConfig{
			CotripDatabaseUser:     getEnvString("COTRIP_DATABASE_USER", "cotrip_database_user"),
			CotripDatabasePassword: getEnvString("COTRIP_DATABASE_PASSWORD", "cotrip_database_password"),
			CotripDatabaseHost:     getEnvString("COTRIP_DATABASE_HOST", "localhost"),
			CotripDatabasePort:     getEnvString("COTRIP_DATABASE_PORT", "5432"),
			CotripDatabaseName:     getEnvString("COTRIP_DATABASE_NAME", "cotrip_database_name"),
		},

		ServerConfig: ServerConfig{
			Port:        getEnvString("COTRIP_SERVER_PORT", "8080"),
			Environment: getEnvString("COTRIP_SERVER_ENVIRONMENT", "dev"),
		},
		MailerConfig: MailerConfig{
			ApiKey:    getEnvString("COTRIP_MAIL_KEY", ""),
			EmailFrom: getEnvString("COTRIP_MAIL_FROM", "john@acme.com"),
		},
	}
}

func getEnvString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
