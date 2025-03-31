package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
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

func scrapeFeeds(s *state) {
	next, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("ERROR GETTING NEXT FEED TO FETCH: %v\n", err)
		return
	}
	fmt.Printf("FETCHING FEED: %v\n", next.Url)
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true}, ID: next.ID})
	if err != nil {
		fmt.Printf("ERROR MARKING FEED AS FETCHED: %v\n", err)
		return
	}
	fetched, err := fetchFeed(context.Background(), next.Url)
	if err != nil {
		fmt.Printf("ERROR FETCHING FEED: %v\n", err)
		return
	}
	fmt.Printf("FEED: %v - %d posts found\n", next.Name, len(fetched.Channel.Item))
	for i := range fetched.Channel.Item {
		pub, err := time.Parse(time.RFC1123Z, fetched.Channel.Item[i].PubDate)
		if err != nil {
			fmt.Printf("ERROR PARSING DATE:%v\n", err)
		}
		err = s.db.CreatePost(context.Background(), database.CreatePostParams{CreatedAt: time.Now(), UpdatedAt: time.Now(), Title: fetched.Channel.Item[i].Title, Url: fetched.Channel.Item[i].Link, Description: sql.NullString{String: fetched.Channel.Item[i].Description, Valid: true}, PublishedAt: sql.NullTime{Time: pub, Valid: true}, FeedID: sql.NullInt64{Int64: next.ID, Valid: true}})
		if err != nil {
			fmt.Printf("ERROR CREATING POST: %v\n", err)
			return
		}
	}
	fmt.Println()
}

func handlerBrowse(s *state, cmd command) error {
	switch len(cmd.arguments) {
	default:
		limit, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			fmt.Printf("ERROR GETTING POST LIMIT: %v\n", err)
			return err
		}
		posts, err := s.db.GetPostsForUser(context.Background(), int32(limit))
		if err != nil {
			fmt.Printf("ERROR GETTING POSTS: %v\n", err)
			return err
		}
		for i := range posts {
			fmt.Println(posts[i].PublishedAt.Time)
			fmt.Printf("%v:\n", posts[i].Title)
			fmt.Println(posts[i].Description.String)
		}
		return nil
	case 0:
		posts, err := s.db.GetPostsForUser(context.Background(), 2)
		if err != nil {
			fmt.Printf("ERROR GETTING POSTS: %v\n", err)
			return err
		}
		for i := range posts {
			fmt.Println(posts[i].PublishedAt.Time)
			fmt.Printf("%v:\n", posts[i].Title)
			fmt.Println(posts[i].Description.String)
		}
		return nil
	}
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	user_id := user.ID
	if len(cmd.arguments) == 0 {
		return errors.New("ERROR: UNFOLLOW COMMAND MUST SPECIFY FEED URL")
	}
	url := cmd.arguments[0]
	feed_to_delete, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		fmt.Printf("ERROR GETTING FEED: %v", err)
		return err
	}
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{UserID: uuid.NullUUID{UUID: user_id, Valid: true}, FeedID: sql.NullInt64{Int64: feed_to_delete.ID, Valid: true}})
	if err != nil {
		fmt.Printf("ERROR DELETING FEED FROM FOLLOWS: %v\n", err)
		return err
	}
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	user_id := user.ID
	following, err := s.db.GetFeedFollowsForUser(context.Background(), uuid.NullUUID{UUID: user_id, Valid: true})
	if err != nil {
		fmt.Printf("ERROR GETTING FOLLOWS FOR USER: %v", err)
		return err
	}
	fmt.Printf("%v IS FOLLOWING:\n", user.Name)
	for i := range following {
		fmt.Printf("- %v\n", following[i].FeedName)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return errors.New("ERROR: FOLLOW COMMAND MUST SPECIFY FEED URL")
	}
	url := cmd.arguments[0]
	user_id := user.ID
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		fmt.Printf("ERROR GETTING FEED: %v\n", err)
		return err
	}
	feed_id := feed.ID
	followed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: uuid.NullUUID{UUID: user_id, Valid: true}, FeedID: sql.NullInt64{Int64: feed_id, Valid: true}})
	if err != nil {
		fmt.Printf("ERROR FOLLOWING FEED:%v\n", err)
		return err
	}
	for i := range followed {
		fmt.Printf("SUCCESS: %v FOLLOWED %v", followed[i].UserName, followed[i].FeedName)
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Printf("ERROR GETTING FEEDS: %v\n", err)
		return err
	}
	for i := range feeds {
		user, err := s.db.GetUserName(context.Background(), (feeds[i].UserID.UUID))
		if err != nil {
			fmt.Printf("ERROR GETTING USER NAME: %v\n", err)
			return err
		}
		fmt.Printf("Feed Name: %v, URL: %v, User Name: %v\n", feeds[i].Name, feeds[i].Url, user)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return errors.New("ERROR: ADD FEED MUST SPECIFY BOTH NAME AND URL OF FEED")
	}
	user_id := user.ID
	entry, err := s.db.CreateFeedEntry(context.Background(), database.CreateFeedEntryParams{Name: cmd.arguments[0], Url: cmd.arguments[1], CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: uuid.NullUUID{UUID: user_id, Valid: true}})
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return err
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: uuid.NullUUID{UUID: user_id, Valid: true}, FeedID: sql.NullInt64{Int64: entry.ID, Valid: true}})
	if err != nil {
		fmt.Printf("ERROR FOLLOWING FEED:%v\n", err)
		return err
	}

	fmt.Printf("SUCCESS! %v WAS ADDED FOR %v\n", entry.Name, user.Name)
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("ERROR: AGG COMMAND MUST SPECIFY TIME DURATION")
	}
	time_between_reqs, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		fmt.Printf("ERROR PARSING TIME: %v\n", err)
		return err
	}
	fmt.Printf("COLLECTING FEEDS EVERY %v\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
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
		return err
	}
	err = s.db.ResetFeed(context.Background())
	if err != nil {
		fmt.Printf("ERROR:%v\n", err)
		return err
	}
	err = s.db.ResetPosts(context.Background())
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return err
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
