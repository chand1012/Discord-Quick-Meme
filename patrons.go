package main

import "database/sql"

// gets the patron status, as in what level of patron they are.
// Currently, there is only levels 0, 1, 2. 0 means not a patron.
// level 1 is single server and level 2 is 3 server.
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

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return status, nil
}
