package config

import (
    "fmt"
    "os"
    "github.com/joho/godotenv"
)

func GetDBConfig() string {
    godotenv.Load()

    host := os.Getenv("DATABASE_HOST")
    user := os.Getenv("DATABASE_USER")
    password := os.Getenv("DATABASE_PASSWORD")
    dbname := os.Getenv("DATABASE_NAME")

    return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", 
        host, user, password, dbname)
}