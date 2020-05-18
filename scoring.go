package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	//"regexp"
	"strings"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type RaceResult struct {
	race   Race
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
	athletes []Athlete
}

type AN struct {
	id   int
	name string
	age  int
}

var athlete_db = make(map[string][]AN)
var athlete_count = 0

// Abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func new_athlete() int {
	athlete_count = athlete_count + 1
	return athlete_count
}

//taking a name an an age, return an athlete ID
func GetId(name string, age int) int {
	name = strings.ToUpper(name)
	name = strings.TrimSpace(name)

	athleteList, ok := athlete_db[name]

	if ok == false {
		id := new_athlete()
		athlete_db[name] = []AN{{id, name, age}}
		return id
	} else {
		if age == 0 {
			return athleteList[0].id
		}
		for aid := range athleteList {
			athlete := &athleteList[aid]
			if athlete.age == 0 {
				athlete.age = age
			}
			if abs(athlete.age-age) < 2 {
				return athlete.id
			}
		}
		id := new_athlete()
		athleteList = append(athleteList, AN{id, name, age})
		athlete_db[name] = athleteList
		return id
	}
}

//turn a CSV line into an athlete
//intended to return null if there is no gender
func athleteFromLine(line []string) *Athlete {
	var athlete Athlete
	name := line[1]
	age, err := strconv.Atoi(line[2])
	if err != nil {
		age = -1
	}
	sex := line[3]
	foreign := false
	if len(sex) > 0 && sex[0] == '*' {
		print("FORREIGN!!")
		foreign = true
		sex = sex[1:]
	}

	if len(sex) > 0 {
		sex = strings.ToUpper(sex)[:1]
		if sex == "F" || sex == "M" {
			id := GetId(name, age)
			athlete = Athlete{id, name, age, sex, foreign, make([]RaceResult, 1)}
		} /*else {
			println("improper sex:", sex)
		}*/
	}
	return &athlete
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
	raceDate, _ := time.Parse(layoutISO, raceDateStr)

	//loc, _ := time.LoadLocation("UTC")

	//println("--------------")
	//println("race:", fn)
	//y,m,d := raceDate.Date()

	//fmt.Printf("%s -- %d %d %d",raceDateStr, y,m,d)

	//if y<1900 {
	//	panic("y")
	//}

	if raceDate.AddDate(1, 0, 0).Before(time.Now()) {
		//println("Skipping")
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

	athletes := make([]Athlete, 0, 500)

	csvfile.Seek(0, io.SeekStart)
	r := csv.NewReader(csvfile)
	r.Read()
	r.Read()
	r.Read()
	r.Read()
	for {
		r.FieldsPerRecord = 0
		record, err := r.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			println("err!")
			log.Fatal(err)
		}

		athlete := athleteFromLine(record)

		if athlete != nil {
			athletes = append(athletes, *athlete)
			//print(".")
		} /*else {
			println("unable to parse line :")
			for i := 0; i < len(record); i++ {
				print(record[i], ",")
			}

		} */
	}

	return &Race{raceName, racePoints, raceDate, athletes}
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
	aCount := 0
	for _, file := range files {
		race := process(file)
		if race != nil {
			fmt.Printf("race: %v\n", race.name)
			fmt.Printf("date: %v\n", race.date)
			fmt.Printf("points: %d\n", race.points)
			fmt.Printf("athlete count: %v\n", len(race.athletes))
			raceCount++
			aCount += len(race.athletes)
		}
	}
	println("%d races and %d athletes", raceCount, athlete_count)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	scanFiles()
	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
