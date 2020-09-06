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
	"sort"
)

type AthleteRaceResult struct {
	race   *Race
	points float32
	rank   int
}

type Athlete struct {
	id      int
	name    string
	age     int
	sex     string
	foreign bool
}

type Race struct {
	name     string
	points   int
	date     time.Time
	athletes []*Athlete
}

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

type CategoryResult struct {
	gender          string
	age_low         int
	age_high        int
	include_foreign bool

	results map[Athlete][]*AthleteRaceResult
}

//taking a name an an age, return an athlete ID
func LookupAthlete(name string, age int, sex string, foreign bool) *Athlete {
	name = strings.ToUpper(name)
	name = strings.TrimSpace(name)

	athleteList := athleteDb[name]

	newAth := func() *Athlete { return &Athlete{newAthlete(), name, age, sex, foreign} }

	if len(athleteList) == 0 {
		na := newAth()
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
		na := newAth()
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

func loadARace(filename string) *Race {
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
		y, m, d := raceDate.Date()
		println(raceDateStr, y, m, d)
		println("skipping ", filename, raceName)
		return nil
	}

	println("Loading ", filename)

	popper()

	racePointsString := popper()
	racePointsString = strings.Split(racePointsString, ",")[0]
	racePointsString = strings.Split(racePointsString, "#")[0]

	var racePoints int

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

	return race
}

func scoreGender(race *Race, gender string, include_foreign bool, result *CategoryResult) {
	basePoints := float64(race.points * 5)
	denom := 5
	athletes := race.athletes
	for i := 0; i < len(athletes); i++ {
		athlete := *athletes[i]
		if athlete.sex == gender &&
			(athlete.foreign == false || include_foreign) &&
			(athlete.age >= result.age_low && athlete.age <= result.age_high) {
			points := basePoints / float64(denom)
			rr := AthleteRaceResult{race, float32(points), denom - 4}
			races := result.results
			if races[athlete] == nil {
				races[athlete] = make([]*AthleteRaceResult, 0)
			}
			races[athlete] = append(races[athlete], &rr)
			denom = denom + 1
		}
	}
}

func computeOverallForCategory(result *CategoryResult) {

	for _, results := range result.results {
		if len(results) > 5 {
			sort.Slice(results, func(i, j int) bool { return results[i].points > results[j].points })
			//println("")
			//println(ath.name)
			//for _,r := range results {
			//println(r.race.name,r.points)
			//}
		}

	}
}

func computeCategory(cr *CategoryResult) {
	for _, race := range races {
		scoreGender(race, cr.gender, cr.include_foreign, cr)
	}
	computeOverallForCategory(cr)
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
		race := loadARace(file)
		if race != nil {
			fmt.Printf("race: %v\n", race.name)
			fmt.Printf("date: %v\n", race.date)
			fmt.Printf("points: %d\n", race.points)
			fmt.Printf("athlete count: %v\n", len(race.athletes))
			races = append(races, race)
			aCount += len(race.athletes)
		}
	}
	//fmt.Printf("%d races and %d athletes", len(races), athleteCount)
}

func computeCategories() {
	genders := []string{"F", "M"}
	tf := []bool{true, false}
	ageRanges := [][]int{{0, 200}, {0, 9}, {10, 19}, {20, 29}, {30, 39}, {40, 49}, {50, 59}, {60, 69}, {70, 79}, {80, 200}}

	for _, gender := range genders {
		for _, foreign := range tf {
			for _, ar := range ageRanges {
				arr := make(map[Athlete][]*AthleteRaceResult, 0)
				cr := &CategoryResult{
					gender:          gender,
					age_low:         ar[0],
					age_high:        ar[1],
					include_foreign: foreign,
					results:         arr,
				}
				computeCategory(cr)
				println("---------------------------------")
				fmt.Printf("gender: %s\n", gender)
				fmt.Printf("age range: %d %d\n", ar[0], ar[1])
				fmt.Printf("Foreign: %t\n", foreign)
				fmt.Printf("count: %d\n", len(cr.results))
			}
		}
	}
}

/*func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}*/

func main() {
	scanFiles()
	computeCategories()
	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
