package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type SshHostConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	PrivateKey string
}

type Config struct {
	Hosts map[string]SshHostConfig `mapstructure:"hosts"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	viper.SetEnvPrefix("ssm")
	viper.AutomaticEnv()

	// Datei einlesen
	if err := v.ReadInConfig(); err != nil {
		// config.yaml optional, nur Warnung
		fmt.Println("Keine config.yaml gefunden, nur ENV wird genutzt")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("konnte Config nicht mappen: %w", err)
	}

	if err := cfg.ValidateAndSetDefaultsSshHostConfig(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) ValidateAndSetDefaultsSshHostConfig() error {
	for name, host := range c.Hosts {
		// Default Port
		if host.Port == 0 {
			host.Port = 22
		}

		if host.Host == "" {
			return fmt.Errorf("ssh.host für Host %s fehlt", name)
		}
		if host.User == "" {
			return fmt.Errorf("ssh.user für Host %s fehlt", name)
		}

		if host.Password == "" && host.PrivateKey == "" {
			return fmt.Errorf("entweder Password oder PrivateKey muss gesetzt sein für Host %s", name)
		}
		c.Hosts[name] = host
	}
	return nil
}
