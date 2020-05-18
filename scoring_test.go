package main

import "testing"

func TestGetId(t *testing.T) {
	if LookupAthlete("david hanley", 47,"M",false).id != 1 {
		t.Error("first ID not one")
	}

	if LookupAthlete("david hanley", 47,"M",false).id != 1 {
		t.Error("should be same dave")
	}

	if LookupAthlete("david hanley", 48,"M",false).id != 1 {
		t.Error("should still be same dave")
	}

	if LookupAthlete("david hanley", 46,"M",false).id != 1 {
		t.Error("should still be same dave")
	}

	if LookupAthlete("david hanley", 0,"M",false).id != 1 {
		t.Error("should still be same dave")
	}

	if LookupAthlete("david hanley", 25,"M",false).id != 2 {
		t.Error("should be a new dave now")
	}

	if LookupAthlete("david hanley", 25,"M",false).id != 2 {
		t.Error("but still that same one")
	}

	if LookupAthlete("erin brand", 0,"F",false).id != 3 {
		t.Error("young erin should be 3")
	}

	if LookupAthlete("erin brand", 50,"F",false).id != 3 {
		t.Error("in-her-prime erin should be 3")
	}

	if athleteDb["ERIN BRAND"][0].age != 50 {
		t.Error("erin should be updated")
	}


}
