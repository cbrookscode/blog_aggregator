package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func CreateConfigFile(url string) error {
	new_config := &Config{
		DbURL: url,
	}

	new_json, err := json.Marshal(new_config)
	if err != nil {
		return fmt.Errorf("error marshalling new config into json: %w", err)
	}

	loc, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error obtaining home dir for this pc: %w", err)
	}

	//create/overwrite existing file with blank slate
	file, err := os.Create(loc + "/.gatorconfig.json")
	if err != nil {
		return fmt.Errorf("error creating gatorconfig file in users home dir: %w", err)
	}

	//write to new blank file
	_, err = file.WriteString(string(new_json))
	if err != nil {
		return fmt.Errorf("error writing new json into gatorconfig file: %w", err)
	}

	return nil
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
	defer file.Close()

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

func (c *Config) SetUser(user string) error {
	// write conifg struct to json file after setting the current_user_name field-
	c.CurrentUserName = user
	new_json, err := json.Marshal(c)
	if err != nil {
		return err
	}

	loc, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	//create/overwrite existing file with blank slate
	file, err := os.Create(loc + "/.gatorconfig.json")
	if err != nil {
		return err
	}

	//write to new blank file
	_, err = file.WriteString(string(new_json))
	if err != nil {
		return err
	}

	return nil
}
