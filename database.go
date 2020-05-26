package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func initDB() (*sql.DB, error) {
	connectionStr := getDBEnv()
	db, err := sql.Open("mysql", connectionStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	return db, err
}

func addChannelToDB(channel string, nsfw bool, name string) error {
	var nsfwInt int

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	insert, err := db.Prepare("INSERT INTO channels (channelID, nsfw, name) VALUES (?, ?, ?)")

	defer insert.Close()

	if nsfw {
		nsfwInt = 1
	} else {
		nsfwInt = 0
	}

	_, err = insert.Exec(channel, nsfwInt, name)

	return err
}

func getChannelFromDB(channel string) (bool, string, error) {

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return false, "", err
	}

	output, err := db.Prepare("SELECT nsfw, name from channels WHERE channelID = ?")

	defer output.Close()

	if err != nil {
		return false, "", err
	}

	var nsfwInt int
	var name string

	err = output.QueryRow(channel).Scan(&nsfwInt, &name)

	return nsfwInt == 1, name, err
}

func setBannedSubreddit(channel string, subreddit string) error {
	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	insert, err := db.Prepare("INSERT INTO banned_subs (channelID, subreddit) VALUES (?, ?)")

	defer insert.Close()

	if err != nil {
		return err
	}

	_, err = insert.Exec(channel, subreddit)

	return err
}

func removeBannedSubreddit(channel string, subreddit string) error {
	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	remove, err := db.Prepare("DELETE FROM banned_subs WHERE channelID = ? AND subreddit = ?")

	defer remove.Close()

	_, err = remove.Exec(channel, subreddit)

	return err
}

// need a get all banned subs function
// also need functions related to the queue
