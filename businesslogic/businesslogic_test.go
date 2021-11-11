package businesslogic

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/MarcBernstein0/match-display/ulits/mocks"
)

var mockClient *mocks.MockClient

func init() {
	mockClient = &mocks.MockClient{}
}

func TestGetTournaments(t *testing.T) {
	testData, err := os.ReadFile("./test-data/testTournamentData.json")
	if err != nil {
		t.Errorf("Failed to read the test file\n%v\n", err)
	}

	expectedResult := map[int]string{
		3953832:  "Guilty Gear -Strive-",
		10469768: "Melty Blood: Type Lumina",
	}
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(testData)),
		}, nil
	}
	mapResult, err := getTournaments(mockClient)
	if err != nil {
		t.Errorf("getTournaments failed\n%v\n", err)
	}
	if !reflect.DeepEqual(expectedResult, mapResult) {
		t.Fatalf("Tournamet list did not come back the same. Expected=%v, got=%v\n", expectedResult, mapResult)
	}
}
