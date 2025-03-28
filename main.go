package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jubilant-gremlin/gator/internal/config"
	"github.com/jubilant-gremlin/gator/internal/database"

	_ "github.com/lib/pq"
)

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
	cmd_map.register("addfeed", handlerAddFeed)
	cmd_map.register("feeds", handlerFeeds)
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
