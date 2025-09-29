package env

import (
	"log"
	"os"

	"github.com/fox998/Chirpy/internal/config/env/platform"
	"github.com/joho/godotenv"
)

type EnvParams struct {
	DbUrl    string
	Platform platform.Type
	Secret   string
}

const (
	DB_URL   string = "DB_URL"
	PLATFORM string = "PLATFORM"
	SECRET   string = "SECRET"
)

func GetEnvParams() EnvParams {
	godotenv.Load()

	params := EnvParams{
		DbUrl:    os.Getenv(DB_URL),
		Platform: platform.GetPlatformType(os.Getenv(PLATFORM)),
		Secret:   os.Getenv(SECRET),
	}

	if params.DbUrl == "" {
		log.Fatalf("Env variable %s is empty\n", DB_URL)
	}

	if params.Secret == "" {
		log.Fatalf("Env variable %s is empty\n", SECRET)
	}

	return params
}
