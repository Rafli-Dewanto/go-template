package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadDatabaseConfig(filePath string) (*DatabaseConfig, error) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ini file: %v", err)
	}

	dbSection := cfg.Section("database")

	config := &DatabaseConfig{
		Driver:   dbSection.Key("driver").String(),
		Host:     dbSection.Key("host").String(),
		Port:     dbSection.Key("port").String(),
		User:     dbSection.Key("user").String(),
		Password: dbSection.Key("password").String(),
		DBName:   dbSection.Key("dbname").String(),
		SSLMode:  dbSection.Key("sslmode").String(),
	}

	return config, nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
}
