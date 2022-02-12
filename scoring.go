package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"math/big"
)

func scoreGender(race *Race, gender string, result *CategoryResult) {
	startFraction := 5
	basePoints := int64(race.points * startFraction)

	denom := int64(startFraction)
	athletes := race.athletes
	for i := 0; i < len(athletes); i++ {
		athlete := *athletes[i]
		if athlete.sex == gender &&
			(athlete.age >= result.ageLow && athlete.age <= result.ageHigh) {
			points := big.NewRat(basePoints, denom)
			athleteRank := AthleteAndPoints{&athlete, *big.NewRat(1, 1)}
			rr := AthleteRaceResult{athleteRank, race, *points, int(denom - 4)}
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

func computeRankForCategory(category *CategoryResult) {
	athletesAndPoints := make([]AthleteAndPoints, 0)
	//first, compute the top five for each athlete
	for _, results := range category.results {
		sort.Slice(results, func(i, j int) bool { return results[i].points.Cmp(&results[j].points) > 0 })
		var points *big.Rat = nil
		for i, r := range results {
			if i >= 3 {
				break
			}
			if points == nil {
				points = new (big.Rat)
				points.Set(&r.points)
			} else {
				points = points.Add(points, &r.points)
			}

		}
		athlete := results[0].ath
		athlete.points = *points
		athletesAndPoints = append(athletesAndPoints, athlete)
	}
	//next, sort the category by top five results per athlete
	sort.Slice(athletesAndPoints, func(i, j int) bool {
		return athletesAndPoints[i].points.Cmp(&athletesAndPoints[j].points) > 0
	})

	category.sortedAthletes = athletesAndPoints
}

func computeCategory( /*waitGroup *sync.WaitGroup, */ cr *CategoryResult, races []*Race) {
	//defer waitGroup.Done()

	for _, race := range races {
		scoreGender(race, cr.gender, cr)
	}
	computeRankForCategory(cr)
}

type CategoryMap map[string]*CategoryResult

func categoryKey(gender string, foreignicity Foreignicity, agecat int) string {
	return fmt.Sprintf("%s%d%d", gender, foreignicity, agecat)
}

func getCategory(categoryMap CategoryMap, gender string, foreignicity Foreignicity, agecat int) *CategoryResult {
	return categoryMap[categoryKey(gender, foreignicity, agecat)]
}

func filterRaceForForeignicity(race *Race, foreignicity Foreignicity) {
	athletes := make([]*Athlete, 0)
	for _, athlete := range race.athletes {

		if athlete.foreign {
			if foreignicity == US_ONLY ||
				(foreignicity == THREE_RACE_FOREIGNERS && athlete.racesDone < 3) {
				continue
			}
		}
		athletes = append(athletes, athlete)
	}
	//fmt.Printf("Filtered %s %d athletes to %d with %d\n", race.name, len(race.athletes), len(athletes), foreignicity)
	race.athletes = athletes
}

func filterRacesForForeignicity(races []*Race, foreignicity Foreignicity) {
	/*if foreignicity == ALL {
		return
	}*/

	for _, race := range races {
		filterRaceForForeignicity(race, foreignicity)
	}
}

func computeCategories(races []*Race) CategoryMap {
	//var waitGroup sync.WaitGroup

	categoryMap := make(map[string]*CategoryResult)

	genders := []string{"F", "M"}
	tf := []Foreignicity{ALL, THREE_RACE_FOREIGNERS, US_ONLY}

	ageRanges := [][]int{{0, 200}, {0, 9}, {10, 19}, {20, 29}, {30, 39}, {40, 49}, {50, 59}, {60, 69}, {70, 79}, {80, 200}}

	for _, foreign := range tf {
		filterRacesForForeignicity(races, foreign)
		for _, gender := range genders {
			for ageIndex, ar := range ageRanges {

				resultMap := make(map[string][]*AthleteRaceResult, 0)
				sorted := make([]AthleteAndPoints, 0)

				var categoryResult = &CategoryResult{
					gender:         gender,
					ageLow:         ar[0],
					ageHigh:        ar[1],
					includeForeign: foreign,
					results:        resultMap,
					sortedAthletes: sorted,
				}
				//because this has no side-effect other than modifying categoryResult,
				//we can run it as its own goroutine..  We just need to wait until we serve results..
				computeCategory( /*waitGroup, */ categoryResult, races)
				//waitGroup.Add(1)

				key := categoryKey(gender, foreign, ageIndex)
				categoryMap[key] = categoryResult
			}
		}
	}
	//waitGroup.Wait()
	return categoryMap
}

type TableRow struct {
	Rank   int
	Name   string
	Age    int
	Points string
	Races  []string
}

//var templ *template.Template

func makeHandler(categoryMap CategoryMap, templ *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.URL.Query()
		g := header["g"][0]
		f, _ := strconv.Atoi(header["f"][0])
		a, _ := strconv.Atoi(header["a"][0])
		rmax, _ := strconv.Atoi(header["r"][0])
		category := getCategory(categoryMap, g, Foreignicity(f), a)
		if category != nil {
			rows := make([]*TableRow, 0)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			for i, athlete := range category.sortedAthletes {
				if i > rmax {
					break
				}
				sa := make([]string, 0)
				results := category.results[athlete.athlete.name]
				for _, rr := range results {
					p, _ := rr.points.Float64()
					sa = append(sa, fmt.Sprintf("%s %.3f", rr.race.name, p))
				}
				p  := athlete.points.FloatString(3)
				r := &TableRow{i + 1, athlete.athlete.name, athlete.athlete.age, p, sa}
				rows = append(rows, r)
			}
			templ.ExecuteTemplate(w, "raceTable.html", rows)
		} else {
			fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
		}
	}
}

func main() {
	var templ *template.Template
	temp2, _ := templ.ParseGlob("raceTable.html")
	db := makeAthleteDB()
	races := scanFiles(db)
	categoryMap := computeCategories(races)
	http.HandleFunc("/scoring", makeHandler(categoryMap, temp2))
	println("ready to serve")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
