package config

type DownloaderConfig interface {
	GetTarget() string
	GetOutputPath() string
}
