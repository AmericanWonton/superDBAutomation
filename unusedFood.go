package main

import (
	"fmt"
	"strconv"

	_ "github.com/go-mysql/errors"
	_ "github.com/go-sql-driver/mysql"
)

func discardFood() {
	//Here is all of our User IDs
	var userIDS []int
	stmt := "SELECT USER_ID from users"
	row, err := db.Query(stmt)
	check(err)

	//Put query results into userIDS
	var anID int
	for row.Next() {
		err = row.Scan(&anID)
		check(err)
		userIDS = append(userIDS, anID)
	}
	row.Close()
	//Collect all food IDs
	var hDogIDS []int
	var hamIDS []int
	var aFoodID int
	hDogStmt := "SELECT USER_ID FROM hot_dogs"
	hrow, err := db.Query(hDogStmt)
	check(err)
	for hrow.Next() {
		err = hrow.Scan(&aFoodID)
		check(err)
		hDogIDS = append(hDogIDS, aFoodID)
	}
	hrow.Close()
	hamStmt := "SELECT USER_ID FROM hamburgers"
	hamrow, err := db.Query(hamStmt)
	check(err)
	for hamrow.Next() {
		err = hamrow.Scan(&aFoodID)
		check(err)
		hamIDS = append(hamIDS, aFoodID)
	}
	hamrow.Close()
	//Go eliminate colleceted Hotdogs,(if they have values that need to be collected)
	eliminateHDogs(userIDS, hDogIDS)
	eliminateHam(userIDS, hamIDS)
	//Print log information
	fmt.Println("DEBUG: Finished discarding food in SQL")
	logWriter("Finished discarding food.")
}

func eliminateHDogs(theUserIDS []int, theHDogs []int) {
	theQuery := ""      //A query to be built for eliminating food with no user.
	foundValue := false //If true, we can run a query
	var allIDS []int
	for z := 0; z < len(theHDogs); z++ {
		foundUser := findID(theUserIDS, theHDogs[z])
		if foundUser == true {
			//Do nothing, it's got a user
		} else {
			allIDS = append(allIDS, theHDogs[z])
			foundValue = true //Needed to have this query run with values inside it
			//For first pass to build query
			if z == 0 {
				stringID := strconv.Itoa(theHDogs[z])
				theQuery = theQuery + "DELETE FROM hamburgers WHERE 1=1\nAND " +
					"USER_ID = " + stringID + "\n"
			} else {
				stringID := strconv.Itoa(theHDogs[z])
				theQuery = theQuery + "OR " + "USER_ID = " + stringID + "\n"
			}
		}
	}
	if foundValue == true {
		//Run the query to remove those values
		delH, err := db.Prepare(theQuery)
		check(err)

		r, err := delH.Exec()
		check(err)

		n, err := r.RowsAffected()
		check(err)

		fmt.Printf("Removed this many rows with no userID found for hotdogs: %v\n", n)
	} else {
		fmt.Println("No hotdogs with missing UserIDS to remove.")
	}
	//Print log information for SQL
	logWriter("Finished removing Hotdogs in SQL.")
	//Remove values for Mongo
	foodDeleteUnusedMongo(1, allIDS) //Delete HotDog Collection IDS
	fmt.Printf("Finished removing Unused Hotdogs in Mongo\n")
	logWriter("Finished removing Hotdogs in Mongo.")
}

func eliminateHam(theUserIDS []int, theHams []int) {
	var allIDS []int
	theQuery := ""      //A query to be built for eliminating food with no user.
	foundValue := false //If true, we can run a query
	for z := 0; z < len(theHams); z++ {
		foundUser := findID(theUserIDS, theHams[z])
		if foundUser == true {
			//Do nothing, it's got a user
		} else {
			allIDS = append(allIDS, theHams[z])
			foundValue = true //Needed to have this query run with values inside it
			//For first pass to build query
			if z == 0 {
				stringID := strconv.Itoa(theHams[z])
				theQuery = theQuery + "DELETE FROM hamburgers WHERE 1=1\nAND " +
					"USER_ID = " + stringID + "\n"
			} else {
				stringID := strconv.Itoa(theHams[z])
				theQuery = theQuery + "OR " + "USER_ID = " + stringID + "\n"
			}
		}
	}

	if foundValue == true {
		fmt.Printf("Here is the query we shall run to get rid of Hamburgers: \n\n%v\n\n", theQuery)
		//Run the query to remove those values
		delHams, err := db.Prepare(theQuery)
		check(err)

		r, err := delHams.Exec()
		check(err)

		n, err := r.RowsAffected()
		check(err)

		fmt.Printf("Removed this many rows with no userID found for hamburgers: %v\n", n)
	} else {
		fmt.Println("No hamburgers with missing UserIDs to remove.")
	}
	//Print log information
	logWriter("Finished removing hamburgers for SQL.")
	//Remove Hamburgers for Mongo
	foodDeleteUnusedMongo(2, allIDS) //Delete HotDog Collection IDS
	logWriter("Finished removing Hamburgers in Mongo.")
}

//If ID is not found in user table, remove from the food table
func findID(theUserIDS []int, possibleInt int) bool {
	theReturn := false
	for i := 0; i < len(theUserIDS); i++ {
		if theUserIDS[i] == possibleInt {
			theReturn = true
		} else {
			theReturn = false
		}
	}
	return theReturn
}
