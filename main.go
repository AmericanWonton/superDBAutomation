package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Here's our User struct
type User struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"` //This was formally a []byte but we are changing our code to fit the database better
	First    string `json:"First"`
	Last     string `json:"Last"`
	Role     string `json:"Role"`
	UserID   int    `json:"UserID"`
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

//Handle all Requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/signup", signUp)
	myRouter.HandleFunc("/mainPage", mainPage)
	myRouter.HandleFunc("/signUpUserUpdated", signUpUserUpdated)
	//Database Stuff
	myRouter.HandleFunc("/deleteFood", deleteFood).Methods("POST")
	myRouter.HandleFunc("/updateFood", updateFood).Methods("POST")           //Update a certain food item
	myRouter.HandleFunc("/insertHotDog", insertHotDog).Methods("POST")       //Post a hotdog!
	myRouter.HandleFunc("/insertHamburger", insertHamburger).Methods("POST") //Post a hamburger!
	myRouter.HandleFunc("/getAllFoodUser", getAllFoodUser).Methods("POST")   //Get all foods for a User ID
	myRouter.HandleFunc("/getHotDog", getHotDog).Methods("GET")              //Get a SINGULAR hotdog
	myRouter.HandleFunc("/insertUser", insertUser).Methods("POST")           //Post a User!
	myRouter.HandleFunc("/getUsers", getUsers).Methods("GET")                //Get a Users!
	myRouter.HandleFunc("/updateUsers", updateUsers).Methods("POST")         //Get a Users!
	myRouter.HandleFunc("/deleteUsers", deleteUsers).Methods("POST")         //DELETE a Users!
	//Validation Stuff
	myRouter.HandleFunc("/checkUsername", checkUsername) //Check Username
	myRouter.HandleFunc("/loadUsernames", loadUsernames) //Loads in Usernames
	//Middleware logging
	myRouter.Handle("/", loggingMiddleware(http.HandlerFunc(logHandler)))
	//Serve our static files
	myRouter.Handle("/", http.FileServer(templatesBox))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(templatesBox)))
	log.Fatal(http.ListenAndServe(":80", myRouter))
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

	//Handle Requests
	handleRequests()
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
