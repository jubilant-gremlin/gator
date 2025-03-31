package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jubilant-gremlin/gator/internal/config"
	"github.com/jubilant-gremlin/gator/internal/database"

	_ "github.com/lib/pq"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		current_user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
		if err != nil {
			fmt.Printf("ERROR GETTING USER: %v", err)
			return err
		}
		err = handler(s, cmd, current_user)
		if err != nil {
			return err
		}
		return nil
	}

}

func main() {
	cfg := config.Read()
	// load db url to config struct and sql.Open() a connection to db
	db, err := sql.Open("postgres", cfg.Db_url)
	dbQueries := database.New(db)
	s := state{
		db:  dbQueries,
		cfg: &cfg,
	}
	// initialize command map
	cmd_map := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	// register commands
	cmd_map.register("login", handlerLogin)
	cmd_map.register("register", handlerRegister)
	cmd_map.register("reset", handlerReset)
	cmd_map.register("users", handlerUsers)
	cmd_map.register("agg", handlerAgg)
	cmd_map.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmd_map.register("feeds", handlerFeeds)
	cmd_map.register("follow", middlewareLoggedIn(handlerFollow))
	cmd_map.register("following", middlewareLoggedIn(handlerFollowing))
	cmd_map.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmd_map.register("browse", handlerBrowse)
	// interpret cli args to command
	args := os.Args
	new_cmd := command{
		name:      args[1],
		arguments: args[2:],
	}
	// run commands
	err = cmd_map.run(&s, new_cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
