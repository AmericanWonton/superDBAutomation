package main

import (
	"fmt"
	"time"
)

const successMessage string = "Successful Insert"
const failureMessage string = "Unsuccessful Insert"

var dbConnectString string

//POST hotdog, Mainpage
func insertHotDogs(aHotdogs []Hotdog) {
	postedHotDogs := aHotdogs

	for x := 0; x < len(aHotdogs); x++ {
		stmt, err := db.Prepare("INSERT INTO hot_dogs(TYPE, CONDIMENT, CALORIES, NAME, USER_ID, FOOD_ID, DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?)")

		r, err := stmt.Exec(postedHotDogs[x].HotDogType, postedHotDogs[x].Condiment, postedHotDogs[x].Calories,
			postedHotDogs[x].Name, postedHotDogs[x].UserID, postedHotDogs[x].FoodID, postedHotDogs[x].DateCreated, postedHotDogs[x].DateUpdated)
		check(err)

		n, err := r.RowsAffected()
		check(err)

		stmt.Close()

		fmt.Printf("DEBUG: %v rows effected.\n", n)
	}
	wg.Done() //For GoRoutines
}

//INSERT Hamburgers
func insertHamburgers(aBurgers []Hamburger) {
	postedHamburgers := aBurgers
	for x := 0; x < len(postedHamburgers); x++ {
		stmt, err := db.Prepare("INSERT INTO hamburgers(TYPE, CONDIMENT, CALORIES, NAME, USER_ID, FOOD_ID, DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?)")

		r, err := stmt.Exec(postedHamburgers[x].BurgerType, postedHamburgers[x].Condiment,
			postedHamburgers[x].Calories, postedHamburgers[x].Name, postedHamburgers[x].UserID,
			postedHamburgers[x].FoodID, postedHamburgers[x].DateCreated, postedHamburgers[x].DateUpdated)
		check(err)

		n, err := r.RowsAffected()
		check(err)
		fmt.Printf("DEBUG: %v rows effected.\n", n)
		stmt.Close()
	}
	wg.Done()
}

//INSERT USER(s)
func insertUsers(theUsers []User) {
	//Marshal it into our type
	postedUsers := theUsers

	//Add User to the SQL Database
	for x := 0; x < len(theUsers); x++ {
		stmt, err := db.Prepare("INSERT INTO users(USERNAME, PASSWORD, FIRSTNAME, LASTNAME, ROLE, USER_ID, DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?)")

		r, err := stmt.Exec(postedUsers[x].UserName, postedUsers[x].Password, postedUsers[x].First,
			postedUsers[x].Last, postedUsers[x].Role, postedUsers[x].UserID,
			postedUsers[x].DateCreated, postedUsers[x].DateUpdated)
		check(err)

		n, err := r.RowsAffected()
		check(err)

		stmt.Close()

		fmt.Printf("Inserted User Record for SQL: %v\n", n)
		//Print log info
		insertionString := "Inserted User record for SQL: " + string(n)
		logWriter(insertionString)
	}

	//Now we insert everything for Mongo
	//Add Users to MongoDB
	insertionUsers := TheUsers{
		Users: []AUser{},
	}
	for j := 0; j < len(theUsers); j++ {
		theTimeNow := time.Now()
		insertionUser := AUser{
			UserName:    postedUsers[j].UserName,
			Password:    postedUsers[j].Password,
			First:       postedUsers[j].First,
			Last:        postedUsers[j].Last,
			Role:        postedUsers[j].Role,
			UserID:      randomIDCreation(),
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
			Hotdogs:     MongoHotDogs{},
			Hamburgers:  MongoHamburgers{},
		}
		insertionUsers.Users = append(insertionUsers.Users, insertionUser)
	}

	insertUsersMongo(insertionUsers)

	//Give Users random food
	for q := 0; q < len(insertionUsers.Users); q++ {
		wg.Add(1)
		go giveRandomFood(insertionUsers.Users[q].UserID, insertionUsers.Users[q])
	}
	wg.Done() //For GoRoutines
}
