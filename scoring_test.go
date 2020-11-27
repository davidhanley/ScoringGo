package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetId(t *testing.T) {
	db := makeAthleteDB()

	if LookupAthlete("david hanley", 47, "M", false, db).id != 1 {
		t.Error("first ID not one")
	}

	if LookupAthlete("david hanley", 47, "M", false, db).id != 1 {
		t.Error("should be same dave")
	}

	if LookupAthlete("david hanley", 48, "M", false, db).id != 1 {
		t.Error("should still be same dave")
	}

	if LookupAthlete("david hanley", 46, "M", false, db).id != 1 {
		t.Error("should still be same dave")
	}

	if LookupAthlete("david hanley", 0, "M", false, db).id != 1 {
		t.Error("should still be same dave")
	}

	if LookupAthlete("david hanley", 25, "M", false, db).id != 2 {
		t.Error("should be a new dave now")
	}

	if LookupAthlete("david hanley", 25, "M", false, db).id != 2 {
		t.Error("but still that same one")
	}

	if LookupAthlete("erin brand", 0, "F", false, db).id != 3 {
		t.Error("young erin should be 3")
	}

	if LookupAthlete("erin brand", 50, "F", false, db).id != 3 {
		t.Error("in-her-prime erin should be 3")
	}

	if db.db["ERIN BRAND"][0].age != 50 {
		t.Error("erin should be updated")
	}
}

func TestFilterForeign(t *testing.T) {

	a1 := &Athlete{0, "dave", 48, "M", false, 5}
	a2 := &Athlete{1, "robero", 48, "M", true, 5}
	a3 := &Athlete{1, "gustavo", 28, "M", true, 2}

	athletes := []*Athlete{a1, a2, a3}

	race := &Race{"STS", 200, time.Now(), athletes}

	//filter out the foreign athletes with fewer than 3 races 
	filterRaceForForeignicity(race, THREE_RACE_FOREIGNERS)
	if len(race.athletes) != 2 {
		t.Error("Three race foreign filter size != 2")
	}

	//dual filter doesn't change things 
	filterRaceForForeignicity(race, THREE_RACE_FOREIGNERS)
	if len(race.athletes) != 2 {
		t.Error("Three race foreign filter size != 2")
	}

	//filter out all foreign athletes 
	filterRaceForForeignicity(race, US_ONLY)
	if len(race.athletes) != 1 {
		t.Error("US onlyfilter size !=1")
	}
}

func TestLoadRace(t *testing.T) {
	db := makeAthleteDB()

	races := make([]*Race, 0)

	raceDate, _ := time.Parse(layoutISO, "2020-04-01")
	races = loadARace("data/2020-scale-the-strat.csv", races, db, raceDate)

	computeCategories(races)

	race := races[0]

	if race.name != "2020 Scale The Strat" {
		t.Error("scale the strat name is wrong")
	}

	if race.points != 350 {
		t.Error("strat points are wrong")
	}

	//check the USA only results
	overallUSA := getCategory("M", US_ONLY, 0)

	winnahUSA := overallUSA.sortedAthletes[0]

	for a := 0; a < 50; a++ {
		ath := overallUSA.sortedAthletes[a].athlete
		fmt.Printf("%d %s %d\n", a,ath.name,ath.age)
	}

	if winnahUSA.points != 350.0 {
		t.Error("wrong winner points")
	}

	dave := overallUSA.sortedAthletes[20]
	if dave.athlete.name != "DAVID HANLEY" {
		t.Error("strat 20 isn't dave")
	}

	if dave.athlete.foreign == true {
		t.Error("Dave isn't foreign")
	}

	if dave.athlete.age != 47 {
		t.Error(fmt.Sprintf("Dave age is wrong %d", dave.athlete.age))
	}

	if dave.points != 70.0 {
		t.Error(fmt.Sprintf("Dave points are wrong (%f)", dave.points))
	}

	//now check foreign

	overall := getCategory("M", ALL, 0)

	winnahOverall := overall.sortedAthletes[0]

	/*for a := 0; a < 50; a++ {
		ath := overall.sortedAthletes[a]
		fmt.Printf("%d %s %f\n", a, ath.name, ath.points)
	}*/

	if winnahOverall.points != 350.0 {
		t.Error("wrong winner points")
	}

	/* dave := race.athletes[27]

	if dave.foreign == true {
		t.Error("dave is not foreign")
	}

	daveFirstResult := dave.raceResults[0]

	if daveFirstResult.rank != 22 {
		t.Error(fmt.Sprintf("dave rank (%d)is wrong...", daveFirstResult.rank))
	}

	if daveFirstResult.points != 67.307693 {
		t.Error(fmt.Printf("dave first points wrong: %f", daveFirstResult.points))
	}
	*/
}
