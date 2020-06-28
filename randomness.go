package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"

	_ "github.com/go-mysql/errors"
	_ "github.com/go-sql-driver/mysql"
)

const min int = 0
const max int = 2
const secondMax int = 5

type randomHotDog struct {
	TypeArray      []string `json:"TypeArray"`
	CondimentArray []string `json:"CondimentArray"`
	CaloriesArray  []int    `json:"CaloriesArray"`
	NameArray      []string `json:"Name"`
}

type randomHamburger struct {
	TypeArray      []string `json:"TypeArray"`
	CondimentArray []string `json:"CondimentArray"`
	CaloriesArray  []int    `json:"CaloriesArray"`
	NameArray      []string `json:"Name"`
}

func randomRole() string {
	theRole := ""
	randomNum := rand.Intn(max-min) + min
	switch randomNum {
	case 0:
		//Role for User
		theRole = "user"
		break
	case 1:
		//Role for Admin
		theRole = "admin"
		break
	case 2:
		//Role for IT
		theRole = "IT"
		break
	default:
		fmt.Println("Error choosing IT Role.")
		theRole = "ERROR"
		break
	}
	return theRole
}

func randomID() int {

	//Make User and USERID
	goodNum := false
	theID := 0

	for goodNum == false {
		//Query the database for all IDS
		row, err := db.Query(`SELECT user_id FROM users;`)
		check(err)
		defer row.Close()
		//Build the random, unique integer to be assigned to this User
		goodNumFound := true //A second checker to break this loop
		randInt := 0         //The random integer added onto ID
		var databaseID int   //The ID returned from the database while searching
		randIntString := ""  //The integer built through a string...
		min, max := 0, 9     //The min and Max value for our randInt
		for i := 0; i < 8; i++ {
			randInt = rand.Intn(max-min) + min
			randIntString = randIntString + strconv.Itoa(randInt)
		}
		theID, err = strconv.Atoi(randIntString)
		if err != nil {
			fmt.Println(err)
		}
		//Check to see if the built number is taken.
		for row.Next() {
			err = row.Scan(&databaseID)
			check(err)
			if databaseID == theID {
				//Found the number, need to create another one!
				goodNumFound = false
				break
			} else {

			}
		}
		//Final check to see if we need to go through this loop again
		if goodNumFound == false {
			goodNum = false
		} else {
			goodNum = true
		}
	}

	return theID
}

func randomPassword(pWord string) string {
	bsString := []byte(pWord)                     //Encode Password
	encodedString := hex.EncodeToString(bsString) //Encode Password Pt2

	return encodedString
}

func giveRandomFood(userID int) {
	//Declare food
	var takenFoods []int
	hotDogArray := randomHotDog{
		TypeArray:      []string{"Alcaholic", "Ordinary", "Nickelback", "Big", "Space"},
		CondimentArray: []string{"Barcardi", "Plainess", "Awfulness", "Giga size", "Blackness"},
		CaloriesArray:  []int{800, 20, 300, 1400, 1},
		NameArray:      []string{"The Boozdog", "The PlainDog", "The Bad Music Dog", "The Bigdog", "The NeildeGrasse Dog"},
	}

	hamburgerArray := randomHamburger{
		TypeArray:      []string{"Alcaholic", "Ordinary", "Nickelback", "Big", "Space"},
		CondimentArray: []string{"Barcardi", "Plainess", "Awfulness", "Giga size", "Blackness"},
		CaloriesArray:  []int{800, 20, 300, 1400, 1},
		NameArray:      []string{"The Boozburger", "The Plainburger", "The Bad Music Burger", "The Big Burger", "The NeildeGrasse Burger"},
	}

	//Assign Hamburger Nums
	takenFoods = takenFoods[:0]
	for x := 0; x < 3; x++ {
		//Select random Hamburger
		goodFood := false //Determines if the food we've assembled is good.
		for goodFood == false {
			allGood := true //Determines if food is found in the 'takenFoods'
			randomNum := rand.Intn(secondMax-min) + min
			//See if food number is in takenFoods
			for j := 0; j < len(takenFoods); j++ {
				if takenFoods[j] == randomNum {
					allGood = false //food found, must start over
					break
				} else {

				}
			}
			if allGood == true {
				goodFood = true //Food not found, we can add it to takenFoods
				takenFoods = append(takenFoods, randomNum)
				fmt.Printf("DEBUG: Added %v to the takenFoods: %v\n", randomNum, takenFoods)
			} else {

			}
		}
	}
	//Give 3 Hamburger
	for z := 0; z < len(takenFoods); z++ {
		newHamburger := Hamburger{hamburgerArray.TypeArray[takenFoods[z]],
			hamburgerArray.CondimentArray[takenFoods[z]],
			hamburgerArray.CaloriesArray[takenFoods[z]],
			hamburgerArray.NameArray[takenFoods[z]],
			userID}
		//wg.Add(1) Need to add our wait groups for the program
		//go
		insertHamburger(newHamburger)
	}
	//Give 3 Hotdogs
	takenFoods = takenFoods[:0]
	for x := 0; x < 3; x++ {
		//Select random Hotdog
		goodFood := false //Determines if the food we've assembled is good.
		for goodFood == false {
			allGood := true //Determines if food is found in the 'takenFoods'
			randomNum := rand.Intn(secondMax-min) + min
			//See if food number is in takenFoods
			for j := 0; j < len(takenFoods); j++ {
				if takenFoods[j] == randomNum {
					allGood = false //food found, must start over
					break
				} else {

				}
			}
			if allGood == true {
				goodFood = true //Food not found, we can add it to takenFoods
				takenFoods = append(takenFoods, randomNum)
				fmt.Printf("DEBUG: Added %v to the takenFoods: %v\n", randomNum, takenFoods)
			} else {

			}
		}
	}
	//Give 3 Hotdogs
	for z := 0; z < len(takenFoods); z++ {
		newHotdog := Hotdog{hotDogArray.TypeArray[takenFoods[z]],
			hotDogArray.CondimentArray[takenFoods[z]],
			hotDogArray.CaloriesArray[takenFoods[z]],
			hotDogArray.NameArray[takenFoods[z]],
			userID}
		//wg.Add(1) Need to add our wait groups for the program
		//go
		insertHotDog(newHotdog)
	}
}
