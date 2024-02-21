package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"regexp"
	"strings"
	//	"log"
	//	s "main/server"
	//	"net/http"
	"os"
)

/*
Make PostgresPlayerStore
I am thinking we will need to store a connection? Done for now, look at more example of how to manage database connections in code similar to this
How many connections do I want to be able to support? for now max out at 10, doesn't really matter atm.
I have a lot of questions to ask myself and there is a lot to implement here. 1. once db calls are done make sure the website handles concurrent users

*/

var db *pgxpool.Pool

func startDB() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("WEB_APP_DB"))
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

type InMemoryPlayerStore struct {
	store map[string]int
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

func sanitizeInput(input string) string {
	// Allow only alphabetical characters (both uppercase and lowercase)
	re := regexp.MustCompile("[^a-zA-Z]")
	return re.ReplaceAllString(input, "")
}

func main() {

	//	myServer := &s.PlayerServer{Store: NewInMemoryPlayerStore()}
	//	log.Fatal(http.ListenAndServe(":5000", myServer))

	db, err := startDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rdr := bufio.NewReader(os.Stdin)
	userInput, _ := rdr.ReadString('\n')
	userInput = strings.TrimSuffix(userInput, "\n")
	cleanedUserInput := sanitizeInput(userInput)
	fmt.Println(cleanedUserInput)
	queryString := fmt.Sprintf("SELECT gamesWon FROM scores WHERE name = '%s';", cleanedUserInput)

	ctx := context.Background()

	rows, err := db.Query(ctx, queryString)
	numbers, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (int32, error) {
		var n int32
		err := row.Scan(&n)
		return n, err
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(numbers)
}
