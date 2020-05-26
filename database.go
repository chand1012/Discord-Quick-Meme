package main

import (
	"database/sql"
	"strings"

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

// AddChannelToDB adds a channel to the database
func AddChannelToDB(channel string, nsfw bool, name string) error {
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

// GetChannelFromDB gets the channel from the database
func GetChannelFromDB(channel string) (bool, string, error) {

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

// add a function that sets the last time that the channel made a request
// If the channel hasn't had a meme sent to it in a month, delete its records

// SetBannedSubreddit adds a sub to the channel bans
func SetBannedSubreddit(channel string, subreddit string) error {
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

// RemoveBannedSubreddit removes a sub from the bans list
func RemoveBannedSubreddit(channel string, subreddit string) error {
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

// GetAllBannedSubs gets all of the banned subreddits
func GetAllBannedSubs(channel string) ([]string, error) {
	db, err := initDB()

	defer db.Close()

	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT subreddit FROM banned_subs WHERE channelID = ?", channel)

	if err != nil {
		return nil, err
	}

	var subs []string
	var sub string

	for rows.Next() {
		err = rows.Scan(&sub)
		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}

	err = rows.Err()

	return subs, err
}

// GetMemeQueue gets data for the meme queue for the specified channel
func GetMemeQueue(channel string) (QueueObj, error) {
	var queue QueueObj
	var nsfwInt int
	var subString string

	db, err := initDB()

	db.Close()

	if err != nil {
		return QueueObj{}, err
	}

	row := db.QueryRow("SELECT interval, subreddits, nsfw, time FROM queue WHERE channelID = ?", channel)
	err = row.Scan(&queue.Interval, &subString, &nsfwInt, &queue.Time)

	if err != nil {
		return QueueObj{}, err
	}

	queue.SubReddits = strings.Split(subString, ",")

	if nsfwInt == 1 {
		queue.NSFW = true
	} else {
		queue.NSFW = false
	}

	return QueueObj{}, nil

}

// DeleteMemeQueue clears the meme queue for the specified channel
func DeleteMemeQueue(channel string) error {
	db, err := initDB()

	db.Close()

	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM queue WHERE channelID = ?", channel)

	return err
}

// SetMemeQueue sets the meme queue for the channel
func SetMemeQueue(channel string, nsfw bool, interval string, subs string) error {
	db, err := initDB()

	db.Close()

	if err != nil {
		return err
	}

	insert, err := db.Prepare("INSERT INTO queue (channelID, nsfw, interval, subreddits) VALUES (?, ?, ?, ?)")

	if err != nil {
		return err
	}

	var nsfwInt int

	if nsfw {
		nsfwInt = 1
	} else {
		nsfwInt = 0
	}

	_, err = insert.Exec(channel, nsfwInt, interval, subs)

	return err
}
