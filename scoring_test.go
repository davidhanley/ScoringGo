package scoring

import "testing"

func TestGetId(t *testing.T) {
	if GetId("david hanley", 47) != 0 {
		t.Error("first ID not zero")
	}

	if GetId("david hanley", 47) != 0 {
		t.Error("should be same dave")
	}

	if GetId("david hanley", 48) != 0 {
		t.Error("should still be same dave")
	}

	if GetId("david hanley", 46) != 0 {
		t.Error("should still be same dave")
	}

	if GetId("david hanley", 0) != 0 {
		t.Error("should still be same dave")
	}
}
