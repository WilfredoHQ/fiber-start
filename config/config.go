package config

import (
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type config struct {
	AccessTokenExpirationMinutes        int
	PasswordResetTokenExpirationMinutes int
	ClientUrl                           string `env:"CLIENT_URL" validate:"omitempty,url"`
	BackendCorsOrigins                  string `env:"BACKEND_CORS_ORIGINS" validate:"required"`
	ProjectName                         string `env:"PROJECT_NAME"`
	SecretKey                           string `env:"SECRET_KEY" validate:"required"`
	UsersOpenRegistration               bool   `env:"USERS_OPEN_REGISTRATION"`
	FirstSuperuser                      string `env:"FIRST_SUPERUSER" validate:"required,email"`
	FirstSuperuserPassword              string `env:"FIRST_SUPERUSER_PASSWORD" validate:"required,min=8"`
	EmailsEnabled                       bool   `env:"EMAILS_ENABLED"`
	EmailsApiKey                        string `env:"EMAILS_API_KEY"`
	DBUser                              string `env:"DB_USER" validate:"required"`
	DBPassword                          string `env:"DB_PASSWORD" validate:"required"`
	DBHost                              string `env:"DB_HOST" validate:"required,hostname"`
	DBPort                              int    `env:"DB_PORT" validate:"required"`
	DBName                              string `env:"DB_NAME" validate:"required"`
}

func parseConfig(conf *config) error {
	t := reflect.TypeOf(*conf)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get("env")
		envValue := os.Getenv(tagValue)

		if envValue == "" {
			continue
		}

		setValue := func(value any) {
			reflect.ValueOf(conf).Elem().FieldByName(field.Name).Set(reflect.ValueOf(value))
		}

		switch field.Type.Kind() {
		case reflect.String:
			setValue(envValue)
		case reflect.Int:
			value, err := strconv.Atoi(envValue)
			if err != nil {
				return err
			}
			setValue(value)
		case reflect.Bool:
			value, err := strconv.ParseBool(envValue)
			if err != nil {
				return err
			}
			setValue(value)
		}
	}

	return nil
}

func getConfig() config {
	conf := config{
		// 60 minutes * 24 hours * 8 days = 8 days
		AccessTokenExpirationMinutes:        60 * 24 * 8,
		PasswordResetTokenExpirationMinutes: 15,
	}

	err := parseConfig(&conf)
	if err != nil {
		log.Fatal(err)
	}

	validate := validator.New()
	if err := validate.Struct(&conf); err != nil {
		log.Fatal(err)
	}

	return conf
}

var Config = getConfig()
