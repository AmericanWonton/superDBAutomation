package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	_ "github.com/go-mysql/errors"
	_ "github.com/go-sql-driver/mysql"
)

var logFile *os.File //used for logging

func logWriter(logMessage string) {
	//Logging info

	wd, _ := os.Getwd()
	logDir := filepath.Join(wd, "logging", "autogenlog.txt")
	logFile, err := os.OpenFile(logDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	defer logFile.Close()

	if err != nil {
		fmt.Println("Failed opening log file")
	}

	log.SetOutput(logFile)

	log.Println(logMessage)
}

//Here is our waitgroup
var wg sync.WaitGroup

//Here's our User struct
type User struct {
	UserName    string `json:"UserName"`
	Password    string `json:"Password"` //This was formally a []byte but we are changing our code to fit the database better
	First       string `json:"First"`
	Last        string `json:"Last"`
	Role        string `json:"Role"`
	UserID      int    `json:"UserID"`
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

type SpecialUser struct {
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MaidenName  string  `json:"maiden_name"`
	BirthData   string  `json:"birth_data"`
	PhoneH      string  `json:"phone_h"`
	PhoneW      string  `json:"phone_w"`
	EmailU      string  `json:"email_u"`
	EmailD      string  `json:"email_d"`
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	Domain      string  `json:"domain"`
	Useragent   string  `json:"useragent"`
	Ipv4        string  `json:"ipv4"`
	Macaddress  string  `json:"macaddress"`
	Plasticcard string  `json:"plasticcard"`
	Cardexpir   string  `json:"cardexpir"`
	Bonus       int     `json:"bonus"`
	Company     string  `json:"company"`
	Color       string  `json:"color"`
	UUID        string  `json:"uuid"`
	Height      int     `json:"height"`
	Weight      int     `json:"weight"`
	Blood       string  `json:"blood"`
	Eye         string  `json:"eye"`
	Hair        string  `json:"hair"`
	Pict        string  `json:"pict"`
	URL         string  `json:"url"`
	Sport       string  `json:"sport"`
	Ipv4URL     string  `json:"ipv4_url"`
	EmailURL    string  `json:"email_url"`
	DomainURL   string  `json:"domain_url"`
}

type UserCollection struct {
	TheUsers User `json:"TheUsers"`
}

//Below is our struct for Hotdogs/Hamburgers
type Hotdog struct {
	HotDogType  string `json:"HotDogType"`
	Condiment   string `json:"Condiment"`
	Calories    int    `json:"Calories"`
	Name        string `json:"Name"`
	UserID      int    `json:"UserID"` //User WHOMST this hotDog belongs to
	FoodID      int    `json:"FoodID"`
	PhotoID     int    `json:"PhotoID"`
	PhotoSrc    string `json:"PhotoSrc"`
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

type Hamburger struct {
	BurgerType  string `json:"BurgerType"`
	Condiment   string `json:"Condiment"`
	Calories    int    `json:"Calories"`
	Name        string `json:"Name"`
	UserID      int    `json:"UserID"` //User WHOMST this hotDog belongs to
	FoodID      int    `json:"FoodID"`
	PhotoID     int    `json:"PhotoID"`
	PhotoSrc    string `json:"PhotoSrc"`
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

/* Mongo No-SQL Variable Declarations */
type AUser struct { //Using this for Mongo
	UserName    string          `json:"UserName"`
	Password    string          `json:"Password"` //This was formally a []byte but we are changing our code to fit the database better
	First       string          `json:"First"`
	Last        string          `json:"Last"`
	Role        string          `json:"Role"`
	UserID      int             `json:"UserID"`
	DateCreated string          `json:"DateCreated"`
	DateUpdated string          `json:"DateUpdated"`
	Hotdogs     MongoHotDogs    `json:"Hotdogs"`
	Hamburgers  MongoHamburgers `json:"Hamburgers"`
}

type TheUsers struct { //Using this for Mongo
	Users []AUser `json:"Users"`
}

type MongoHotDog struct {
	HotDogType  string   `json:"HotDogType"`
	Condiments  []string `json:"Condiments"`
	Calories    int      `json:"Calories"`
	Name        string   `json:"Name"`
	FoodID      int      `json:"FoodID"`
	UserID      int      `json:"UserID"` //User WHOMST this hotDog belongs to
	PhotoID     int      `json:"PhotoID"`
	PhotoSrc    string   `json:"PhotoSrc"`
	DateCreated string   `json:"DateCreated"`
	DateUpdated string   `json:"DateUpdated"`
}

type MongoHotDogs struct {
	Hotdogs []MongoHotDog `json:"Hotdogs"`
}

type MongoHamburger struct {
	BurgerType  string   `json:"BurgerType"`
	Condiments  []string `json:"Condiments"`
	Calories    int      `json:"Calories"`
	Name        string   `json:"Name"`
	FoodID      int      `json:"FoodID"`
	UserID      int      `json:"UserID"` //User WHOMST this hotDog belongs to
	PhotoID     int      `json:"PhotoID"`
	PhotoSrc    string   `json:"PhotoSrc"`
	DateCreated string   `json:"DateCreated"`
	DateUpdated string   `json:"DateUpdated"`
}

type MongoHamburgers struct {
	Hamburgers []MongoHamburger `json:"Hamburgers"`
}

//Mongo DB Declarations
var mongoClient *mongo.Client

//mySQL database declarations
var db *sql.DB
var err error

// Handle Errors
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

func main() {
	//open SQL connection
	db, err = sql.Open("mysql",
		"joek1:fartghookthestrong69@tcp(food-database.cd8ujtto1hfj.us-east-2.rds.amazonaws.com)/food-database-schema?charset=utf8")
	check(err)
	//db.SetMaxOpenConns(1) //Needed for DB
	defer db.Close()

	err = db.Ping()
	check(err)
	//Print to logs
	logWriter("Connected to SQL DB starting process")

	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed

	//Connect to MongoDB
	mongoClient = connectDB()

	//Launch our creations into their own goroutine
	//https://www.udemy.com/course/learn-how-to-code/learn/lecture/11922316#overview
	fmt.Printf("Launching go routines...\n")
	fmt.Printf("OS: %v\n", runtime.GOOS)
	fmt.Printf("ARCH: %v\n", runtime.GOARCH)
	fmt.Printf("CPUs: %v\n", runtime.NumCPU())

	fmt.Printf("Number of goRoutines: %v\n", runtime.NumGoroutine())
	//Here's our GoRoutines
	/*** SWEAR WORD REMOVERS ***/
	wg.Add(1)
	go swearUserRemoverHDog()
	wg.Add(1)
	go swearUserRemoverHam()
	/*** CREATE AND ADD USERS ***/
	wg.Add(1)
	go userCreator()
	/*** REMOVE UNUSED FOOD ***/
	wg.Add(1)
	go discardFood()
	//Need to tell our main program to wait for goroutines
	wg.Wait()
}

//Check errors in our mySQL errors
func check(err error) {
	if err != nil {
		fmt.Printf("Error in SQLDB: \n%v\n", err.Error())
		log.SetOutput(logFile)
		failureString := "Error with SQL: " + err.Error()
		logWriter(failureString)
	}
}

func userCreator() {
	fmt.Println("Making random Users.")
	var theUsers []User

	url := "https://api.namefake.com"
	method := "GET"
	//Make 5 Users
	for v := 0; v < 5; v++ {
		fmt.Printf("DEBUG: Starting to append a user.\n")
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
		}

		res, err := client.Do(req)
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		var theNewUser SpecialUser
		json.Unmarshal(body, &theNewUser)

		nameArray := strings.Fields(theNewUser.Name)
		fName := ""
		lName := ""
		for v := 0; v < len(nameArray); v++ {
			if v == 0 {
				fName = nameArray[v]
			} else if v >= 2 {
				break
			} else {
				lName = nameArray[v]
			}
		}

		//Create a User
		theTimeNow := time.Now()
		newUser := User{
			UserName:    theNewUser.Username,
			Password:    randomPassword(theNewUser.Password),
			First:       fName,
			Last:        lName,
			Role:        randomRole(),
			UserID:      randomID(),
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		theUsers = append(theUsers, newUser)
		fmt.Printf("DEBUG: Here is our newUser: %v\n", newUser)
	}
	//Go give Users food
	fmt.Printf("DEBUG: Users made, inserting Users.\n")
	wg.Add(1)
	go insertUsers(theUsers) //User inserted
	//Print to logs
	logWriter("Done creating Users.")
	wg.Done()
}
