package main

import (
	"fmt"
	"os"

	config "github.com/cbrookscode/blog_aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	// This will be a map of command names to their handler functions.
	cmdnames map[string]func(*state, command) error
}

// This method registers a new handler function for a command name.
func (c *commands) register(name string, f func(*state, command) error) {
	if _, exists := c.cmdnames[name]; exists {
		fmt.Println("The handler you are trying to register already exists")
	}

	c.cmdnames[name] = f
	return
}

// This method runs a given command with the provided state if it exists.
func (c *commands) run(s *state, cmd command) error {
	cmd_func, exists := c.cmdnames[cmd.name]
	if !exists {
		return fmt.Errorf("command doesn't exist")
	}
	cmd_func(s, cmd)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 || len(cmd.arguments) > 1 {
		return fmt.Errorf("only one string input is valid")
	}

	err := s.cfg.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set")
	return nil
}

func main() {
	test_struct, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(test_struct)

	app_state := state{&test_struct}
	var cmds map[string]func(*state, command) error
	mycmds := commands{cmds}
	mycmds.register("login", handlerLogin)

	arguments := os.Args

	err = test_struct.SetUser("Hilda Brown")
	if err != nil {
		fmt.Println(err)
		return
	}

	test_struct, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(test_struct)
}
