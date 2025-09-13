package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type WatchedDir struct {
	Name         string `mapstructure:"name"`
	Path         string `mapstructure:"path"`
	InodeCount   *bool  `mapstructure:"inode-count"`
	FolderSize   *bool  `mapstructure:"folder-size"`
	LastModified *bool  `mapstructure:"last-modified"`
}

type RecursiveWatchedDirs struct {
	MaxDepth     int    `mapstructure:"max-depth"`
	Path         string `mapstructure:"path"`
	InodeCount   *bool  `mapstructure:"inode-count"`
	FolderSize   *bool  `mapstructure:"folder-size"`
	LastModified *bool  `mapstructure:"last-modified"`
}

type SshHostConfig struct {
	Host              string                 `mapstructure:"host"`
	Port              int                    `mapstructure:"port"`
	User              string                 `mapstructure:"user"`
	Password          string                 `mapstructure:"password"`
	PrivateKey        string                 `mapstructure:"privateKey"`
	FilesystemMetrics *bool                  `mapstructure:"filesystem_metrics"`
	WatchedDirs       []WatchedDir           `mapstructure:"watched_dirs"`
	RecursiveDirs     []RecursiveWatchedDirs `mapstructure:"recursive_watched_dirs"`
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

	if err := cfg.validateAndSetDefaultsSshHostConfig(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validateAndSetDefaultsSshHostConfig() error {
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

		for i := range host.WatchedDirs {
			wd := &host.WatchedDirs[i]
			if err := wd.validateAndSetDefaultsWatchedDir(); err != nil {
				return fmt.Errorf("host %s, WatchedDir %s: %w", name, wd.Name, err)
			}
		}

		for i := range host.RecursiveDirs {
			wd := &host.RecursiveDirs[i]
			if err := wd.validateAndSetDefaultsRecursiveWatchedDir(); err != nil {
				return fmt.Errorf("host %s, RecursiveWatchedDir %s: %w", name, wd.Path, err)
			}
		}

		if host.FilesystemMetrics == nil {
			b := true
			host.FilesystemMetrics = &b
		}

		c.Hosts[name] = host
	}
	return nil
}

func (wd *WatchedDir) validateAndSetDefaultsWatchedDir() error {
	if wd.Path == "" {
		return fmt.Errorf("pfad darf nicht leer sein")
	}
	if wd.Name == "" {
		wd.Name = wd.Path
	}

	if wd.InodeCount == nil {
		b := true
		wd.InodeCount = &b
	}

	if wd.FolderSize == nil {
		b := true
		wd.FolderSize = &b
	}

	if wd.LastModified == nil {
		b := true
		wd.LastModified = &b
	}
	return nil
}

func (wd *RecursiveWatchedDirs) validateAndSetDefaultsRecursiveWatchedDir() error {
	if wd.Path == "" {
		return fmt.Errorf("pfad darf nicht leer sein")
	}

	if wd.InodeCount == nil {
		b := true
		wd.InodeCount = &b
	}

	if wd.FolderSize == nil {
		b := true
		wd.FolderSize = &b
	}

	if wd.LastModified == nil {
		b := true
		wd.LastModified = &b
	}
	return nil
}
