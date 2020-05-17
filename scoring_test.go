package scoring

import "testing"

func TestGetId(t *testing.T) {
	if GetId("david hanley", 47) != 1 {
		t.Error("first ID not one")
	}

	if GetId("david hanley", 47) != 1 {
		t.Error("should be same dave")
	}

	if GetId("david hanley", 48) != 1 {
		t.Error("should still be same dave")
	}

	if GetId("david hanley", 46) != 1 {
		t.Error("should still be same dave")
	}

	if GetId("david hanley", 0) != 1 {
		t.Error("should still be same dave")
	}

	if GetId("david hanley", 25) != 2 {
		t.Error("should be a new dave now")
	}

	if GetId("david hanley", 25) != 2 {
		t.Error("but still that same one")
	}

	if GetId("erin brand", 0) != 3 {
		t.Error("young erin should be 3")
	}

	if GetId("erin brand", 50) != 3 {
		t.Error("in her prime erin should be 3")
	}

	if athlete_db["erin brand"][0].age != 50 {
		t.Error("erin should be updated")
	}


}
