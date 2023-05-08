/*
A simple go script for doing phonelookup at gulesider.no at the commandline without using API.
Takes an 8 digit norwegian phonenumber as commandline argument.
..will probably stop working with the slightest change in the gulesider.no webpage.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func checkErr(e error) {
	if e != nil {
		fmt.Println("Error:", e.Error())
		os.Exit(1)
	}
}

func exit() {
	fmt.Println("Specify an 8 digit phonenumber to search for")
	os.Exit(0)
}

func isEightDigits(s string) bool {
	if len(s) != 8 {
		return false
	}
	_, err := strconv.Atoi(s)
	return err == nil
}

func parseArgs(a []string) string {
	var nr string
	if len(a) == 2 {
		if isEightDigits(a[1]) {
			nr = a[1]
		} else {
			exit()
		}
	} else {
		exit()
	}
	return nr
}

type Person struct {
	FirstName    string `json:"firstName"`
	MiddleName   string `json:"middleName"`
	LastName     string `json:"lastName"`
	StreetName   string `json:"streetName"`
	StreetNumber string `json:"streetNumber"`
	PostalCode   string `json:"postalCode"`
	PostalArea   string `json:"postalArea"`
}

func printPerson(t string, p Person) {
	fmt.Println("Tlf:", t)
	if p.FirstName != "" {
		fmt.Print(p.FirstName + " ")
	}
	if p.MiddleName != "" {
		fmt.Print(p.MiddleName + " ")
	}
	if p.LastName != "" {
		fmt.Println(p.LastName)
	}
	if p.StreetName != "" {
		fmt.Print(p.StreetName + " ")
	}
	if p.StreetNumber != "" {
		fmt.Println(p.StreetNumber)
	}
	if p.PostalCode != "" {
		fmt.Print(p.PostalCode + " ")
	}
	if p.PostalArea != "" {
		fmt.Println(p.PostalArea)
	}
}

func makePerson(xb []byte) Person {
	name := findName(xb)
	adr := findAddress(xb)
	sp := name + adr
	var p Person
	err := json.Unmarshal([]byte(sp), &p)
	checkErr(err)
	return p
}

func findName(xb []byte) string {
	bodyStr := string(xb)
	firstSplit := "name\":{\"firstName"
	secondSplit := "},\"phones\":[{"
	bodyParts := strings.Split(bodyStr, firstSplit)
	if len(bodyParts) == 1 {
		fmt.Println("No results found")
		os.Exit(0)
	}
	bp2 := strings.Split(bodyParts[1], secondSplit)
	bp := bp2[0]
	return "{\"firstName" + bp + ","
}

func findAddress(xb []byte) string {
	bodyStr := string(xb)
	firstSplit := "addresses\":[{"
	secondSplit := ",\"municipality\":"
	bodyParts := strings.Split(bodyStr, firstSplit)
	if len(bodyParts) == 1 {
		fmt.Println("No results found")
		os.Exit(0)
	}
	bp2 := strings.Split(bodyParts[1], secondSplit)
	return bp2[0] + "}"
}

func main() {
	nr := parseArgs(os.Args)
	// search string "https://www.gulesider.no/95925407/personer"
	searchUrl := "https://www.gulesider.no/" + nr + "/personer"
	search, err := http.Get(searchUrl)
	checkErr(err)

	body, err := io.ReadAll(search.Body)
	checkErr(err)
	p := makePerson(body)
	printPerson(nr, p)

}
