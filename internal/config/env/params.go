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
	PolkaKey string
}

const (
	DB_URL    string = "DB_URL"
	PLATFORM  string = "PLATFORM"
	SECRET    string = "SECRET"
	POLKA_KEY string = "POLKA_KEY"
)

func GetEnvParams() EnvParams {
	godotenv.Load()

	params := EnvParams{
		DbUrl:    os.Getenv(DB_URL),
		Platform: platform.GetPlatformType(os.Getenv(PLATFORM)),
		Secret:   os.Getenv(SECRET),
		PolkaKey: os.Getenv(POLKA_KEY),
	}

	if params.DbUrl == "" {
		log.Fatalf("Env variable %s is empty\n", DB_URL)
	}

	if params.Secret == "" {
		log.Fatalf("Env variable %s is empty\n", SECRET)
	}

	if params.PolkaKey == "" {
		log.Fatalf("Env variable %s is empty\n", POLKA_KEY)
	}

	return params
}
