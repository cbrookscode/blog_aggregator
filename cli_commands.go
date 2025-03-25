package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	config "github.com/cbrookscode/blog_aggregator/internal/config"
	"github.com/cbrookscode/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

// DB URL - "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

// Map of command names to their handler functions.
type commands struct {
	cmdnames map[string]func(*state, command) error
}

// Registers a new handler function for a command name.
func (c *commands) register(name string, f func(*state, command) error) {
	if _, exists := c.cmdnames[name]; exists {
		fmt.Println("The handler you are trying to register already exists")
		return
	}

	c.cmdnames[name] = f
}

// Runs a given command with the provided state if it exists.
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

// Logs on a user which simply means adjusting the config file with the users name. Will only be done if user has been registered
func handlerLogin(s *state, cmd command) error {
	// check for expected length of arguments
	if len(cmd.arguments) == 0 || len(cmd.arguments) > 1 {
		return fmt.Errorf("can only provide one string to represent login name. please try again. user was not registered for login")
	}
	username := cmd.arguments[0]

	// ensure user has been registerd in DB already
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println("user attempting to logon has not be registered")
		os.Exit(1)
	}

	// edit config file so that it shows the new user logging on
	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	// Inform user of success
	fmt.Println("User has been set")
	return nil
}

// Register a new user in the DB.
func handlerRegister(s *state, cmd command) error {
	// Check length of expected inputs
	if len(cmd.arguments) == 0 || len(cmd.arguments) > 1 {
		return fmt.Errorf("can only provide one string to represent user. please try again. user was not registered")
	}
	new_username := cmd.arguments[0]

	// Check if user exists already
	_, err := s.db.GetUser(context.Background(), new_username)
	if err == nil {
		return fmt.Errorf("user is already registered")
	}

	// Create new user in DB
	new_user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: new_username})
	if err != nil {
		return err
	}

	// Log user on by changing config file
	err = s.cfg.SetUser(new_username)
	if err != nil {
		return err
	}

	// Prompt user of success and show new user info
	fmt.Println("User has been registered")
	fmt.Printf("%v\n", new_user)
	return nil
}

// Completely deletes all rows in users DB table. FOR TESTING PURPOSES ONLY
func handlerReset(s *state, cmd command) error {
	// Check for expected length of arguements
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("no arguments allowed for reset command")
	}

	// Delete all rows in users DB table
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	// Prompt of success
	fmt.Println("Users table has been deleted")
	return nil
}

func cli() (int, error) {
	// Read file to use for app state initialization, and print output of config file to get before snapshot
	my_config, err := config.Read()
	if err != nil {
		return 1, err
	}
	fmt.Println(my_config)

	// Open connection to the DB and intialize app_state
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable")
	if err != nil {
		return 1, err
	}
	defer db.Close()
	dbQueries := database.New(db)
	app_state := state{dbQueries, &my_config}

	// initialize map of commands, and struct that holds map. Register commands
	cmds_map := make(map[string]func(*state, command) error)
	mycmds := commands{cmds_map}
	mycmds.register("login", handlerLogin)
	mycmds.register("register", handlerRegister)
	mycmds.register("reset", handlerReset)

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
	my_config, err = config.Read()
	if err != nil {
		return 1, err
	}
	fmt.Println(my_config)

	return 0, nil
}
