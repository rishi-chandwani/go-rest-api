package configs

import (
	"os"

	"github.com/joho/godotenv"
)

func loadEnvVariables() {
	envLoadErr := godotenv.Load()
	if envLoadErr != nil {
		panic(envLoadErr)
	}
}

// GetEnvMongoLink function loads the environment variables mentioned in .env file.
//
// While loading variables informtion if there was any error it will panic and stop execution from there.
// Otherwise, it will fetch MongoDB Connection URL and return to the caller.
func GetEnvMongoLink() string {
	loadEnvVariables()

	return os.Getenv("MONGOURI")
}

func GetApplicationPort() string {
	return os.Getenv("PORT")
}
