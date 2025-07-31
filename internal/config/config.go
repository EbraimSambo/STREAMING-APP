package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TranscodingConfig struct {
	Qualities    []Quality         `yaml:"qualities"`
	FFmpegParams []string          `yaml:"ffmpeg_params"`
}

type Quality struct {
	Name    string `yaml:"name"`
	Width   int    `yaml:"width"`
	Height  int    `yaml:"height"`
	Bitrate int    `yaml:"bitrate"`
}

type Config struct {
	Transcoding TranscodingConfig `yaml:"transcoding"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
