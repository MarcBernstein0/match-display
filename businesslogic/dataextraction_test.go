package businesslogic

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"reflect"
	"sync"
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

	expectedResult := map[int]tournament{
		3953832: {
			tournamentID:   3953832,
			tournamentGame: "Guilty Gear -Strive-",
			participants:   make(map[string]int),
		},
		10469768: {
			tournamentID:   10469768,
			tournamentGame: "Melty Blood: Type Lumina",
			participants:   make(map[string]int),
		},
	}
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(testData)),
		}, nil
	}
	tournamentsResult, err := getTournaments(mockClient)
	if err != nil {
		t.Errorf("getTournaments failed\n%v\n", err)
	}
	if !reflect.DeepEqual(expectedResult, tournamentsResult) {
		t.Fatalf("Tournamet list did not come back the same. Expected=%v, got=%v\n", expectedResult, tournamentsResult)
	}
}

func TestMultipleApiCalls(t *testing.T) {
	// mock data and client
	testData, err := os.ReadFile("./test-data/testMultipleApiCalls.json")
	if err != nil {
		t.Errorf("Failed to read the test file\n%v\n", err)
	}
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(testData)),
		}, nil
	}

	// expected result
	var expectedResult []map[string]map[string]interface{}
	if err = json.Unmarshal([]byte(testData), &expectedResult); err != nil {
		t.Errorf("failed to unmarshal json data\n%v", err)
	}

	// test channels
	testResultChan := make(chan result)

	// test waitgoup
	testWaitGroup := new(sync.WaitGroup)
	testWaitGroup.Add(1)
	go challongeApiMultiCall(
		mockClient,
		"tournaments/10469768/participants",
		nil,
		testResultChan,
		testWaitGroup,
	)

	go func() {
		testWaitGroup.Wait()
		close(testResultChan)
	}()

	for res := range testResultChan {
		if res.err != nil {
			t.Errorf("getMultipleApiCalls failed\n%v\n", res.err)
		}
		if !reflect.DeepEqual(res.data, expectedResult) {
			t.Fatalf("Expected result did not match result. Expected=%v, got=%v\n", expectedResult, res)
		}

	}

}

func TestGetParticipants(t *testing.T) {
	mockTournaments := map[int]tournament{
		10469768: {
			tournamentID:   10469768,
			tournamentGame: "Melty Blood: Type Lumina",
			participants:   make(map[string]int),
		},
	}
	expectedResult := map[int]tournament{
		10469768: {
			tournamentID:   10469768,
			tournamentGame: "Melty Blood: Type Lumina",
			participants: map[string]int{
				"Marc":       158464100,
				"KosherSalt": 158464107,
				"Bernstein":  158464116,
				"test":       158464118,
				"test2":      158464119,
				"test3":      158464121,
				"test4":      158464124,
			},
		},
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

	t.Run("Get Participants from one tournament", func(t *testing.T) {
		participantsResults, err := getParticipants(mockTournaments, mockClient)
		if err != nil {
			t.Errorf("getParticipants failed\n%v\n", err)
		}
		if !reflect.DeepEqual(expectedResult, participantsResults) {
			t.Fatalf("Participants list did not come back the same. Expected=%v, got=%v\n", expectedResult, participantsResults)
		}
	})

	mockMultipleTournaments := map[int]tournament{
		10469768: {
			tournamentID:   10469768,
			tournamentGame: "Melty Blood: Type Lumina",
			participants:   make(map[string]int),
		},
		10469769: {
			tournamentID:   10469769,
			tournamentGame: "Melty Blood: Type Lumina",
			participants:   make(map[string]int),
		},
	}

	expectedResultMultipleTournaments := map[int]tournament{
		10469768: {
			tournamentID:   10469768,
			tournamentGame: "Melty Blood: Type Lumina",
			participants: map[string]int{
				"Marc":       158464100,
				"KosherSalt": 158464107,
				"Bernstein":  158464116,
				"test":       158464118,
				"test2":      158464119,
				"test3":      158464121,
				"test4":      158464124,
			},
		},
		10469769: {
			tournamentID:   10469769,
			tournamentGame: "Melty Blood: Type Lumina",
			participants:   make(map[string]int),
		},
	}

	t.Run("Get participants from multiple tournaments", func(t *testing.T) {
		participantsResults, err := getParticipants(mockMultipleTournaments, mockClient)
		if err != nil {
			t.Errorf("getParticipants failed\n%v\n", err)
		}
		if !reflect.DeepEqual(expectedResultMultipleTournaments, participantsResults) {
			t.Fatalf("Participants list did not come back the same. Expected=%v, got=%v\n", expectedResultMultipleTournaments, participantsResults)
		}
	})

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("Testing error failure")
	}
	t.Run("Get paritipants error occurs", func(t *testing.T) {
		_, err := getParticipants(mockTournaments, mockClient)
		if err != nil {
			if err.Error() != "request failed in getParticipants call.\nfailed to received response from challonge api.\nTesting error failure" {
				t.Fatalf("Error did not come back as expected. Expected='request failed in getParticipants call.\nfailed to received response from challonge api.\nTesting error failure', got=%v\n", err)
			}
		} else {
			t.Fatalf("Error came back empty\n")
		}
	})

}

func TestGetMatches(t *testing.T) {

}
