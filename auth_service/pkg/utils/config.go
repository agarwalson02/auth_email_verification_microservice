package utils

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "config/config-docker.yaml"
	}
	return "config/config-local.yaml"
}