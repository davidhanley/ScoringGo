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
	race := loadARace("data/2020-scale-the-strat.csv")

	if race.name != "2020 Scale The Strat" {
		t.Error("scale the strat name is wrong")
	}

	if race.points != 350 {
		t.Error("strat points are wrong")
	}

	resultmap := make(map[string][]*AthleteRaceResult, 0)
	sorted := make([]Athlete, 0)

	var cr = &CategoryResult{
		gender:          "F",
		age_low:         20,
		age_high:        50,
		include_foreign: true,
		results:         resultmap,
		sortedAthletes:  sorted,
	}

	scoreGender(race, "M", true, cr)

	first := race.athletes[0]

	if first.foreign == false {
		t.Error("the winner was foreign...")
	}

	if len(cr.sortedAthletes) == 0 {
		t.Error("didn't add sorted athletes")
	}
	if cr.sortedAthletes[0].points != 350.0 {
		println("actual points : " , cr.sortedAthletes[0].points)
		t.Error("the winner points are wrong...")
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
