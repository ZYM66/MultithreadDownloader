package config

type Config interface {
	GetTarget() string
	GetOutputPath() string
}

type BaseConfig struct {
	Target     string
	OutputPath string
}

func (c *BaseConfig) GetTarget() string {
	return c.Target
}

func (c *BaseConfig) GetOutputPath() string {
	return c.OutputPath
}

type MultiThreadConfig struct {
	BaseConfig
	ChunkSize int
}
