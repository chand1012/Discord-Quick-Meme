package main

type agentFile struct {
	UserAgent    string `yaml:"user_agent"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

func generateAgentFile() error {

}
