package internal

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	// reads json gator config file and returns Config struct. should read file from home directory then decode json string in a new config struct. os.UserHomeDir to get location of HOME.
	loc, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(loc + "/.gatorconfig.json")
	if err != nil {
		return Config{}, err
	}

	mybytes, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(mybytes, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
