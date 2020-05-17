// +build ignore

package main

import (
	"fmt"
	"log"
	//"net/http"
	//"strings"
	"io"
	"encoding/csv"
	"os"
	"time"
	"strconv"
	"path/filepath"
	"net/http"
)

type Athlete struct {
	name       string
	birth_year int
	sex        string
}

type Race struct {
	name     string
	points   string
	date     time.Time
	athletes []Athlete
}

func athleteFromLine(line []string) Athlete {
	name := line[1]
	age, err := strconv.Atoi(line[2])
	if err != nil {
		age = -1
	}
	sex = line[2]
	return Athlete{name, 2020 - age}
}

func process(fn string) Race {
	println("Loading", fn)
	csvfile, err := os.Open(fn)

	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	r := csv.NewReader(csvfile)

	popper := func() string {
		line, _ := r.Read()
		return line[0]
	}

	raceName := popper()
	layoutISO := "2006-01-02"
	raceDate, _ := time.Parse(layoutISO, popper())

	popper()

	racePoints := popper()

	athletes := make([]Athlete, 0, 500)

	for {
		r.FieldsPerRecord = 0
		record, err := r.Read()
		if err == io.EOF {
			println("eof")
			break
		}
		if err != nil {
			println("err!")
			log.Fatal(err)
		}

		athlete := athleteFromLine(record)

		athletes = append(athletes, athlete)
	}

	csvfile.Close()

	return Race{raceName, racePoints, raceDate, athletes}
}

func scanFiles() {
	var files []string

	root := "data/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	files = files[1:]

	if err != nil {
		panic(err)
	}

	raceCount := 0

	for _, file := range files {
		race := process(file)
		fmt.Printf("race: %v\n", race.name)
		fmt.Printf("date: %v\n", race.date)
		fmt.Printf("points: %v\n", race.points)
		fmt.Printf("athlete count: %v\n", len(race.athletes))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	scanFiles()
	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}

