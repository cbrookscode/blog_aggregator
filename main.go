package main

import (
	"fmt"
	"os"

	config "github.com/cbrookscode/blog_aggregator/internal/config"
	_ "github.com/lib/pq"
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
}

// This method runs a given command with the provided state if it exists.
func (c *commands) run(s *state, cmd command) error {
	cmd_func, exists := c.cmdnames[cmd.name]
	if !exists {
		return fmt.Errorf("command doesn't exist")
	}
	err := cmd_func(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 || len(cmd.arguments) > 1 {
		return fmt.Errorf("can only provide one string to represent login name. please try again. user was not registered for login")
	}

	err := s.cfg.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set")
	return nil
}

func cli() (int, error) {
	// Read file to use for app state initialization, and print output of config file to get before snapshot
	test_struct, err := config.Read()
	if err != nil {
		return 1, err
	}
	fmt.Println(test_struct)

	// initialize app state, map of commands, and struct that holds map. Register login command
	app_state := state{&test_struct}
	cmds := make(map[string]func(*state, command) error)
	mycmds := commands{cmds}
	mycmds.register("login", handlerLogin)

	// build command struct based on inputs from user when running program. first arg is always program name, second is assumed to be command name, rest are arguements for command
	cmd := command{}
	args := os.Args
	if len(args) < 2 {
		return 1, fmt.Errorf("need to provide command when running program")
	} else if len(args) > 2 {
		cmd.name = args[1]
		cmd.arguments = args[2:]
	} else {
		cmd.name = args[1]
	}

	// run command
	err = mycmds.run(&app_state, cmd)
	if err != nil {
		return 1, err
	}

	// test output to see if config file was changed
	test_struct, err = config.Read()
	if err != nil {
		return 1, err
	}
	fmt.Println(test_struct)
	return 0, nil
}

func main() {
	err_val, err := cli()
	if err != nil {
		fmt.Println(err)
		os.Exit(err_val)
	}
}
