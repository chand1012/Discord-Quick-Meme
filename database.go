package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Channel struct {
	ChannelID string             `json:"channel"`
	Nsfw      bool               `json:"nsfw"`
	Name      string             `json:"name"`
	Time      int64              `json:"time"`
	ID        primitive.ObjectID `bson:"_id,omitempty"`
}

type BannedSub struct {
	ChannelID string             `json:"channelID"`
	Sub       string             `json:"subreddit"`
	ID        primitive.ObjectID `bson:"_id,omitempty"`
}

func ConnectMongo() (*mongo.Client, context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_CONNECT_STR")))

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, ctx

}

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

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	var channelObject Channel
	channelCache := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("channels")

	filter := bson.M{
		"channel": channel,
	}

	err := channelCache.FindOne(context.TODO(), filter, options.FindOne()).Decode(&channelObject)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}

	if err == nil { // if there is no error that means that the entry must be updated. Ideally this won't happen often.
		fmt.Println("Already in database, updating existing record...")
		channelObject.Time = time.Now().Unix()
		filter = bson.M{"_id": bson.M{"$eq": channelObject.ID}}
		update := bson.M{
			"$set": bson.M{
				"channel": channelObject.ChannelID,
				"nsfw":    channelObject.Nsfw,
				"time":    channelObject.Time,
				"name":    channelObject.Name,
			},
		}

		_, err = channelCache.UpdateOne(context.TODO(), filter, update)
	} else if err == mongo.ErrNoDocuments {
		fmt.Println("Adding new entry...")
		channelObject.ChannelID = channel
		channelObject.Name = name
		channelObject.Nsfw = nsfw
		channelObject.Time = time.Now().Unix()
		_, err = channelCache.InsertOne(context.TODO(), channelObject)
	}
	fmt.Println("Done adding to DB.")
	return err
}

// GetChannelFromDB gets the channel from the database
func GetChannelFromDB(channel string) (bool, string, error) {

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	var channelObject Channel
	channelCache := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("channels")

	filter := bson.M{
		"channel": channel,
	}

	err := channelCache.FindOne(context.TODO(), filter, options.FindOne()).Decode(&channelObject)

	if err != nil {
		return false, "", err
	}

	return channelObject.Nsfw, channelObject.Name, nil

}

// RemoveChannelFromDB Removes the channel from the database
func RemoveChannelFromDB(channel string) error {

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	channelCache := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("channels")

	filter := bson.M{
		"channel": channel,
	}

	_, err := channelCache.DeleteOne(context.TODO(), filter)

	return err
}

// UpdateChannelTime updates the time for the channel
func UpdateChannelTime(channel string) error {
	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	channelCache := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("channels")

	filter := bson.M{
		"channel": channel,
	}

	update := bson.M{
		"$set": bson.M{
			"time": time.Now().Unix(),
		},
	}

	_, err := channelCache.UpdateOne(context.TODO(), filter, update)

	return err

}

// RemoveDormantChannels removes all channels that have not made a meme request in a month
func RemoveDormantChannels() error {

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	channelCache := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("channels")

	filter := bson.M{
		"$lte": bson.M{
			"time": time.Now().Unix() - 2592000,
		},
	}

	_, err := channelCache.DeleteMany(context.TODO(), filter)

	return err
}

// SetBannedSubreddit adds a sub to the channel bans
func SetBannedSubreddit(channel string, subreddit string) error {

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	bannedSubs := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("banned_subs")

	var bannedSub BannedSub

	bannedSub.ChannelID = channel
	bannedSub.Sub = subreddit

	_, err := bannedSubs.InsertOne(context.TODO(), bannedSub)

	return err

}

// RemoveBannedSubreddit removes a sub from the bans list
func RemoveBannedSubreddit(channel string, subreddit string) error {

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	bannedSubs := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("banned_subs")

	filter := bson.M{
		"channel":   channel,
		"subreddit": subreddit,
	}

	_, err := bannedSubs.DeleteOne(context.TODO(), filter)

	return err
}

// GetAllBannedSubs gets all of the banned subreddits
func GetAllBannedSubs(channel string) ([]string, error) {
	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	bannedSubs := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("banned_subs")l, err
	
	filter := bson.M{
		"channel": channel,
	}

	cursor, err := bannedSubs.Find(dbContext, filter)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(dbContext)

	var subs []string

	for cursor.Next(dbContext) {
		var bannedSub BannedSub
		if err = cursor.Decode(&bannedSub); err != nil {
			return nil, err
		}
		subs = append(subs, bannedSub.Sub)
	}

	return subs, nil
	
}

// GetMemeQueue gets data for the meme queue for the specified channel
func GetMemeQueue(channel string) (QueueObj, error) {
	var queue QueueObj
	var subString string

	
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
		_, err = db.Exec("UPDATE queue SET nsfw = ?, timeInterval = ?, subreddits = ? WHERE channelID = ?", nsfwInt, interval, subs, channel)

	} else if err == sql.ErrNoRows {
		// this isn't working.
		// Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'interval, subreddits) VALUES (?, ?, ?, ?)' at line 1
		insert, err := db.Prepare("INSERT INTO queue (channelID, nsfw, timeInterval, subreddits) VALUES (?, ?, ?, ?)")
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
