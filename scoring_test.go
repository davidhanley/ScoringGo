package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSAF(t *testing.T) {
	s, f := SAF("M")
	if s != "M" || f != false {
		t.Error("first wrong")
	}

	s, f = SAF("*M")
	if s != "M" || f != true {
		t.Error("second wrong")
	}


}

func TestGetId(t *testing.T) {
	db := makeAthleteDB()

	if lookupAthlete("david hanley", 47, "M", false, db).id != 1 {
		t.Error("first ID not one")
	}

	if lookupAthlete("david hanley", 47, "M", false, db).id != 1 {
		t.Error("should be same dave")
	}

	if lookupAthlete("david hanley", 48, "M", false, db).id != 1 {
		t.Error("should still be same dave")
	}

	if lookupAthlete("david hanley", 46, "M", false, db).id != 1 {
		t.Error("should still be same dave")
	}

	if lookupAthlete("david hanley", 0, "M", false, db).id != 1 {
		t.Error("should still be same dave")
	}

	if lookupAthlete("david hanley", 25, "M", false, db).id != 2 {
		t.Error("should be a new dave now")
	}

	if lookupAthlete("david hanley", 25, "M", false, db).id != 2 {
		t.Error("but still that same one")
	}

	if lookupAthlete("erin brand", 0, "F", false, db).id != 3 {
		t.Error("young erin should be 3")
	}

	if lookupAthlete("erin brand", 50, "F", false, db).id != 3 {
		t.Error("in-her-prime erin should be 3")
	}

	if db.db["ERIN BRAND"][0].age != 50 {
		t.Error("erin should be updated")
	}

	l := strings.Split(",WAI CHING SOH,27,*M,\"KUALA LUMPUR MALAYSIA,\",10:46",",")

    a := athleteFromLine(l,db)

    if ( a.foreign != true || a.sex!= "M") {
		t.Error("waiching not parsed")
	}
}

func TestNameDedupe(t *testing.T) {
	db := makeAthleteDB()

	if lookupAthlete("JOSH DUNCAN", 42, "M", false, db).id != 1 {
		t.Error("first ID not one")
	}

	if lookupAthlete("JOSHUA DUNCAN", 42, "M", false, db).id != 1 {
		t.Error("should be same josh")
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

/*
func TestLoadRace(t *testing.T) {
	db := makeAthleteDB()

	races := make([]*Race, 0)

	raceDate, _ := time.Parse(layoutISO, "2020-04-01")
	races = loadARace("data/2020-scale-the-strat.csv", races, db, raceDate)

	categoryMap := computeCategories(races)

	race := races[0]

	if race.name != "2020 Scale The Strat" {
		t.Error("scale the strat name is wrong")
	}

	if race.points != 350 {
		t.Error("strat points are wrong")
	}

	//check the USA only results
	overallUSA := getCategory(categoryMap, "M", US_ONLY, 0)

	winnahUSA := overallUSA.sortedAthletes[0]


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

	overall := getCategory(categoryMap, "M", ALL, 0)

	winnahOverall := overall.sortedAthletes[0]

	if winnahOverall.points != 350.0 {
		t.Error("wrong winner points")
	}

	daveAAP := overall.sortedAthletes[27]

	if daveAAP.points != 54.6875 {
		t.Error(fmt.Sprintf("dave points wrong: %f", daveAAP.points))
	}

} */

func TestLoadRace2(t *testing.T) {
	db := makeAthleteDB()

	races := make([]*Race, 0)

	scoringDate, _ := time.Parse(layoutISO, "2022-02-08")

	races = loadARace("data/2021-esbru.csv", races, db, scoringDate)

	categoryMap := computeCategories(races)

	race := races[0]

	if race.name != "2021 ESBRU" {
		t.Error("ESBRU name is wrong: " + race.name)
	}

	//verify the de-duping code is right

	category := getCategory(categoryMap, "M", ALL, 0)

	for a := 0 ; a < 20 ; a ++ {
		r := category.sortedAthletes[a]
		fmt.Print(r.athlete.name)
	}
	waiching := category.results["WAI CHING SOH"]

	if len(waiching) != 1 {
		t.Error("Waiching not present")
	}
}
