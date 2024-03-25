package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go"
)

func LoadEnv() {
	err := godotenv.Load()
	stripe.Key = os.Getenv("STRIPE_KEY")

	if err != nil {
		log.Fatal("loading error: ", err)
	}
}
