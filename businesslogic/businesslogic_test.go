package businesslogic

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestGetTournaments(t *testing.T) {
	testData, err := ioutil.ReadFile("./test-data/testTournamentData.json")
	if err != nil {
		t.Errorf("Failed to read the test file\n%v\n", err)
	}
	var data []map[string]map[string]interface{}
	if err = json.Unmarshal([]byte(testData), &data); err != nil {
		t.Errorf("Failed to unmarshal json data\n%v\n", err)
	}
	expectedResult := map[int]string{
		3953832:  "Guilty Gear -Strive-",
		10469768: "Melty Blood: Type Lumina",
	}
	mapResult, err := getTournaments()
	if err != nil {
		t.Errorf("getTournaments failed\n%v\n", err)
	}
	if !reflect.DeepEqual(expectedResult, mapResult) {
		t.Fatalf("Tournamet list did not come back the same. Expected=%v, got=%v\n", expectedResult, mapResult)
	}
}
