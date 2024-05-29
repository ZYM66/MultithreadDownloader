package downloaderconfig

type MultiThreadConfig struct {
	HeaderConfig HeaderConfig
	Target       []string
	OutputPath   string
	NumChunk     int
}

func (c MultiThreadConfig) GetTarget() []string {
	return c.Target
}

func (c MultiThreadConfig) GetOutputPath() string {
	return c.OutputPath
}
