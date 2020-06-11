package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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
	fmt.Println("Adding channel to database: " + channel)
	var nsfwInt int

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	if nsfw {
		nsfwInt = 1
	} else {
		nsfwInt = 0
	}

	// gets the channel from the database to see if it exists.
	output, err := db.Prepare("SELECT channelID from channels WHERE channelID = ?")
	defer output.Close()

	if err != nil {
		return err
	}

	err = output.QueryRow(channel).Scan(&channel)

	if err == nil { // if there is no error that means that the entry must be updated. Ideally this won't happen often.
		fmt.Println("Already in database, updating existing record...")
		chanTime := time.Now().Unix()
		_, err = db.Exec("UPDATE channels SET name = ?, nsfw = ?, time = ? WHERE channelID = ?", name, nsfwInt, chanTime, channel)

	} else if err == sql.ErrNoRows {
		fmt.Println("Adding new entry...")
		insert, err := db.Prepare("INSERT INTO channels (channelID, nsfw, name, time) VALUES (?, ?, ?, ?)")

		if err != nil {
			fmt.Println(err)
			return err
		}

		chanTime := time.Now().Unix()
		_, err = insert.Exec(channel, nsfwInt, name, chanTime)
		insert.Close()
	}
	fmt.Println("Done adding to DB.")
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

// RemoveChannelFromDB Removes the channel from the database
func RemoveChannelFromDB(channel string) error {

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM queue WHERE channelID = ?", channel)

	return err
}

// UpdateChannelTime updates the time for the channel
func UpdateChannelTime(channel string) error {
	db, err := initDB()

	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	nowTime := time.Now().Unix()

	_, err = db.Exec("UPDATE channels SET time = ? WHERE channelID = ?", nowTime, channel)

	if err != nil {
		fmt.Println(err)
	}
	return err
}

// RemoveDormantChannels removes all channels that have not made a meme request in a month
func RemoveDormantChannels() error {

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	monthAgo := time.Now().Unix() - 2592000

	_, err = db.Exec("DELETE FROM channels WHERE time <= ?", monthAgo)

	return err

}

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

	defer db.Close()

	if err != nil {
		return QueueObj{}, err
	}

	row := db.QueryRow("SELECT time_interval, subreddits, nsfw, time FROM queue WHERE channelID = ?", channel)
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
	queue.Type = "media" // this needs to be implemented. Either make it automatic or make the user specify
	return queue, nil

}

// DeleteMemeQueue clears the meme queue for the specified channel
func DeleteMemeQueue(channel string) error {
	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM queue WHERE channelID = ?", channel)

	return err
}

// SetMemeQueue sets the meme queue for the channel
func SetMemeQueue(channel string, nsfw bool, interval string, subs string) error {
	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	var nsfwInt int

	if nsfw {
		nsfwInt = 1
	} else {
		nsfwInt = 0
	}

	output, err := db.Prepare("SELECT channelID from queue WHERE channelID = ?")
	defer output.Close()

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	err = output.QueryRow(channel).Scan(&channel)

	if err == nil {
		// update
		_, err = db.Exec("UPDATE queue SET nsfw = ?, time_interval = ?, subreddits = ? WHERE channelID = ?", nsfwInt, interval, subs, channel)

	} else if err == sql.ErrNoRows {
		// this isn't working.
		// Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'interval, subreddits) VALUES (?, ?, ?, ?)' at line 1
		insert, err := db.Prepare("INSERT INTO queue (channelID, nsfw, time_interval, subreddits) VALUES (?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, err = insert.Exec(channel, nsfwInt, interval, subs)
		insert.Close()
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// UpdateMemeQueueTime updates the time for a queue item.
func UpdateMemeQueueTime(channel string, setTime int64) error {

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE queue SET time = ? WHERE channelID = ?", setTime, channel)

	return err
}

// GetAllQueueChannels gets all of the queue channels
func GetAllQueueChannels() ([]string, error) {

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT channelID FROM queue")

	var channels []string
	var channel string

	for rows.Next() {
		err = rows.Scan(&channel)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}

	return channels, err
}

//RemoveChannelFromDBAllTables removes a channel from all tables in the database. For examples such as if someone deletes a guild or channel
func RemoveChannelFromDBAllTables(channel string) error {
	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM queue WHERE channelID = ?", channel)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	_, err = db.Exec("DELETE FROM channels WHERE channelID = ?", channel)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	_, err = db.Exec("DELETE FROM banned_subs WHERE channelID = ?", channel)

	if err == sql.ErrNoRows {
		return nil
	}
	return err
}
