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
	mockTournaments := map[int]string{
		10469768: "Melty Blood: Type Lumina",
	}
	expectedResult := map[string][]int{
		"Marc":       {158464100},
		"KosherSalt": {158464107},
		"Bernstein":  {158464116},
		"test":       {158464118},
		"test2":      {158464119},
		"test3":      {158464121},
		"test4":      {158464124},
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

	mockMultipleTournaments := map[int]string{
		10469768: "Melty Blood: Type Lumina",
		10469769: "Melty Blood: Type Lumina",
	}

	expectedResult2expectedResult2 := map[string][]int{
		"Marc":       {158464100, 158464100},
		"KosherSalt": {158464107, 158464107},
		"Bernstein":  {158464116, 158464116},
		"test":       {158464118, 158464118},
		"test2":      {158464119, 158464119},
		"test3":      {158464121, 158464121},
		"test4":      {158464124, 158464124},
	}

	// testData2, err := os.ReadFile("./test-data/testParticipantsData2.json")
	// if err != nil {
	// 	t.Errorf("Failed to read the test file\n%v\n", err)
	// }
	t.Run("Get participants from multiple tournaments", func(t *testing.T) {
		participantsResults, err := getParticipants(mockMultipleTournaments, mockClient)
		if err != nil {
			t.Errorf("getParticipants failed\n%v\n", err)
		}
		if !reflect.DeepEqual(expectedResult2expectedResult2, participantsResults) {
			t.Fatalf("Participants list did not come back the same. Expected=%v, got=%v\n", expectedResult2expectedResult2, participantsResults)
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
