package bootstrap

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/joho/godotenv"
)

type Env struct {
	ServerAddress string
	Port          string
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
	DBPass        string
}

func NewEnv() *Env {
	_ = godotenv.Load()

	viper.AutomaticEnv()

	fmt.Printf("🔧 DB_HOST: %s\n", viper.GetString("DB_HOST"))

	fmt.Printf("🔧 DB_USER: %s\n", viper.GetString("DB_USER"))
	fmt.Printf("🔧 DB_PASS: %s\n", viper.GetString("DB_PASS"))
	fmt.Printf("🔧 DB_NAME: %s\n", viper.GetString("DB_NAME"))

	return &Env{
		ServerAddress: viper.GetString("SERVER_ADDRESS"),
		Port:          viper.GetString("PORT"),
		DBHost:        viper.GetString("DB_HOST"),
		DBPort:        viper.GetString("DB_PORT"),
		DBName:        viper.GetString("DB_NAME"),
		DBUser:        viper.GetString("DB_USER"),
		DBPass:        viper.GetString("DB_PASS"),
	}
}
