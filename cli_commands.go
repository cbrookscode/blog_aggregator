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

// Completely deletes all rows in users and DB table. FOR TESTING PURPOSES ONLY
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

func handlerUsers(s *state, cmd command) error {
	// Check for expected length of arguements
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("no arguments allowed for users command")
	}
	// Get all users from users table
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return fmt.Errorf("there are no users in the database yet")
	}

	// Print user names out in special format
	for i := 0; i < len(users); i++ {
		if users[i].Name == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", users[i].Name)
		} else {
			fmt.Printf("* %v\n", users[i].Name)
		}
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	// Check for expected length of arguements
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("need two arguements for add feed command- name, url")
	}
	name_string := cmd.arguments[0]
	url_string := cmd.arguments[1]

	// create feed
	new_feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name_string,
		Url:       url_string,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error with create feed: %w", err)
	}

	// create feed follow record for user addign feed
	_, err = s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    new_feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	fmt.Println(new_feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	// Check for expected length of arguements
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("no arguements needed for feeds command")
	}

	// Get feeds from db
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds from db: %w", err)
	}

	// Print out feeds name, url, and username that created the feed
	for i := 0; i < len(feeds); i++ {
		user, err := s.db.GetUserByID(context.Background(), feeds[i].UserID)
		if err != nil {
			return fmt.Errorf("error getting user by id: %w", err)
		}

		fmt.Printf("* %v\n", feeds[i].Name)
		fmt.Printf("* %v\n", feeds[i].Url)
		fmt.Printf("* %v\n", user.Name)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	// Check for expected length of arguements
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("need one arguement - url")
	}

	// Grab Feed info
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error getting feed by url: %w", err)
	}

	// Create Feed Follow Record
	feed_follow_info, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	// Notify of success
	fmt.Printf("* %v\n", feed_follow_info[0].FeedName)
	fmt.Printf("* %v\n", feed_follow_info[0].UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	// Check for expected length of arguements
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("no arguements needed for this function")
	}

	// Grab feed row info for given user
	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feed follow info for given user: %w", err)
	}

	// Print out username and feed follow names for that user
	fmt.Printf("%v is following the below feeds:\n", user.Name)
	for i := 0; i < len(feed_follows); i++ {
		fmt.Printf("* %v\n", feed_follows[i].FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	// Check for expected length of arguements
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("need one arguement - url")
	}

	// Grab Feed info
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error getting feed by url: %w", err)
	}

	err = s.db.DeleteFeedFollowRecordByUserFeedurlCombo(
		context.Background(),
		database.DeleteFeedFollowRecordByUserFeedurlComboParams{
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("user doesn't follow this feed")
	}

	fmt.Printf("%v has been unfollowed for %v\n", feed.Name, user.Name)

	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %w", err)
	}

	err = s.db.MarkFeedFetched(
		context.Background(),
		feed.ID,
	)
	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %w", err)
	}

	fetched_feed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	for i := 0; i < len(fetched_feed.Channel.Item); i++ {
		// convert the title into a nullable string type for db compatability
		title := fetched_feed.Channel.Item[i].Title
		titleNull := sql.NullString{String: title, Valid: title != ""}

		parsed_time, err := time.Parse(time.Layout, fetched_feed.Channel.Item[i].PubDate)
		if err != nil {
			return fmt.Errorf("error parsing published at string from post: %w", err)
		}
		timeNull := sql.NullTime{Time: parsed_time, Valid: parsed_time.IsZero()}

		_, err = s.db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       titleNull,
				Url:         fetched_feed.Channel.Item[i].Link,
				Description: fetched_feed.Channel.Item[i].Description,
				PublishedAt: timeNull,
				FeedID:      feed.ID,
			},
		)
		if err != nil {
			return fmt.Errorf("error creating post: %w", err)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	// Check for expected length of arguements
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("need one arguement - time duration")
	}

	// parse duration string into time duration value
	time_between_reqs := cmd.arguments[0]
	duration, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return fmt.Errorf("error parsing time duration string into duration value: %w", err)
	}

	// Setup a new ticker and logger that prints every interval of the duration value. Call scrapefeeds each time ticker ticks.
	fmt.Printf("Collecting feeds every %v\n", duration)
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerBrowse(s *state, cmd command) error {
	posts, err := s.db.GetPostsForUser(context.Background(), 5)
	if err != nil {
		return fmt.Errorf("error getting posts for user: %w", &err)
	}

	for i := 0; i < len(posts); i++ {
		fmt.Printf("* %v\n* %v\n", posts[i].Title, posts[i].Description)
	}

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return nil
		}
		return handler(s, cmd, user)
	}
}

func cli() (int, error) {
	// Read file to use for app state initialization, and print output of config file to get before snapshot
	my_config, err := config.Read()
	if err != nil {
		return 1, err
	}

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
	mycmds.register("users", handlerUsers)
	mycmds.register("agg", handlerAgg)
	mycmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	mycmds.register("feeds", handlerFeeds)
	mycmds.register("follow", middlewareLoggedIn(handlerFollow))
	mycmds.register("following", middlewareLoggedIn(handlerFollowing))
	mycmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	mycmds.register("browse", handlerBrowse)

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

	return 0, nil
}
