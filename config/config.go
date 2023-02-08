package config

type Config struct{
	GatewayEndpoint string // hostname to connect to
	Retries uint
};

const defaultGWEndpoint = "api.atmoperator.co.uk"
const defaultRetries = uint(3)

func GetConfig() Config {
	return Config{
	// we could programmatically process a config file (YAML or JSON)
		GatewayEndpoint: defaultGWEndpoint,
		Retries: defaultRetries,
	}
}
