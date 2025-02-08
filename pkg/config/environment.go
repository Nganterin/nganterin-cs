package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

type Env struct {
	DB_USER            string
	DB_PASSWORD        string
	DB_HOST            string
	DB_PORT            string
	DB_NAME            string
	ENVIRONMENT        string
	PORT               string
	JWT_SECRET         string
	SMTP_SERVER        string
	SMTP_PORT          string
	SMTP_EMAIL         string
	SMTP_PASSWORD      string
	WEBCLIENT_BASE_URL string
	DASHBOARD_BASE_URL string
}

func InitEnvCheck() {
	environment := Env{
		DB_USER:            os.Getenv("DB_USER"),
		DB_PASSWORD:        os.Getenv("DB_PASSWORD"),
		DB_HOST:            os.Getenv("DB_HOST"),
		DB_PORT:            os.Getenv("DB_PORT"),
		DB_NAME:            os.Getenv("DB_NAME"),
		ENVIRONMENT:        os.Getenv("ENVIRONMENT"),
		PORT:               os.Getenv("PORT"),
		JWT_SECRET:         os.Getenv("JWT_SECRET"),
		SMTP_SERVER:        os.Getenv("SMTP_SERVER"),
		SMTP_PORT:          os.Getenv("SMTP_PORT"),
		SMTP_EMAIL:         os.Getenv("SMTP_EMAIL"),
		SMTP_PASSWORD:      os.Getenv("SMTP_PASSWORD"),
		WEBCLIENT_BASE_URL: os.Getenv("WEBCLIENT_BASE_URL"),
		DASHBOARD_BASE_URL: os.Getenv("DASHBOARD_BASE_URL"),
	}

	isEmpty, emptyFields := checkEmptyFields(environment)
	if isEmpty {
		log.Panicln("ERROR: The following environment variables are missing or empty:", emptyFields)
	} else {
		fmt.Println("INFO: All environment variables are set")
	}
}

func checkEmptyFields(env Env) (bool, []string) {
	v := reflect.ValueOf(env)
	typeOfEnv := v.Type()
	var emptyFields []string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			emptyFields = append(emptyFields, typeOfEnv.Field(i).Name)
		}
	}

	return len(emptyFields) > 0, emptyFields
}
