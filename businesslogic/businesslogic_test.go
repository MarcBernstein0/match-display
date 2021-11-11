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

func TestGetParticipants(t *testing.T) {
	mockTournaments := map[int]string{
		10469768: "Melty Blood: Type Lumina",
	}
	expectedResult := map[int]string{
		158464100: "Marc",
		158464107: "KosherSalt",
		158464116: "Bernstein",
		158464118: "Test",
		158464119: "Test2",
		158464121: "Test3",
		158464124: "Test4",
	}
	testData, err := os.ReadFile("./test-data/testParticipantsData.json")
	if err != nil {
		t.Errorf("Failed to read the test file\n%v\n", err)
	}
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(testData)),
		}, nil
	}
	participantsResults, err := getParticipants(mockTournaments, mockClient)
	if err != nil {
		t.Errorf("getParticipants failed\n%v\n", err)
	}
	if !reflect.DeepEqual(expectedResult, participantsResults) {
		t.Fatalf("Participants list did not come back the same. Expected=%v, got=%v\n", expectedResult, participantsResults)
	}
}
