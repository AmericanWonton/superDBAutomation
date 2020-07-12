package main

import (
	"fmt"
	"time"
)

const successMessage string = "Successful Insert"
const failureMessage string = "Unsuccessful Insert"

//POST hotdog, Mainpage
func insertHotDog(aHotdogs []Hotdog) {
	postedHotDogs := aHotdogs

	for x := 0; x < len(aHotdogs); x++ {
		stmt, err := db.Prepare("INSERT INTO hot_dogs(TYPE, CONDIMENT, CALORIES, NAME, USER_ID) VALUES(?,?,?,?,?)")
		defer stmt.Close()

		r, err := stmt.Exec(postedHotDogs[x].HotDogType, postedHotDogs[x].Condiment, postedHotDogs[x].Calories,
			postedHotDogs[x].Name, postedHotDogs[x].UserID)
		check(err)

		n, err := r.RowsAffected()
		check(err)

		fmt.Printf("DEBUG: %v rows effected.\n", n)
	}
}

//INSERT HOTDOG
func insertHamburgers(aBurgers []Hamburger) {
	postedHamburgers := aBurgers
	for x := 0; x < len(postedHamburgers); x++ {
		stmt, err := db.Prepare("INSERT INTO hamburgers(TYPE, CONDIMENT, CALORIES, NAME, USER_ID) VALUES(?,?,?,?,?)")
		defer stmt.Close()

		r, err := stmt.Exec(postedHamburgers[x].BurgerType, postedHamburgers[x].Condiment,
			postedHamburgers[x].Calories, postedHamburgers[x].Name, postedHamburgers[x].UserID)
		check(err)

		n, err := r.RowsAffected()
		check(err)
		fmt.Printf("DEBUG: %v rows effected.\n", n)
	}
}

//INSERT USER(s)
func insertUser(aUser User) {
	defer wg.Done() //For Wait Group
	//Marshal it into our type
	postedUser := aUser

	//Add User to the SQL Database
	stmt, err := db.Prepare("INSERT INTO users(USERNAME, PASSWORD, FIRSTNAME, LASTNAME, ROLE, USER_ID) VALUES(?,?,?,?,?,?)")
	defer stmt.Close()

	r, err := stmt.Exec(postedUser.UserName, postedUser.Password, postedUser.First,
		postedUser.Last, postedUser.Role, postedUser.UserID)
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Printf("Inserted User Record for SQL: %v\n", n)
	//Print log info
	insertionString := "Inserted User record for SQL: " + string(n)
	logWriter(insertionString)

	//Now we insert everything for Mongo
	//Add User to MongoDB
	fmt.Printf("DEBUG: Adding User to MongoDB\n")
	theTimeNow := time.Now()
	var insertionUser AUser = AUser{
		UserName:    postedUser.UserName,
		Password:    postedUser.Password,
		First:       postedUser.First,
		Last:        postedUser.Last,
		Role:        postedUser.Role,
		UserID:      randomIDCreation(),
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		Hotdogs:     MongoHotDogs{},
		Hamburgers:  MongoHamburgers{},
	}
	insertionUsers := TheUsers{
		Users: []AUser{insertionUser},
	}
	insertUsers(insertionUsers)
}
