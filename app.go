package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func main(){

	config := loadConfig("local", "./resources/")

	log.Println("starting server...")
	server := NewAPIServer(fmt.Sprintf(":%v", config.Server.Port))
	server.Run()
}

type Config struct{
	*Server
	*Database
	Jwt_secret_key string `yaml:"jwt_secret_key"`
}

type Server struct{
	Port string`yaml:"port"`
}

type Database struct{
	Host string `yaml:host`
	Port string`yaml:"port"`
	User string `yaml:"username"`
	Password string `yaml:"password"`
	Timeout uint64 `yaml:"timeout"`
	DB_name string `yaml:"db_name"`
	SSLmode string `yaml:"sslmode"`
}

func loadConfig(profile, path string) *Config{
	if profile == ""{
		profile = "local"
	}

	fileName := profile + ".yaml"

	log.Printf("application profile %v detected\n", profile)
	log.Printf("loading file %v", fileName)

	viper.SetConfigFile(path + fileName)

	if err := viper.ReadInConfig(); err != nil{
		log.Fatalf("Couldn`t read file %v. Error: %v", fileName, err)
	}

	log.Println("config file loaded successfully")
	
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil{
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	log.Printf("Server config: %+v\n", config.Server)
	log.Printf("Database config: %+v\n", config.Database)

	return config
}