package main

import (
	"testing"
)

func TestGetId(t *testing.T) {
	if LookupAthlete("david hanley", 47, "M", false).id != 1 {
		t.Error("first ID not one")
	}

	if LookupAthlete("david hanley", 47, "M", false).id != 1 {
		t.Error("should be same dave")
	}

	if LookupAthlete("david hanley", 48, "M", false).id != 1 {
		t.Error("should still be same dave")
	}

	if LookupAthlete("david hanley", 46, "M", false).id != 1 {
		t.Error("should still be same dave")
	}

	if LookupAthlete("david hanley", 0, "M", false).id != 1 {
		t.Error("should still be same dave")
	}

	if LookupAthlete("david hanley", 25, "M", false).id != 2 {
		t.Error("should be a new dave now")
	}

	if LookupAthlete("david hanley", 25, "M", false).id != 2 {
		t.Error("but still that same one")
	}

	if LookupAthlete("erin brand", 0, "F", false).id != 3 {
		t.Error("young erin should be 3")
	}

	if LookupAthlete("erin brand", 50, "F", false).id != 3 {
		t.Error("in-her-prime erin should be 3")
	}

	if athleteDb["ERIN BRAND"][0].age != 50 {
		t.Error("erin should be updated")
	}
}

func TestLoadRace(t *testing.T) {

	races := make([]*Race, 0)

	races = loadARace("data/2020-scale-the-strat.csv", races)

	computeCategories(races)

	race := races[0]
	if race.name != "2020 Scale The Strat" {
		t.Error("scale the strat name is wrong")
	}

	if race.points != 350 {
		t.Error("strat points are wrong")
	}

	//check the USA only results
	overallUSA := getCategory("M", false, 0)

	winnahUSA := overallUSA.sortedAthletes[0]

	/*for a := 0; a < 50; a++ {
		ath := overallUSA.sortedAthletes[a]
		fmt.Printf("%d %s %f\n", a, ath.name, ath.points)
	}*/

	if winnahUSA.points != 350.0 {
		t.Error("wrong winner points")
	}

	/*dave := overallUSA.sortedAthletes[20]
	if dave.athlete.name != "DAVID HANLEY" {
		t.Error("strat 20 isn't dave")
	}

	if dave.athlete.foreign == true {
		t.Error("Dave isn't foreign")
	}

	if dave.points != 70.0 {
		t.Error(fmt.Sprintf("Dave points are wrong (%f)", dave.points))
	}*/

	//now check foreign

	overall := getCategory("M", true, 0)

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
