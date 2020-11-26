package main

import (
	"bufio"
	//"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
	//"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//turn a CSV line into an athlete
//intended to return null if there is no gender
func athleteFromLine(line []string, db *AthleteDB) *Athlete {
	var athlete *Athlete
	if len(line) >= 4 {
		name := line[1]
		age, err := strconv.Atoi(line[2])
		if err != nil {
			age = 0
		}
		sex := line[3]
		foreign := false
		if len(sex) > 0 && sex[0] == '*' {
			foreign = true
			sex = sex[1:]
		}

		if len(sex) > 0 {
			sex = strings.ToUpper(sex)[:1]
			if sex == "F" || sex == "M" {
				athlete = LookupAthlete(name, age, sex, foreign, db)
			}
		}
	}
	return athlete
}


func loadARace(filename string, races []*Race, db *AthleteDB) []* Race{
	//fmt.Printf("LAR: %s\n",filename)

	csvfile, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer csvfile.Close()

	reader := bufio.NewReader(csvfile)

	popper := func() string {
		line, _ := reader.ReadString('\n')
		return strings.TrimSpace(line)
	}

	raceName := popper()
	layoutISO := "2006-1-2"
	raceDateStr := popper()
	raceDateStr = strings.Split(raceDateStr, ",")[0]
	raceDate, _ := time.Parse(layoutISO, raceDateStr)

	if raceDate.AddDate(1, 0, 0).Before(time.Now()) {
		//y, m, d := raceDate.Date()
		//println(raceDateStr, y, m, d)
		//println("skipping ", filename, raceName)
		return races
	}

	//println("Loading ", filename)

	popper()

	racePointsString := popper()
	racePointsString = strings.Split(racePointsString, ",")[0]
	racePointsString = strings.Split(racePointsString, "#")[0]

	racePoints := 0

	fmt.Sscanf(racePointsString, "%d", &racePoints)

	if racePoints == 0 {
		panic(fmt.Sprintf("unable to parse points string %s for race %s", racePointsString, filename))
	}

	athletes := make([]*Athlete, 0)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		record := strings.Split(line, ",")
		for i := 0; i < len(record); i++ {
			record[i] = strings.TrimSpace(record[i])
		}

		athlete := athleteFromLine(record,db )

		if athlete != nil {
			athlete.racesDone = athlete.racesDone + 1
			athletes = append(athletes, athlete)
		}
	}

	//fmt.Printf("Loaded %d athletes\n",len(athletes))

	race := &Race{raceName, racePoints, raceDate, athletes}

	races = append(races, race)

	return races
}

func scanFiles(db *AthleteDB) []*Race {
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

	races := make([]*Race, 0)

	for _, file := range files {
		races = loadARace(file, races, db)
	}

	return races
}


