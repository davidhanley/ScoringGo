package main

import (
	"bufio"
	//"encoding/csv"
	"fmt"
	"log"
	"strings"
	"io"
	//"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type RaceResult struct {
	race   *Race
	points float32
}

type Athlete struct {
	id          int
	name        string
	age         int
	sex         string
	foreign     bool
	raceResults []RaceResult
}

type Race struct {
	name     string
	points   int
	date     time.Time
	athletes []*Athlete
}

/*type AN struct {
	id   int
	name string
	age  int
}*/

var athleteDb = make(map[string][]*Athlete)
var athleteCount = 0
var races = make([]*Race, 0)

// Abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func newAthlete() int {
	athleteCount = athleteCount + 1
	return athleteCount
}

//taking a name an an age, return an athlete ID
func LookupAthlete(name string, age int, sex string, foreign bool) *Athlete {
	name = strings.ToUpper(name)
	name = strings.TrimSpace(name)

	athleteList := athleteDb[name]

	if len(athleteList) == 0 {
		id := newAthlete()
		na := &Athlete{id, name, age, sex, foreign, make([]RaceResult, 0)}
		athleteList = append(athleteList, na)
		athleteDb[name] = athleteList
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
		id := newAthlete()
		na := &Athlete{id, name, age, sex, foreign, make([]RaceResult, 0)}
		athleteList = append(athleteList, na)
		athleteDb[name] = athleteList
		return na
	}
}

//turn a CSV line into an athlete
//intended to return null if there is no gender
func athleteFromLine(line []string) *Athlete {
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
				athlete = LookupAthlete(name, age, sex, foreign)
			}
		}
	}
	return athlete
}

func process(fn string) *Race {
	csvfile, err := os.Open(fn)
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
		y, m, d := raceDate.Date()
		println(raceDateStr, y, m, d)
		println("skipping ", fn, raceName)
		return nil
	}

	println("Loading ", fn)

	popper()

	racePointsString := popper()
	racePointsString = strings.Split(racePointsString, ",")[0]
	racePointsString = strings.Split(racePointsString, "#")[0]

	var racePoints int
	fmt.Sscanf(racePointsString, "%d", &racePoints)
	if racePoints == 0 {
		panic(fmt.Sprintf("unable to parse points string %s for race %s", racePointsString, fn))
	}

	athletes := make([]*Athlete, 0)

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}
		record := strings.Split(line,",")
		for i := 0; i < len(record) ; i++ {
			record[i] = strings.TrimSpace(record[i])
		}

		if err != nil {
			println("err!")
			log.Fatal(err)
		}

		athlete := athleteFromLine(record)

		if athlete != nil {
			athletes = append(athletes, athlete)
		}
	}

	race := &Race{raceName, racePoints, raceDate, athletes}
	scoreGender(race, "F")
	scoreGender(race, "M")
	return race
}

func scoreGender(race *Race, gender string) {
	basePoints := float64(race.points * 5)
	denom := 5
	athletes := race.athletes
	for i := 0; i < len(athletes); i++ {
		athlete := athletes[i]
		if athlete.sex == gender {
			points := basePoints / float64(denom)
			//fmt.Printf("%d %s %s %f\n", athlete.id, athlete.name, race.name, points)
			athlete.raceResults = append(athlete.raceResults, RaceResult{race, float32(points)})
			denom = denom + 1
		}
	}
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

	aCount := 0
	for _, file := range files {
		race := process(file)
		if race != nil {
			fmt.Printf("race: %v\n", race.name)
			fmt.Printf("date: %v\n", race.date)
			fmt.Printf("points: %d\n", race.points)
			fmt.Printf("athlete count: %v\n", len(race.athletes))
			races = append(races, race)
			aCount += len(race.athletes)
		}
	}
	fmt.Printf("%d races and %d athletes", len(races), athleteCount)

	println("-----------------------------------")
	jda := athleteDb["DAVID HANLEY"]
	for j := 0; j < len(jda); j++ {
		jd := jda[j]
		println(jd.name, jd.age)
		results := jd.raceResults
		for i := 0; i < len(results); i++ {
			r := results[i]
			fmt.Printf("%s %f\n", r.race.name, r.points)
		}
	}

}

/*func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}*/

func main() {
	scanFiles()
	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
