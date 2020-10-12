package main

import (
	"database/sql"
	"fmt"
)

func resetPatronCache() {
	fmt.Println("Resetting Patron cache...")
	PatronCache = make(map[string]uint8)
}

// gets the patron status, as in what level of patron they are.
// Currently, there is only levels 0, 1, 2. 0 means not a patron.
// level 1 is single server and level 2 is three server.
func getPatronStatus(userID string, cache bool) (uint8, error) {
	var status uint8

	if cache {
		status, ok := PatronCache[userID]

		if ok {
			return status, nil
		}
	}

	db, err := initDB()

	defer db.Close()

	if err != nil {
		return 0, err
	}

	output, err := db.Prepare("SELECT status FROM patrons WHERE userID = ?")

	defer output.Close()

	err = output.QueryRow(userID).Scan(&status)

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if cache {
		PatronCache[userID] = status
	}

	return status, nil
}
