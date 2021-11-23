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

	expectedResult := &tournaments{
		tournamentList: map[int]tournament{
			3953832: {
				tournamentID:       3953832,
				tournamentGame:     "Guilty Gear -Strive-",
				participantsByName: make(map[string]int),
				participantsByID:   make(map[int]string),
			},
			10469768: {
				tournamentID:       10469768,
				tournamentGame:     "Melty Blood: Type Lumina",
				participantsByName: make(map[string]int),
				participantsByID:   make(map[int]string),
			},
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
		t.Fatalf("Tournamet list did not come back the same. Expected=%+v, got=%+v\n", expectedResult, tournamentsResult)
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
	mockTournaments := tournaments{
		tournamentList: map[int]tournament{
			10469768: {
				tournamentID:       10469768,
				tournamentGame:     "Melty Blood: Type Lumina",
				participantsByName: make(map[string]int),
				participantsByID:   make(map[int]string),
			},
		},
	}

	expectedResult := tournaments{
		map[int]tournament{
			10469768: {
				tournamentID:   10469768,
				tournamentGame: "Melty Blood: Type Lumina",
				participantsByName: map[string]int{
					"Marc":       158464100,
					"KosherSalt": 158464107,
					"Bernstein":  158464116,
					"test":       158464118,
					"test2":      158464119,
					"test3":      158464121,
					"test4":      158464124,
				},
				participantsByID: map[int]string{
					158464100: "Marc",
					158464107: "KosherSalt",
					158464116: "Bernstein",
					158464118: "test",
					158464119: "test2",
					158464121: "test3",
					158464124: "test4",
				},
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
		err := mockTournaments.getParticipants(mockClient)
		if err != nil {
			t.Errorf("getParticipants failed\n%v\n", err)
		}
		if !reflect.DeepEqual(expectedResult, mockTournaments) {
			t.Fatalf("Participants list did not come back the same. Expected=%v, got=%v\n", expectedResult, mockTournaments)
		}
	})

	mockMultipleTournaments := tournaments{
		tournamentList: map[int]tournament{
			10469768: {
				tournamentID:       10469768,
				tournamentGame:     "Melty Blood: Type Lumina",
				participantsByName: make(map[string]int),
				participantsByID:   make(map[int]string),
			},
			10469769: {
				tournamentID:       10469769,
				tournamentGame:     "Melty Blood: Type Lumina",
				participantsByName: make(map[string]int),
				participantsByID:   make(map[int]string),
			},
		},
	}

	expectedResultMultipleTournaments := tournaments{
		tournamentList: map[int]tournament{
			10469768: {
				tournamentID:   10469768,
				tournamentGame: "Melty Blood: Type Lumina",
				participantsByName: map[string]int{
					"Marc":       158464100,
					"KosherSalt": 158464107,
					"Bernstein":  158464116,
					"test":       158464118,
					"test2":      158464119,
					"test3":      158464121,
					"test4":      158464124,
				},
				participantsByID: map[int]string{
					158464100: "Marc",
					158464107: "KosherSalt",
					158464116: "Bernstein",
					158464118: "test",
					158464119: "test2",
					158464121: "test3",
					158464124: "test4",
				},
			},
			10469769: {
				tournamentID:       10469769,
				tournamentGame:     "Melty Blood: Type Lumina",
				participantsByName: make(map[string]int),
				participantsByID:   make(map[int]string),
			},
		},
	}

	t.Run("Get participants from multiple tournaments", func(t *testing.T) {
		err := mockMultipleTournaments.getParticipants(mockClient)
		if err != nil {
			t.Errorf("getParticipants failed\n%v\n", err)
		}
		if !reflect.DeepEqual(expectedResultMultipleTournaments, mockMultipleTournaments) {
			t.Fatalf("Participants list did not come back the same. Expected=%v, got=%v\n", expectedResultMultipleTournaments, mockMultipleTournaments)
		}
	})

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("Testing error failure")
	}
	t.Run("Get paritipants error occurs", func(t *testing.T) {
		err := mockTournaments.getParticipants(mockClient)
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
	testData, err := os.ReadFile("./test-data/testMatchData.json")
	if err != nil {
		t.Errorf("Failed to read the test file\n%v\n", err)
	}
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(testData)),
		}, nil
	}

	mockTournaments := tournaments{
		map[int]tournament{
			10469768: {
				tournamentID:   10469768,
				tournamentGame: "Melty Blood: Type Lumina",
				participantsByName: map[string]int{
					"KosherSalt": 158464107,
					"test":       158464118,
					"test2":      158464119,
					"test4":      158464124,
				},
				participantsByID: map[int]string{
					158464107: "KosherSalt",
					158464118: "test",
					158464119: "test2",
					158464124: "test4",
				},
			},
		},
	}

	expectedResult := []match{
		{
			player1ID:          158464118,
			player1Name:        "test",
			player2ID:          158464119,
			player2Name:        "test2",
			tournamentID:       10469768,
			tournamentGamename: "Melty Blood: Type Lumina",
		},
		{
			player1ID:          158464107,
			player1Name:        "KosherSalt",
			player2ID:          158464124,
			player2Name:        "test4",
			tournamentID:       10469768,
			tournamentGamename: "Melty Blood: Type Lumina",
		},
	}
	t.Run("Get Matches Single Tournaments", func(t *testing.T) {
		result, err := mockTournaments.getMatches(mockClient)
		if err != nil {
			t.Errorf("getMatches failed\n%v\n", err)
		}
		if !reflect.DeepEqual(result, expectedResult) {
			t.Fatalf("Matches did not come back as expected. Expected: %v, got=%v\n", expectedResult, result)
		}
	})
}
