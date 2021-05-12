package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Hard coded token for now to deal with special sample data
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I"
	fmt.Printf("Expected token: %s\n", token)

	// Test the user exists
	f, err := os.Open("../sample_data/test_user.json")
	if err != nil {
		// handle err
		panic(err.Error())
	}
	defer f.Close()
	req, err := http.NewRequest("POST", "http://localhost:1323/n3/graphql", f)
	if err != nil {
		// handle err
		panic(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		panic(err.Error())
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		// handle err
		panic(err.Error())
	}
	// TODO - check return code = 500 or better error parsing, or json
	if strings.Contains(string(b), "Failed to retrieve context") {
		fmt.Printf("User does not exist, creating\n")
		body := strings.NewReader(`userName=n3Demo&contextName=mySchool`)
		req, err := http.NewRequest("POST", "http://localhost:1323/admin/newdemocontext", body)
		if err != nil {
			// handle err
			panic(err.Error())
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
			panic(err.Error())
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			// handle err
			panic(err.Error())
		}

		fmt.Printf("User Created Success: %s\n", b)
	} else {
		fmt.Printf("User already exists\n")
	}

	// POST to X with token
	//  If not auth
	// curl -s -X POST -d 'userName=n3Demo' -d 'contextName=mySchool' \ localhost:1323/admin/newdemocontext > token.json
	//  Check token created matches - ERROR here
	// Try X again - ERROR here

	// Read CSV file
	//  post_token,URL,filename

	// Open the file
	csvfile, err := os.Open("../sample_data/index.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if record[0] == "post_token" {
			fmt.Printf("post_auth: %s => %s\n", record[2], record[1])
			f, err := os.Open(record[2])
			if err != nil {
				// handle err
				panic(err.Error())
			}
			defer f.Close()
			req, err := http.NewRequest("POST", record[1], f)
			if err != nil {
				// handle err
				panic(err.Error())
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				panic(err.Error())
			}
			defer resp.Body.Close()

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				// handle err
				panic(err.Error())
			}

			fmt.Printf("Success: %s\n", b)
		} else {
			fmt.Printf("Invalid method type in CSV, use post_token (val=%s)\n", record[0])
		}
	}

}
