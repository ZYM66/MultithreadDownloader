package downloaderconfig

type MultiThreadConfig struct {
	Target     string
	OutputPath string
	ChunkSize  int
}

func (c MultiThreadConfig) GetTarget() string {
	return c.Target
}

func (c MultiThreadConfig) GetOutputPath() string {
	return c.OutputPath
}
