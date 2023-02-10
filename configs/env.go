package configs

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type envVars struct {
	ClientUrl                           string   `env:"CLIENT_URL" validate:"url"`
	BackendCorsOrigins                  []string `env:"BACKEND_CORS_ORIGINS" validate:"dive,url"`
	ProjectName                         string   `env:"PROJECT_NAME"`
	SecretKey                           string   `env:"SECRET_KEY" validate:"gte=8"`
	AccessTokenExpirationMinutes        int      `env:"ACCESS_TOKEN_EXPIRATION_MINUTES" validate:"numeric"`
	PasswordResetTokenExpirationMinutes int      `env:"PASSWORD_RESET_TOKEN_EXPIRATION_MINUTES" validate:"numeric"`
	UsersOpenRegistration               bool     `env:"USERS_OPEN_REGISTRATION" validate:"boolean"`
	FirstSuperuser                      string   `env:"FIRST_SUPERUSER" validate:"email"`
	FirstSuperuserPassword              string   `env:"FIRST_SUPERUSER_PASSWORD" validate:"gte=8"`
	EmailsEnabled                       bool     `env:"EMAILS_ENABLED" validate:"boolean"`
	EmailsApiKey                        string   `env:"EMAILS_API_KEY"`
	DBUser                              string   `env:"DB_USER"`
	DBPassword                          string   `env:"DB_PASSWORD"`
	DBHost                              string   `env:"DB_HOST" validate:"hostname"`
	DBPort                              string   `env:"DB_PORT" validate:"numeric"`
	DBName                              string   `env:"DB_NAME"`
}

func parseEnv(vars *envVars, filenames ...string) error {
	if err := godotenv.Load(filenames...); err != nil {
		return err
	}

	t := reflect.TypeOf(*vars)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get("env")
		envValue := os.Getenv(tagValue)

		if envValue == "" {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			reflect.ValueOf(vars).Elem().FieldByName(field.Name).Set(reflect.ValueOf(envValue))
		case reflect.Int:
			value, err := strconv.Atoi(envValue)
			if err != nil {
				return err
			}
			reflect.ValueOf(vars).Elem().FieldByName(field.Name).Set(reflect.ValueOf(value))
		case reflect.Bool:
			value, err := strconv.ParseBool(envValue)
			if err != nil {
				return err
			}
			reflect.ValueOf(vars).Elem().FieldByName(field.Name).Set(reflect.ValueOf(value))
		case reflect.Slice:
			values := []string{}
			if err := json.Unmarshal([]byte(envValue), &values); err != nil {
				return err
			}
			reflect.ValueOf(vars).Elem().FieldByName(field.Name).Set(reflect.ValueOf(values))
		}
	}

	if err := validator.New().Struct(vars); err != nil {
		return err
	}

	return nil
}

func getEnvVars() envVars {
	vars := envVars{
		// 60 minutes * 24 hours * 8 days = 8 days
		AccessTokenExpirationMinutes:        60 * 24 * 8,
		PasswordResetTokenExpirationMinutes: 15,
	}

	err := parseEnv(&vars, ".env")
	if err != nil {
		log.Fatal(err)
	}

	return vars
}

var Env = getEnvVars()
