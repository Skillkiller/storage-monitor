package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type WatchedDir struct {
	Name         string `mapstructure:"name"`
	Path         string `mapstructure:"path"`
	FolderSize   *bool  `mapstructure:"folder-size"`
	LastModified *bool  `mapstructure:"last-modified"`
}

type RecursiveWatchedDir struct {
	Name         string `mapstructure:"name"`
	MaxDepth     int    `mapstructure:"max-depth"`
	Path         string `mapstructure:"path"`
	FolderSize   *bool  `mapstructure:"folder-size"`
	LastModified *bool  `mapstructure:"last-modified"`
}

type Config struct {
	WatchedDirs          []WatchedDir          `mapstructure:"watched_dirs"`
	RecursiveDirs        []RecursiveWatchedDir `mapstructure:"recursive_watched_dirs"`
	RcloneExecutablePath string                `mapstructure:"rclone_executable_path"`
}

type InternalConfig struct {
	RecursiveDirs        []RecursiveWatchedDir
	RcloneExecutablePath string `mapstructure:"rclone_executable_path"`
}

func Load() (*InternalConfig, error) {
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

	internalCfg, err := cfg.validateAndConvertConfig()
	if err != nil {
		return nil, err
	}

	return internalCfg, nil
}

func (c *Config) validateAndConvertConfig() (*InternalConfig, error) {
	var recursiveDirs []RecursiveWatchedDir

	// WatchedDirs -> RecursiveDirs
	for _, d := range c.WatchedDirs {
		recursiveDirs = append(recursiveDirs, RecursiveWatchedDir{
			Name:         d.Name,
			MaxDepth:     0,
			Path:         d.Path,
			FolderSize:   d.FolderSize,
			LastModified: d.LastModified,
		})
	}

	// Copy RecursiveDirs
	for _, d := range c.RecursiveDirs {
		recursiveDirs = append(recursiveDirs, d)
	}

	for i := range recursiveDirs {
		if err := recursiveDirs[i].setDefaultsAndValidate(); err != nil {
			return nil, err
		}
	}

	RcloneExecutablePath := c.RcloneExecutablePath

	if RcloneExecutablePath == "" {
		RcloneExecutablePath = "rclone"
	}

	internal := &InternalConfig{
		RecursiveDirs:        recursiveDirs,
		RcloneExecutablePath: RcloneExecutablePath,
	}

	return internal, nil
}

func (rc *RecursiveWatchedDir) setDefaultsAndValidate() error {
	if rc.Name == "" {
		return fmt.Errorf("name darf nicht leer sein")
	}

	if rc.Path == "" {
		return fmt.Errorf("path darf nicht leer sein")
	}

	if rc.FolderSize == nil {
		def := false
		rc.FolderSize = &def
	}
	if rc.LastModified == nil {
		def := false
		rc.LastModified = &def
	}
	return nil
}
