package main

import (
	"fmt"
)

/*********

WARNING! THIS CODE CONTAINS DEROGATORY TERMS, RACIAL/ETHNIC/SEXUAL SLURS,
AND OTHER OFFENSIVE CONTENT. THE PURPOSE IS TO REMOVE THIS CONTENT OFF OF
MY PLATFORM. IF ANY OF THIS CONTENT OFFENDS YOU, I APOLOGIZE; PLEASE STAY OFF
OF THIS PAGE!!!

*******/

/* DEFINED SLURS */
var slurs []string = []string{"penis", "vagina", "dick", "cunt", "asshole", "fag", "faggot",
	"nigglet", "nigger", "beaner", "wetback", "wet back", "chink", "tranny", "bitch", "slut",
	"whore", "fuck", "damn",
	"shit", "piss", "cum", "jizz"}

func swearUserRemoverHDog() {
	//Build the query
	var theQuery string
	theQuery = buildQueryHDog()
	//Run the query
	delDog, err := db.Prepare(theQuery)
	check(err)

	r, err := delDog.Exec()
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Printf("%v\n", n)
	//Print log files
	logWriter("Done removing hotdog bad words for SQL.")

	//Delete Records for Mongo
	wg.Add(1)
	go foodDeleteMongo(1, slurs) //1 for deleting Hotdogs
	logWriter("Done removing hotdog bad words for Mongo.")
	wg.Done() //For Go Routines
}

func buildQueryHDog() string {
	//build theQuery to return to 'swearUserRemover'
	slurQuery := ""
	slurQuery = slurQuery + "DELETE FROM hot_dogs WHERE 1=1" + "\n"
	for j := 0; j < len(slurs); j++ {
		//first pass is for the 'and' statement
		if j == 0 {
			slurQuery = slurQuery + "AND lower(hot_dogs.TYPE) LIKE " + "'%" + slurs[j] + "%'" + "\n" +
				"OR lower(hot_dogs.CONDIMENT) LIKE " + "'%" + slurs[j] + "%'" + "\n" +
				"OR lower(hot_dogs.NAME) LIKE " + "'%" + slurs[j] + "%'" + "\n"
		} else {
			slurQuery = slurQuery + "OR lower(hot_dogs.TYPE) LIKE " + "'%" + slurs[j] + "%'" + "\n" +
				"OR lower(hot_dogs.CONDIMENT) LIKE " + "'%" + slurs[j] + "%'" + "\n" +
				"OR lower(hot_dogs.NAME) LIKE " + "'%" + slurs[j] + "%'" + "\n"
		}
	}

	return slurQuery
}

func swearUserRemoverHam() {
	//Build the query
	var theQuery string
	theQuery = buildQueryHam()
	//Run the query
	delHam, err := db.Prepare(theQuery)
	check(err)

	r, err := delHam.Exec()
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Printf("%v\n", n)
	//Print log info
	logWriter("Done removing Hamburger bad words.")

	//Delete Records for Mongo
	wg.Add(1)
	go foodDeleteMongo(2, slurs) //2 for deleting Hamburgers
	logWriter("Done removing Hamburger bad words for Mongo.")
	wg.Done() //For GoRoutines
}

func buildQueryHam() string {
	//build theQuery to return to 'swearUserRemover'
	slurQuery := ""
	slurQuery = slurQuery + "DELETE FROM hamburgers WHERE 1=1" + "\n"
	for j := 0; j < len(slurs); j++ {
		//first pass is for the 'and' statement
		if j == 0 {
			slurQuery = slurQuery + "AND lower(hamburgers.TYPE) LIKE " + "'%" + slurs[j] + "%'" + "\n" +
				"OR lower(hamburgers.CONDIMENT) LIKE " + "'%" + slurs[j] + "%'" + "\n" +
				"OR lower(hamburgers.NAME) LIKE " + "'%" + slurs[j] + "%'" + "\n"
		} else {
			slurQuery = slurQuery + "OR lower(hamburgers.TYPE) LIKE " + "'%" + slurs[j] + "%'" + "\n" +
				"OR lower(hamburgers.CONDIMENT) LIKE " + "'%" + slurs[j] + "%'" + "\n" +
				"OR lower(hamburgers.NAME) LIKE " + "'%" + slurs[j] + "%'" + "\n"
		}
	}

	return slurQuery
}
