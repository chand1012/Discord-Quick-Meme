package main

import (
	"database/sql"
	"fmt"
	"time"
)

func getPatronStatus(userID string) (uint8, error) {
	db, err := initDB()

	defer db.Close()

	if err != nil {
		return 0, err
	}

	output, err := db.Prepare("SELECT status FROM patron WHERE userID = ?")

	defer output.Close()

	var status uint8

	err = output.QueryRow(userID).Scan(&status)

	if err != nil {
		return 0, err
	}

	return status, nil
}

func setBenefitServer(userID string, guild string) error {
	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	insert, err := db.Prepare("INSERT INTO boosted (userID, status, guilds, cooldown) VALUES ?, ?, ?, ?")

	if err != nil {
		fmt.Println(err)
		return err
	}

	// this can also be used to check whether
	// a user is a patron or not
	status, err := getPatronStatus(userID)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if status == 0 {
		return sql.ErrNoRows
	}

	cooldown := time.Now().Unix() + 2700000

	_, err = insert.Exec(userID, status, guild, cooldown)
	insert.Close()

	return err
}

func removeBenefitServer(userID string, guild string) error {
	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	status, err := getPatronStatus(userID)

	if err != nil {
		return err
	}

	if status == 0 {
		return sql.ErrNoRows
	}

	_, err = db.Exec("DELETE FROM boosted WHERE guildID = ?", guild)

	return err

}
