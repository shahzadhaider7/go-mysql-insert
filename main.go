package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// connecting with database
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/")

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Printf("Error %s occured when opening the database\n", err)
	} else {
		fmt.Println("Successfully connected to database")
	}

	// creating database here. It will not create database, if it already exists
	_, err = db.Exec("CREATE DATABASE test1")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully created Database")
	}

	// now selecting which database to USE
	_, err = db.Exec("USE test1")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully selected Database")
	}

	// creating required table in the database. It will not create the table if it already exists
	stmt, err := db.Prepare("CREATE Table Persons(firstname varchar(50), lastname varchar(50), age int, bloodgroup text);")
	if err != nil {
		fmt.Println(err.Error())
	}

	// now executing the above Prepared query here using Exec()
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully created Table")
	}
	// defer the close till after the main function has finished executing
	defer db.Close() // this is important, we should close the conenction once the work is finished

	// ------- Reading from CSV part ---------

	// opening and reading the CSV file now
	csvFile, err := os.Open("data.csv")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully opened the CSV file")
		reader := csv.NewReader(csvFile)

		for {
			line, error := reader.Read()
			if error == io.EOF {
				break // we reach end of the file, so stop reading
			} else if error != nil {
				log.Fatal(error)
			}

			// here we are printing what we have read from the csv file
			fmt.Println(line[0], line[1], line[2], line[3])

			//	the desired query for insertion : INSERT INTO `persons` (`firstname`, `lastname`, `age`, `bloodgroup`) VALUES ('shahzad', 'haider', '23', 'O+');

			stmt1, _ := db.Prepare("INSERT INTO persons (firstname, lastname, age, bloodgroup) VALUES( ? , ? , ? , ? )")

			stmt1.Exec(line[0], line[1], line[2], line[3])
			// executing the INSERT query by putting the values collected from the csv file
		}

		// Below is read from database operation
		// I have done this to make sure there exists data in the database
		// we are reading the data present in the database

		result, err := db.Query("SELECT * FROM persons") // selecting all columns from the table
		defer result.Close()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\n\nNow printing the data fetched from database:")
		fmt.Println("\nFirst name\tLast name\tAge\tBlood group")

		// now traverse and print each element fetched from the database
		for result.Next() {

			var fname string
			var lname string
			var age1 int
			var bgroup string

			err := result.Scan(&fname, &lname, &age1, &bgroup)

			if err != nil {
				log.Fatal(err)
			}

			// now printing the rows
			fmt.Print("\n")
			fmt.Print(fname, "\t\t")
			fmt.Print(lname, "\t\t")
			fmt.Print(age1, "\t")
			fmt.Print(bgroup)
		}
	}

}
