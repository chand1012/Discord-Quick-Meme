package main

import (
	"database/sql"
	"fmt"
	"time"
)

func setBenefitServer(userID string, status uint8, guild string) error {

	if status == 0 { // in theory, this should never happen
		return sql.ErrNoRows
	}

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return err
	}

	insert, err := db.Prepare("INSERT INTO boosted (userID, status, guildID, cooldown) VALUES ?, ?, ?, ?")

	if err != nil {
		fmt.Println(err)
		return err
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

	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM boosted WHERE guildID = ?", guild)

	return err

}

func getBenefitServer(userID string, guildID string) (uint8, int64, error) {
	var status uint8
	var cooldown int64

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return 0, 0, err
	}

	output, err := db.Prepare("SELECT status, cooldown FROM boosted WHERE guildID = ? AND userID = ?")

	defer output.Close()

	if err != nil {
		return 0, 0, err
	}

	err = output.QueryRow(guildID, userID).Scan(&status, &cooldown)

	return status, cooldown, err
}

func getAllBenefitsForUser(userID string) (uint8, []string, error) {
	var status uint8
	var guildID string
	var guildIDs []string

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return 0, nil, err
	}

	rows, err := db.Query("SELECT status, guildID FROM boosted WHERE userID = ?", userID)

	if err != nil {
		return 0, nil, err
	}

	for rows.Next() {
		err = rows.Scan(&status, &guildID)
		if err != nil {
			return 0, nil, err
		}
		guildIDs = append(guildIDs, guildID)
	}

	err = rows.Err()

	return status, guildIDs, err

}
