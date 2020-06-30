package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"

	"github.com/gobuffalo/packr/v2"
	"github.com/sirupsen/logrus"

	_ "github.com/go-mysql/errors"
	_ "github.com/go-sql-driver/mysql"
)

//Here is our waitgroup
var wg sync.WaitGroup

//Here's our User struct
type User struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"` //This was formally a []byte but we are changing our code to fit the database better
	First    string `json:"First"`
	Last     string `json:"Last"`
	Role     string `json:"Role"`
	UserID   int    `json:"UserID"`
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
	HotDogType string `json:"HotDogType"`
	Condiment  string `json:"Condiment"`
	Calories   int    `json:"Calories"`
	Name       string `json:"Name"`
	UserID     int    `json:"UserID"` //User WHOMST this hotDog belongs to
}

type Hamburger struct {
	BurgerType string `json:"BurgerType"`
	Condiment  string `json:"Condiment"`
	Calories   int    `json:"Calories"`
	Name       string `json:"Name"`
	UserID     int    `json:"UserID"` //User WHOMST this hotDog belongs to
}

//mySQL database declarations
var db *sql.DB
var err error

//Loading our templates in for ParseGlob: https://github.com/gobuffalo/packr/issues/16
var templatesBox = packr.New("templates", "./static")

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
	defer db.Close()

	err = db.Ping()
	check(err)

	fmt.Println("Test string.")
	//Launch our creations into their own goroutine
	//https://www.udemy.com/course/learn-how-to-code/learn/lecture/11922316#overview
	fmt.Printf("Launching go routines...\n")
	fmt.Printf("OS: %v\n", runtime.GOOS)
	fmt.Printf("ARCH: %v\n", runtime.GOARCH)
	fmt.Printf("CPUs: %v\n", runtime.NumCPU())
	wg.Add(16) //Need to add our wait groups for the program(should be three with main)
	go discardFood()
	go userCreator()
	go swearUserRemoverHDog()
	go swearUserRemoverHam()
	fmt.Printf("Number of goRoutines: %v\n", runtime.NumGoroutine())
	//Need to tell our main program to wait for goroutines
	wg.Wait()
}

//Check errors in our mySQL errors
func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//Some stuff for logging
func logHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Package main, son")
	fmt.Fprint(w, "package main, son.")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logrus.Infof("uri: %v\n", req.RequestURI)
		next.ServeHTTP(w, req)
	})
}

func userCreator() {
	fmt.Println("Making random Users.")

	url := "https://api.namefake.com"
	method := "GET"
	//Make 5 Users
	for v := 0; v < 5; v++ {
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
		newUser := User{
			UserName: theNewUser.Username,
			Password: randomPassword(theNewUser.Password),
			First:    fName,
			Last:     lName,
			Role:     randomRole(),
			UserID:   randomID(),
		}
		go insertUser(newUser) //User inserted
		//Give User some food
		fmt.Printf("Giving this User,(#%v) some food: %v\n", v+1, newUser.UserID)
		go giveRandomFood(newUser.UserID)
	}
	wg.Done()
}
