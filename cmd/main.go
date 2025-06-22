package main

import (
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/tracey-io/tracey/api"
	"os"
)

func main() {
	log.SetLevel(log.DebugLevel)

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")
	appEnv := os.Getenv("ENVIRONMENT")

	secretKey := os.Getenv("SECRET_KEY")

	var requiredEnvironmentVariables = map[string]string{
		"SERVER_HOST":     serverHost,
		"SERVER_PORT":     serverPort,
		"APP_ENVIRONMENT": appEnv,
		"SECRET_KEY":      secretKey,
	}

	envValidator(requiredEnvironmentVariables, func(missingVariables []string) {
		log.Fatal("Missing environment variables.", "variables", missingVariables)
	})

	serverConfig := &api.ServerConfig{
		Address: &api.ServerAddress{
			Host: serverHost,
			Port: serverPort,
		},
		Environment: api.Environment(appEnv),
	}

	if err := api.StartServer(serverConfig); err != nil {
		log.Fatal(err)
	}
}

func envValidator(variables map[string]string, fn func(missingVariables []string)) {
	var missingVariables []string

	for name, value := range variables {
		if value == "" {
			missingVariables = append(missingVariables, name)
		}
	}

	if len(missingVariables) > 0 {
		fn(missingVariables)
	}
}
