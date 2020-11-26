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
	"sort"
	"strconv"
	"time"
)

type AthleteRaceResult struct {
	ath    AthleteAndPoints
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

type AthleteAndPoints struct {
	athlete *Athlete
	points  float32
}

type Race struct {
	name     string
	points   int
	date     time.Time
	athletes []*Athlete
}

type CategoryResult struct {
	gender         string
	ageLow         int
	ageHigh        int
	includeForeign bool
	results        map[string][]*AthleteRaceResult
	sortedAthletes []AthleteAndPoints
}

type AthleteDB struct {
  db  map[string][]*Athlete
  athleteCount int
}

func makeAthleteDB() *AthleteDB {
	return &AthleteDB{ make(map[string][]*Athlete),0 }
}

// Abs returns the absolute value of x.
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

//taking a name an an age, return an athlete ID
func LookupAthlete(name string, age int, sex string, foreign bool,db *AthleteDB) *Athlete {
	name = strings.ToUpper(name)
	name = strings.TrimSpace(name)

	athleteList := db.db[name]

	newAth := func() *Athlete { return &Athlete{newAthleteID(db), name, age, sex, foreign} }

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
	fmt.Printf("LAR: %s\n",filename)

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
			athletes = append(athletes, athlete)
		}
	}

	fmt.Printf("Loaded %d athletes\n",len(athletes))

	race := &Race{raceName, racePoints, raceDate, athletes}

	races = append(races, race)

	return races
}

func scoreGender(race *Race, gender string, include_foreign bool, result *CategoryResult) {
	startFraction := 5
	basePoints := float64(race.points * startFraction)
	denom := startFraction
	athletes := race.athletes
	for i := 0; i < len(athletes); i++ {
		athlete := *athletes[i]
		if athlete.sex == gender &&
			(athlete.foreign == false || include_foreign) &&
			(athlete.age >= result.ageLow && athlete.age <= result.ageHigh) {
			points := basePoints / float64(denom)
			athleteRank := AthleteAndPoints{&athlete, 0.0}
			rr := AthleteRaceResult{athleteRank, race, float32(points), denom - 4}
			athletesRaces := result.results
			if athletesRaces == nil {
				panic("race nil")
			}

			if athletesRaces[athlete.name] == nil {
				athletesRaces[athlete.name] = make([]*AthleteRaceResult, 0)
			}
			athletesRaces[athlete.name] = append(athletesRaces[athlete.name], &rr)
			result.sortedAthletes = append(result.sortedAthletes, athleteRank)
			denom = denom + 1
		}
	}
}

func computeOverallForCategory(category *CategoryResult) {
	aresults := make([]AthleteAndPoints, 0)
	//first, compute the top five for each athlete
	for _, results := range category.results {
		sort.Slice(results, func(i, j int) bool { return results[i].points > results[j].points })
		points := float32(0.0)
		for i, r := range results {
			if i < 5 {
				points = points + r.points
			}
		}
		ath := results[0].ath
		ath.points = points
		aresults = append(aresults, ath)
	}
	//next, sort the category by top five results per athlete
	sort.Slice(aresults, func(i, j int) bool {
		return aresults[i].points > aresults[j].points
	})

	category.sortedAthletes = aresults
}

func computeCategory(cr *CategoryResult, races []*Race) {
	for _, race := range races {
		scoreGender(race, cr.gender, cr.includeForeign, cr)
	}
	computeOverallForCategory(cr)
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

func min(i int, j int) int {
	if i < j {
		return i
	} else {
		return j
	}
}

var categoryMap = make(map[string]*CategoryResult)

func categoryKey(gender string, foreign bool, agecat int) string {
	return fmt.Sprintf("%s%t%d", gender, foreign, agecat)
}

func getCategory(gender string, foreign bool, agecat int) *CategoryResult {
	return categoryMap[categoryKey(gender, foreign, agecat)]
}

func computeCategories(races []*Race) {
	genders := []string{"F", "M"}
	tf := []bool{true, false}
	ageRanges := [][]int{{0, 200}, {0, 9}, {10, 19}, {20, 29}, {30, 39}, {40, 49}, {50, 59}, {60, 69}, {70, 79}, {80, 200}}

	for _, foreign := range tf {
		for _, gender := range genders {
			for ageIndex, ar := range ageRanges {

				resultmap := make(map[string][]*AthleteRaceResult, 0)
				sorted := make([]AthleteAndPoints, 0)

				var cr = &CategoryResult{
					gender:         gender,
					ageLow:         ar[0],
					ageHigh:        ar[1],
					includeForeign: foreign,
					results:        resultmap,
					sortedAthletes: sorted,
				}
				computeCategory(cr, races)

				key := categoryKey(gender, foreign, ageIndex)

				categoryMap[key] = cr

				println("---------------------------------")
				fmt.Printf("gender: %s\n", gender)
				fmt.Printf("age range: %d %d\n", ar[0], ar[1])
				fmt.Printf("Foreign: %t\n", foreign)
				fmt.Printf("count: %d\n", len(cr.results))
				l := len(cr.sortedAthletes)
				ml := min(5, l)
				for a := 0; a < ml; a++ {
					ath := cr.sortedAthletes[a]
					athlete := ath.athlete
					fmt.Printf("%s %d %f\n", athlete.name, athlete.age, ath.points)
				}
			}
		}
	}
}

/*func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}*/

func main() {
	db := makeAthleteDB()
	races := scanFiles(db)
	computeCategories(races)
	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
