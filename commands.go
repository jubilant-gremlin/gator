package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jubilant-gremlin/gator/internal/config"
	"github.com/jubilant-gremlin/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Println("ERROR GETTING FEEDS")
		return err
	}
	for i := range feeds {
		user, err := s.db.GetUserName(context.Background(), (feeds[i].UserID.UUID))
		if err != nil {
			fmt.Println("ERROR GETTING USER NAME")
			return err
		}
		fmt.Printf("Feed Name: %v, URL: %v, User Name: %v\n", feeds[i].Name, feeds[i].Url, user)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) < 2 {
		fmt.Println("ERROR: ADD FEED MUST SPECIFY BOTH NAME AND URL OF FEED")
		os.Exit(1)
	}
	current_user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	user_id := current_user.ID
	if err != nil {
		fmt.Println("ERROR GETTING CURRENT USER")
	}
	entry, err := s.db.CreateFeedEntry(context.Background(), database.CreateFeedEntryParams{Name: cmd.arguments[0], Url: cmd.arguments[1], CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: uuid.NullUUID{UUID: user_id, Valid: true}})
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
	fmt.Println(entry)
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Println("ERROR FETCHING FEED")
		return err
	}
	fmt.Println(*feed)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Println("ERROR GETTING USERS")
		os.Exit(1)
	}
	for _, user := range users {
		if user == s.cfg.Current_user_name {
			fmt.Printf("* %v (current)", user)
			continue
		}
		fmt.Printf("* %v\n", user)
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("ERROR: must have user name to log in")
	}
	name := cmd.arguments[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		fmt.Println("ERROR: user does not exist")
		os.Exit(1)

	}
	s.cfg.SetUser(name)
	fmt.Println("SUCCESS! USER SET")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("ERROR: user must have name")
	}
	name := cmd.arguments[0]
	_, err := s.db.GetUser(context.Background(), name)
	// if user is not in db
	if err != nil {
		user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name})
		if err != nil {
			return err
		}
		s.cfg.SetUser(name)
		fmt.Printf("USER CREATED:%v\n", user)
	} else {
		// if user is in db
		fmt.Println("ERROR: user already in system")
		os.Exit(1)
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		fmt.Printf("ERROR:%v\n", err)
	}
	err = s.db.ResetFeed(context.Background())
	if err != nil {
		fmt.Printf("ERROR:%v\n", err)
	}
	fmt.Println("DATABASE RESET SUCCESSFUL")
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
	_, ok := c.cmds[name]
	if !ok {
		fmt.Println("ERROR REGISTERING COMMAND")
	}
}

func (c *commands) run(s *state, cmd command) error {
	handlerName := cmd.name
	handler, ok := c.cmds[handlerName]
	if !ok {
		return errors.New("ERROR: command not found")
	}
	handler(s, cmd)
	return nil
}
