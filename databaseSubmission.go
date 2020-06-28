package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const successMessage string = "Successful Insert"
const failureMessage string = "Unsuccessful Insert"

//POST hotdog, Mainpage
func insertHotDog(aHotdog Hotdog) {
	postedHotDog := aHotdog

	//Protections for the hotdog name
	if strings.Compare(postedHotDog.HotDogType, "DEBUGTYPE") == 0 {
		postedHotDog.HotDogType = "NONE"
	}

	stmt, err := db.Prepare("INSERT INTO hot_dogs(TYPE, CONDIMENT, CALORIES, NAME, USER_ID) VALUES(?,?,?,?,?)")
	defer stmt.Close()

	r, err := stmt.Exec(postedHotDog.HotDogType, postedHotDog.Condiment, postedHotDog.Calories, postedHotDog.Name, postedHotDog.UserID)
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Printf("DEBUG: %v rows effected.\n", n)

	if err != nil {
		fmt.Println(failureMessage)
	} else {
		hDogMarshaled, err := json.Marshal(postedHotDog)
		if err != nil {
			fmt.Printf("Error with %v\n", hDogMarshaled)
		}
		hDogSuccessMSG := successMessage + string(hDogMarshaled)
		fmt.Printf(hDogSuccessMSG)
	}
	//wg.Done() //Done with this wait group
}

//INSERT HOTDOG
func insertHamburger(aBurger Hamburger) {
	postedHamburger := aBurger

	//Protections for the hamburger name
	if strings.Compare(postedHamburger.BurgerType, "DEBUGTYPE") == 0 {
		postedHamburger.BurgerType = "NONE"
	}

	fmt.Printf("DEBUG: HERE IS OUR postedHamburger: \n%v\n", postedHamburger)

	stmt, err := db.Prepare("INSERT INTO hamburgers(TYPE, CONDIMENT, CALORIES, NAME, USER_ID) VALUES(?,?,?,?,?)")
	defer stmt.Close()

	r, err := stmt.Exec(postedHamburger.BurgerType, postedHamburger.Condiment,
		postedHamburger.Calories, postedHamburger.Name, postedHamburger.UserID)
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Printf("DEBUG: %v rows effected.\n", n)

	if err != nil {
		fmt.Printf(failureMessage)
	} else {
		hamMarshaled, err := json.Marshal(postedHamburger)
		if err != nil {
			fmt.Printf("Error with %v\n", hamMarshaled)
		}
		hamSuccessMSG := successMessage + string(hamMarshaled)
		fmt.Println(hamSuccessMSG)
	}
	//wg.Done() //Done with wait group
}

//INSERT USER(s)
func insertUser(aUser User) {

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

	fmt.Printf("Inserted Record: %v\n", n)
}
