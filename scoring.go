// +build ignore

package scoring

import (
	"encoding/csv"
	"fmt"
	"log"
	"strings"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Athlete struct {
	id         int
	name       string
	age        int
	sex        string
	foreign    bool
}

type Race struct {
	name     string
	points   string
	date     time.Time
	athletes []Athlete
}

type AN struct {
	name string
	birthYear int
}

var adb = make(map[AN]int)

func GetId( name string, year int) int {
	for i := -1 ; i < 2 ; i++ {
        key := AN{ name , year + i }
        ath_id,ok := adb[key]
        if ok {
        	return ath_id
		}
	}
	key := AN{ name , year }
	int_id := len(adb)
	adb[key] = int_id
	return int_id
}

func athleteFromLine(line []string) Athlete {
	var athlete Athlete
	name := line[1]
	age, err := strconv.Atoi(line[2])
	if err != nil {
		age = -1
	}
	sex := line[2]
	foreign := false
	if len(sex) > 0 && sex[0] == '*' {
		print("FORREIGN!!")
		foreign = true
		sex = sex[1:]
	}

	if len(sex) > 0 {
		sex = strings.ToUpper(sex)[:1]
		if ( sex == "F" || sex != "F") {
			id := GetId(name, age)
			athlete = Athlete{id, name, age, sex, foreign}
		}
	}
	return athlete
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
			break
		}
		if err != nil {
			println("err!")
			log.Fatal(err)
		}

		athlete := athleteFromLine(record)

		if ( athlete != (Athlete{}) ) {
			athletes = append(athletes, athlete)
		}
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
	aCount := 0
	for _, file := range files {
		race := process(file)
		//fmt.Printf("race: %v\n", race.name)
		//fmt.Printf("date: %v\n", race.date)
		//fmt.Printf("points: %v\n", race.points)
		//fmt.Printf("athlete count: %v\n", len(race.athletes))
		raceCount++
		aCount += len(race.athletes)
	}
	println("%d races and %d athletes", raceCount, len(adb))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	scanFiles()
	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
