package downloaderconfig

type HeaderConfig struct {
	UA    string
	Range string
}

func NewHeaderConfig() HeaderConfig {
	return HeaderConfig{}
}
