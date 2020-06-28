package main

func discardFood(){
	//Here is all of our User IDs
	var userIDS []int
	stmt := "SELECT USER_ID from users"
	row, err := db.Query(stmt)
	check(err)
	defer rows.Close()
	//Put query results into userIDS
	var anID int
	for row.Next() {
		err = row.Scan(&anID)
		check(err)
		userIDS = append(userIDS, anID)
	}
	//Collect all food IDs
	var hDogIDS []int
	var hamIDS []int
	var aFoodID int
	hDogStmt := "SELECT USER_ID FROM hot_dogs"
	hrow, err := db.Query(hDogStmt)
	check(err)
	defer hrow.Close()
	for hrow.Next() {
		err = hrow.Scan(&aFoodID)
		check(err)
		hDogIDS = append(hDogIDS, aFoodID)
	}
	hamStmt := "SELECT USER_ID FROM hamburgers"
	hamrow, err := db.Query(hamStmt)
	check(err)
	defer hamrow.Close()
	for hamrow.Next() {
		err = hamrow.Scan(&aFoodID)
		check(err)
		hamIDS = append(hamIDS, aFoodID)
	}
	wg.Add(2)
	go eliminateHDogs(userIDS, hDogIDS)
	go eliminateHam(userIDS, hamIDS)

	wg.Done()
}

func eliminateHDogs(theUserIDS []int, theHDogs []int){
	var theQuery string //A query to be built for eliminating food with no user.
	for z := 0; z < len(theHDogs); z++{
		foundUser := findID(theUserIDS, theHDogs[z])
		if foundUser == true {
			//Do nothing, it's got a user
		} else {
			//For first pass to build query
			if z == 0 {
				stringID := strconv.Itoa(theHDogs[z])
				theQuery = "DELETE FROM hamburgers WHERE 1=1\nAND" + 
				"USER_ID = " + stringID + "\n"
			} else{
				stringID := strconv.Itoa(theHDogs[z])
				theQuery = "OR " + "USER_ID = " + stringID + "\n"
			}
		}
	}
	//Run the query to remove those values
	delH, err := db.Prepare(theQuery)
	check(err)

	r, err := delH.Exec()
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Printf("Removed this many rows with no userID found for hotdogs: %v\n", n)
	wg.Done()
}

func eliminateHam(theUserIDS []int, theHams []int){
	var theQuery string //A query to be built for eliminating food with no user.
	for z := 0; z < len(theHams); z++{
		foundUser := findID(theUserIDS, theHams[z])
		if foundUser == true {
			//Do nothing, it's got a user
		} else {
			//For first pass to build query
			if z == 0 {
				stringID := strconv.Itoa(theHams[z])
				theQuery = "DELETE FROM hamburgers WHERE 1=1\nAND" + 
				"USER_ID = " + stringID + "\n"
			} else{
				stringID := strconv.Itoa(theHams[z])
				theQuery = "OR " + "USER_ID = " + stringID + "\n"
			}
		}
	}
	//Run the query to remove those values
	delHams, err := db.Prepare(theQuery)
	check(err)

	r, err := delHams.Exec()
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Printf("Removed this many rows with no userID found for hamburgers: %v\n", n)
	wg.Done()
}
//If ID is not found in user table, remove from the food table
func findID(theUserIDS []int, possibleInt int)bool{
	for i, item := range theUserIDS{
		if item == possibleInt{
			return true
		} else{
			return false
		}
	}
}