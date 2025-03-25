package initializer

import(
	"fmt"
	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	fmt.Println("Loading Environment Variables...")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	fmt.Println("Environment Variables Loaded Successfully")
}
