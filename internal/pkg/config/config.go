package config

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
)

type Config struct {
	MAX_MEMORY   int64  `type:"required"`
	PARSE_PERIOD int64  `type:"required"`
	WORK_LIMIT   int64  `type:"required"`
	MAX_ROUTINES int    `type:"required"`
	DB_NAME      string `type:"required"`
	DB_PASSWORD  string `type:"required"`
	DB_USER      string `type:"required"`
	DB_HOST      string `type:"required"`
	GRPC_HOST    string `type:"required"`
}

func (conf *Config) parseEnv() {
	var projectDir string

	if flag.Lookup("source") != nil {
		flag.StringVar(&projectDir, "source", "", "Provide project .go files absolute path")
		flag.Parse()
	}

	if projectDir == "" {
		ex, err := os.Executable()
		if err != nil {
			ex, _ = filepath.EvalSymlinks(ex)
		}
		projectDir = filepath.Dir(ex) + "/../"
	}

	if err := godotenv.Load(projectDir + ".env"); err != nil {
		fmt.Println("File .env not found, reading configuration from ENV", err)
	}

	val := reflect.ValueOf(conf).Elem()
	for i := 0; i < val.NumField(); i++ {
		var (
			envVal string
			name   string
		)
		name = val.Type().Field(i).Name
		envVal = os.Getenv(name)

		if tag := val.Type().Field(i).Tag.Get("type"); envVal == "" && tag == "required" {
			panic("Required config value does not set!")
		}

		switch val.Field(i).Kind() {
		case reflect.Int64, reflect.Int:
			if value, err := strconv.ParseInt(envVal, 10, 64); err == nil {
				val.Field(i).SetInt(value)
			}
		case reflect.String:
			val.Field(i).SetString(envVal)
		}

	}
}

func New() *Config {
	config := new(Config)
	config.parseEnv()
	return config
}
