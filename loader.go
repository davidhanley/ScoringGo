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
	"regexp"
)

type AthleteDB struct {
	db           map[string][]*Athlete
	athleteCount int
}

func makeAthleteDB() *AthleteDB {
	return &AthleteDB{make(map[string][]*Athlete), 0}
}

// Abs returns the absolute value of x. Let's add gnerics, golang!
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func newAthleteID(athleteDB *AthleteDB) int {
	athleteDB.athleteCount = athleteDB.athleteCount + 1
	return athleteDB.athleteCount
}

//taking a name and an age, return an athlete ID
func lookupAthlete(name string, age int, sex string, foreign bool, db *AthleteDB) *Athlete {
	name = translateName(strings.ToUpper(name))
	name = strings.TrimSpace(name)

	athleteList := db.db[name]

	newAth := func() *Athlete { return &Athlete{newAthleteID(db), name, age, sex, foreign, 0} }

	if len(athleteList) == 0 {
		na := newAth()
		athleteList = append(athleteList, na)
		db.db[name] = athleteList
		return na
	} else {
		if age == 0 {
			return athleteList[0]
		}
		for aid := range athleteList {
			athlete := athleteList[aid]
			if athlete.age == 0 {
				athlete.age = age
			}
			if abs(athlete.age-age) < 2 {
				return athlete
			}
		}
		na := newAth()
		athleteList = append(athleteList, na)
		db.db[name] = athleteList
		return na
	}
}

type nametrans struct {
	exp *regexp.Regexp
	to  string
}

var transTable = loadTranslation()

func loadTranslation() []nametrans {
	transFile, err := os.Open("translate.dat")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer transFile.Close()

	reader := bufio.NewReader(transFile)

	nt := make([]nametrans, 0)

	for {
		line, err := reader.ReadString('\n')
		line = strings.ToUpper(line)
		if err == io.EOF {
			break
		}
		parts := strings.Split(line, ",")

		reg, _ := regexp.Compile(parts[0])

		nt = append(nt, nametrans{reg, parts[1]})
	}
	return nt
}

func translateName(stringIn string) string {
	b := []byte(stringIn)
	for _, r := range transTable {
		if r.exp.Match(b) {
			return r.to
		}
	}
	return stringIn
}

func SAF(sex string) (string,bool) {
	foreign := false
	if len(sex) > 0 && sex[0] == '*' {
		foreign = true
		sex = sex[1:]
	}
	return sex,foreign
}

//turn a CSV line into an athlete
//intended to return null if there is no gender
func athleteFromLine(line []string, db *AthleteDB) *Athlete {
	//fmt.Print(line)
	var athlete *Athlete
	if len(line) >= 4 {
		name := line[1]
		age, err := strconv.Atoi(line[2])
		if err != nil {
			age = 0
		}
		_sex := line[3]
		sex,foreign := SAF(_sex)

		if len(sex) > 0 {
			sex = strings.ToUpper(sex)[:1]
			if sex == "F" || sex == "M" {
				athlete = lookupAthlete(name, age, sex, foreign, db)
			}
		}
	}
	//fmt.Printf("%s %s %v\n",athlete.name,athlete.sex,athlete.foreign)
	return athlete
}

var layoutISO = "2006-1-2"

func loadARace(filename string, races []*Race, db *AthleteDB, now time.Time) []*Race {
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
	raceName = strings.Split(raceName, ",")[0]
	raceDateStr := popper()
	raceDateStr = strings.Split(raceDateStr, ",")[0]
	raceDate, _ := time.Parse(layoutISO, raceDateStr)

	if raceDate.AddDate(1, 0, 0).Before(now) || raceDate.After(now) {
		//y, m, d := raceDate.Date()
		//println(raceDateStr, y, m, d)
		//println("skipping ", filename, raceName)
		return races
	}
	println("Loading ", filename)

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

	athleteMap := make(map[int]bool)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		fmt.Printf("%s\n",line)

		record := strings.Split(line, ",")
		for i := 0; i < len(record); i++ {
			record[i] = strings.TrimSpace(record[i])
		}

		athlete := athleteFromLine(record, db)

		if athlete != nil && athleteMap[athlete.id] == false {
			athleteMap[athlete.id] = true
			athlete.racesDone = athlete.racesDone + 1
			athletes = append(athletes, athlete)
		}
	}

	fmt.Printf("Loaded %d athletes\n", len(athletes))

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
		races = loadARace(file, races, db, time.Now())
	}

	return races
}
