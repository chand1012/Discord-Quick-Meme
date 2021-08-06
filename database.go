package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
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
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_CONNECT_STR")))

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

// GetAllChannelNames gets all channel names
func GetAllChannelNames() {
	fmt.Println("Getting all current channel names and NSFW statuses...")
	starttime := GetMillis()

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	var channelObjects []Channel
	channelCache := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("channels")

	cursor, err := channelCache.Find(dbContext, bson.M{}, options.Find())

	if err != nil {
		fmt.Println(err)
		return
	}

	err = cursor.All(dbContext, &channelObjects)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, channelObject := range channelObjects {
		ServerMap[channelObject.ChannelID] = channelObject.Name
		NSFWMap[channelObject.ChannelID] = channelObject.Nsfw
	}
	endtime := GetMillis()
	t := endtime - starttime
	fmt.Println("Time to get all current channel names and NSFW status: " + strconv.FormatInt(t, 10) + "ms")
}

// AddChannelToDB adds a channel to the database
func AddChannelToDB(channel string, nsfw bool, name string) error {

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	var channelObject Channel
	channelCache := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("channels")

	filter := bson.M{
		"channelid": channel,
	}

	err := channelCache.FindOne(dbContext, filter, options.FindOne()).Decode(&channelObject)

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

		_, err = channelCache.UpdateOne(dbContext, filter, update)
	} else if err == mongo.ErrNoDocuments {
		fmt.Println("Adding new entry...")
		channelObject.ChannelID = channel
		channelObject.Name = name
		channelObject.Nsfw = nsfw
		channelObject.Time = time.Now().Unix()
		_, err = channelCache.InsertOne(dbContext, channelObject)
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
		"channelid": channel,
	}

	err := channelCache.FindOne(dbContext, filter, options.FindOne()).Decode(&channelObject)

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
		"channelid": channel,
	}

	_, err := channelCache.DeleteOne(dbContext, filter)

	return err
}

// UpdateChannelTime updates the time for the channel
func UpdateChannelTime(channel string) error {
	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	channelCache := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("channels")

	filter := bson.M{
		"channelid": channel,
	}

	update := bson.M{
		"$set": bson.M{
			"time": time.Now().Unix(),
		},
	}

	_, err := channelCache.UpdateOne(dbContext, filter, update)

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

	_, err := channelCache.DeleteMany(dbContext, filter)

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

	_, err := bannedSubs.InsertOne(dbContext, bannedSub)

	return err

}

// RemoveBannedSubreddit removes a sub from the bans list
func RemoveBannedSubreddit(channel string, subreddit string) error {

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	bannedSubs := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("banned_subs")

	filter := bson.M{
		"channelid": channel,
		"subreddit": subreddit,
	}

	_, err := bannedSubs.DeleteOne(dbContext, filter)

	return err
}

// GetAllBannedSubs gets all of the banned subreddits
func GetAllBannedSubs(channel string) ([]string, error) {
	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	bannedSubs := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("banned_subs")

	filter := bson.M{
		"channelid": channel,
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

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	memeQueue := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("queues")

	filter := bson.M{
		"channelid": channel,
	}

	err := memeQueue.FindOne(dbContext, filter).Decode(&queue)

	if err != nil {
		return QueueObj{}, err
	}

	return queue, nil
}

// DeleteMemeQueue clears the meme queue for the specified channel
func DeleteMemeQueue(channel string) error {
	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	memeQueue := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("queues")

	filter := bson.M{
		"channelid": channel,
	}

	_, err := memeQueue.DeleteOne(dbContext, filter)

	return err
}

// SetMemeQueue sets the meme queue for the channel
func SetMemeQueue(channel string, nsfw bool, interval string, subs string) error {
	var queue QueueObj
	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	memeQueue := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("queues")

	filter := bson.M{
		"channelid": channel,
		"nsfw":      nsfw,
		"interval":  interval,
		"subs":      subs,
	}

	err := memeQueue.FindOne(dbContext, filter).Decode(&queue)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}

	if err == mongo.ErrNoDocuments {
		queue.Channel = channel
		queue.NSFW = nsfw
		queue.Interval = interval
		queue.SubReddits = strings.Split(subs, ",")

		_, err = memeQueue.InsertOne(dbContext, queue)
	} else {
		update := bson.M{
			"$set": bson.M{
				"nsfw":     nsfw,
				"interval": interval,
				"subs":     subs,
			},
		}

		_, err = memeQueue.UpdateOne(dbContext, filter, update)
	}

	return err

}

// UpdateMemeQueueTime updates the time for a queue item.
func UpdateMemeQueueTime(channel string, setTime int64) error {

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	memeQueue := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("queues")

	filter := bson.M{
		"channelid": channel,
	}

	update := bson.M{
		"$set": bson.M{
			"time": setTime,
		},
	}

	_, err := memeQueue.UpdateOne(dbContext, filter, update)

	return err
}

// GetAllQueueChannels gets all of the queue channels
func GetAllQueueChannels() ([]string, error) {

	var channels []string

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	memeQueue := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("queues")

	cursor, err := memeQueue.Find(dbContext, bson.M{})

	if err != nil {
		return nil, err
	}

	for cursor.Next(dbContext) {
		var queue QueueObj
		if err = cursor.Decode(&queue); err != nil {
			return nil, err
		}
		channels = append(channels, queue.Channel)
	}

	return channels, nil
}

//RemoveChannelFromDBAllTables removes a channel from all tables in the database. For examples such as if someone deletes a guild or channel
func RemoveChannelFromDBAllTables(channel string) error {

	dbClient, dbContext := ConnectMongo()
	defer dbClient.Disconnect(dbContext)
	defer dbContext.Done()

	memeQueue := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("queues")
	bannedSubs := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("banned_subs")
	channelCache := dbClient.Database(os.Getenv("MONGO_DATABASE")).Collection("channels")

	filter := bson.M{
		"channelid": channel,
	}

	_, err := memeQueue.DeleteMany(dbContext, filter)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}

	_, err = bannedSubs.DeleteMany(dbContext, filter)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}

	_, err = channelCache.DeleteMany(dbContext, filter)

	if err == mongo.ErrNoDocuments {
		return nil
	}

	return err

}
