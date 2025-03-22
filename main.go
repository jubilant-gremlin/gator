package main

import (
	"fmt"
	"os"
	"github.com/jubilant-gremlin/gator/internal/config"
)

func main() {
	foo := config.Read()
	s := state{
	cfg: &foo,
	}
	bar := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	bar.register("login", handlerLogin)
	args := os.Args
	if len(args) < 3 {
		fmt.Println("ERROR: not enough arguments")
	}
	new_cmd := command {
		name: args[1],
		arguments: args[2:],
	}
	err := bar.run(&s, new_cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
