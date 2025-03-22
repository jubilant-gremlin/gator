package main

import (
	"errors"
	"fmt"
	"github.com/jubilant-gremlin/gator/internal/config"
)

type state struct {
	cfg *config.Config 
}

type command struct {
	name string
	arguments []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("ERROR: login expects username arg")
	}
	fmt.Println(cmd.arguments)
	s.cfg.SetUser(cmd.arguments[0])
	fmt.Println("SUCCESS! USER SET")
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
	_, ok := c.cmds[name]
	if !ok {
		fmt.Println("ERROR REGISTERING COMMAND")
	}
}

func(c*commands) run(s *state, cmd command) error {
	handlerName := cmd.name
	handler, ok := c.cmds[handlerName]
	if !ok {
		return errors.New("ERROR: command not found")
	}
	handler(s, cmd)
	return nil
}

