package main

import (
	"database/sql"
	"testing"
	"time"
)

// Tests

func TestBenefitsSingle(t *testing.T) {
	testUserID := "696969696969696969"
	testGuild := "420420420420420420"

	// There will be two sub-tests
	// First one should return sql.ErrNoRows

	err := setBenefitServer(testUserID, 0, testGuild)

	if err != sql.ErrNoRows {
		t.Errorf("Error! Test should return sql.ErrNowRows: %v", err)
	}

	_, _, err = getBenefitServer(testUserID, testGuild)

	if err != sql.ErrNoRows {
		t.Errorf("Error! Test should return sql.ErrNowRows: %v", err)
	}

	// This test is for a normal user
	err = setBenefitServer(testUserID, 1, testGuild)

	if err != nil {
		t.Errorf("There was an error setting the benefits for the server: %v", err)
	}

	status, cooldown, err := getBenefitServer(testUserID, testGuild)

	if err != nil {
		t.Errorf("There was an error getting the benefits for the server: %v", err)
	}

	if status != 1 {
		t.Errorf("There was an error getting the status for the server, expected 1, got %d", status)
	}

	if cooldown < time.Now().Unix() {
		t.Errorf("There was an error getting the cooldown for the server, expected to be less than %d, but got %d", time.Now().Unix(), cooldown)
	}

	err = removeBenefitServer(testGuild)

	if err != nil {
		t.Errorf("There was an error removing the server from the benefits database: %v", err)
	}

}

func TestBenefitsMulti(t *testing.T) {
	var err error
	var testGuilds []string
	var returnGuilds []string
	var returnStatus uint8
	var found int

	const status = 2

	testUserID := "696969696969696969"

	for i := 0; i < 3; i++ {
		testGuilds = append(testGuilds, randString(18))
	}

	for _, testGuild := range testGuilds {
		err = setBenefitServer(testUserID, status, testGuild)
		if err != nil {
			t.Errorf("There was an error setting the benefits for the server: %v", err)
		}
	}

	returnStatus, returnGuilds, err = getAllBenefitsForUser(testUserID)

	if err != nil {
		t.Errorf("There was an error getting the benefits for the servers: %v", err)
	}

	if returnStatus != status {
		t.Errorf("There was an error getting the status for the server, expected %d, got %d", status, returnStatus)
	}

	found = 0
	for _, returnGuild := range returnGuilds {
		// this is because the database may not keep or return
		// the test IDs in the same order as they went in
		if stringInSlice(returnGuild, testGuilds) {
			found++
		}
	}

	if found != len(testGuilds) {
		t.Errorf("The returned guilds does not equal the original list of guilds. Got %v, expected %v", returnGuilds, testGuilds)
	}

	err = removeAllBenefitsForUser(testUserID)

	if err != nil {
		t.Errorf("There was an error removing the servers from the benefits database: %v", err)
	}

}
