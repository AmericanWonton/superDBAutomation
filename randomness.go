package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-mysql/errors"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2/bson"
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
	fmt.Printf("DEBUG: Creating Random ID for User/Food\n")
	finalID := 0        //The final, unique ID to return to the food/user
	randInt := 0        //The random integer added onto ID
	randIntString := "" //The integer built through a string...
	min, max := 0, 9    //The min and Max value for our randInt
	foundID := false
	for foundID == false {
		randInt = 0
		randIntString = ""
		//Create the random number, convert it to string

		for i := 0; i < 8; i++ {
			randInt = rand.Intn(max-min) + min
			randIntString = randIntString + strconv.Itoa(randInt)
		}
		//Once we have a string of numbers, we can convert it back to an integer
		theID, err := strconv.Atoi(randIntString)
		if err != nil {
			fmt.Printf("We got an error converting a string back to a number, %v\n", err)
			fmt.Printf("Here is randInt: %v\n and randIntString: %v\n", randInt, randIntString)
			fmt.Println(err)
			log.Fatal(err)
		}
		//Search all our collections to see if this UserID is unique
		canExit := []bool{true, true, true}
		fmt.Printf("DEBUG: We are going to see if this ID is in our food or User DBs: %v\n", theID)
		//User collection
		userCollection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
		var testAUser AUser
		theErr := userCollection.FindOne(theContext, bson.M{"userid": theID}).Decode(&testAUser)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				fmt.Printf("It's all good, this document wasn't found for User and our ID is clean.\n")
				canExit[0] = true
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: \n%v\n", theErr)
				canExit[0] = false
				log.Fatal(theErr)
			}
		}
		//Check hotdog collection
		hotdogCollection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
		var testHotdog MongoHotDog
		//Give 0 values to determine if these IDs are found
		theFilter := bson.M{
			"$or": []interface{}{
				bson.M{"userid": theID},
				bson.M{"foodid": theID},
			},
		}
		theErr = hotdogCollection.FindOne(theContext, theFilter).Decode(&testHotdog)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				fmt.Printf("It's all good, this document wasn't found for User/Hotdog and our ID is clean.\n")
				canExit[1] = true
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: \n%v\n", theErr)
				canExit[1] = false
			}
		}
		//Check hamburger collection
		hamburgerCollection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
		var testBurger MongoHamburger
		//Give 0 values to determine if these IDs are found
		theFilter2 := bson.M{
			"$or": []interface{}{
				bson.M{"userid": theID},
				bson.M{"foodid": theID},
			},
		}
		theErr = hamburgerCollection.FindOne(theContext, theFilter2).Decode(&testBurger)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				canExit[2] = true
				fmt.Printf("It's all good, this document wasn't found for User/hamburger and our ID is clean.\n")
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: \n%v\n", theErr)
				canExit[2] = false
			}
		}
		//Final check to see if we can exit this loop
		if canExit[0] == true && canExit[1] == true && canExit[2] == true {
			finalID = theID
			foundID = true
		} else {
			foundID = false
		}
	}

	return finalID
}

func randomPassword(pWord string) string {
	bsString := []byte(pWord)                     //Encode Password
	encodedString := hex.EncodeToString(bsString) //Encode Password Pt2

	return encodedString
}

func giveRandomFood(userID int, newUser AUser) {
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
			} else {

			}
		}
	}
	var theHamburgers []Hamburger
	//Give 3 Hamburger
	for z := 0; z < len(takenFoods); z++ {
		theTimeNow := time.Now()
		newHamburger := Hamburger{hamburgerArray.TypeArray[takenFoods[z]],
			hamburgerArray.CondimentArray[takenFoods[z]],
			hamburgerArray.CaloriesArray[takenFoods[z]],
			hamburgerArray.NameArray[takenFoods[z]],
			userID, randomIDCreation(), theTimeNow.Format("2006-01-02 15:04:05"), theTimeNow.Format("2006-01-02 15:04:05")}
		theHamburgers = append(theHamburgers, newHamburger)
	}

	wg.Add(1)
	go insertHamburgers(theHamburgers)
	//Assign Hotdog Nums
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
			} else {

			}
		}
	}
	var theHotdogs []Hotdog
	//Give 3 Hotdogs
	for z := 0; z < len(takenFoods); z++ {
		theTimeNow := time.Now()
		newHotdog := Hotdog{hotDogArray.TypeArray[takenFoods[z]],
			hotDogArray.CondimentArray[takenFoods[z]],
			hotDogArray.CaloriesArray[takenFoods[z]],
			hotDogArray.NameArray[takenFoods[z]],
			userID, randomIDCreation(), theTimeNow.Format("2006-01-02 15:04:05"), theTimeNow.Format("2006-01-02 15:04:05")}
		theHotdogs = append(theHotdogs, newHotdog)
	}

	wg.Add(1)
	go insertHotDogs(theHotdogs)
	//Print log info
	logWriter("Finished giving random food for SQL.")

	//Give Food for Mongo DB as well
	var insertHotDogs MongoHotDogs
	var insertHamburgers MongoHamburgers
	//Put previous Hotdog/Hamburgers into "MongoHotdog/Hamburger"
	for i := 0; i < len(theHotdogs); i++ {
		theTimeNow := time.Now()
		newMongoDog := MongoHotDog{
			HotDogType:  theHotdogs[i].HotDogType,
			Condiments:  []string{theHotdogs[i].Condiment},
			Name:        theHotdogs[i].Name,
			FoodID:      randomIDCreation(),
			UserID:      theHotdogs[i].UserID,
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		insertHotDogs.Hotdogs = append(insertHotDogs.Hotdogs, newMongoDog)
	} //For Hotdogs
	for p := 0; p < len(theHamburgers); p++ {
		theTimeNow := time.Now()
		newMongoHamb := MongoHamburger{
			BurgerType:  theHotdogs[p].HotDogType,
			Condiments:  []string{theHotdogs[p].Condiment},
			Name:        theHotdogs[p].Name,
			FoodID:      randomIDCreation(),
			UserID:      theHotdogs[p].UserID,
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		insertHamburgers.Hamburgers = append(insertHamburgers.Hamburgers, newMongoHamb)
	}
	//give newUser the food to update in Mongo
	newUser.Hotdogs = insertHotDogs
	newUser.Hamburgers = insertHamburgers
	//InsertHotdog/Hamburgers for Mongo
	wg.Add(1)
	go updateUserMongo(newUser)
	wg.Add(1)
	go insertHotDogsMongo(insertHotDogs)
	wg.Add(1)
	go insertHamburgersMongo(insertHamburgers)
	logWriter("Finished giving random food for Mongo.")
	wg.Done() //For GoRoutine
}
