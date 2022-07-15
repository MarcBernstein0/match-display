package mainlogic

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	fetch Fetch
)

func mockEndPoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	testData := `
	[
    {
        "tournament": {
            "id": 10878303,
            "name": "BP GGST 3/4 ",
            "url": "4p7r215v",
            "description": "",
            "tournament_type": "double elimination",
            "started_at": "2022-03-04T22:09:03.530-03:00",
            "completed_at": null,
            "require_score_agreement": false,
            "notify_users_when_matches_open": true,
            "created_at": "2022-03-04T19:02:33.781-03:00",
            "updated_at": "2022-03-05T02:48:04.667-03:00",
            "state": "awaiting_review",
            "open_signup": false,
            "notify_users_when_the_tournament_ends": true,
            "progress_meter": 100,
            "quick_advance": false,
            "hold_third_place_match": false,
            "pts_for_game_win": "0.0",
            "pts_for_game_tie": "0.0",
            "pts_for_match_win": "1.0",
            "pts_for_match_tie": "0.5",
            "pts_for_bye": "1.0",
            "swiss_rounds": 0,
            "private": false,
            "ranked_by": "match wins",
            "show_rounds": true,
            "hide_forum": false,
            "sequential_pairings": false,
            "accept_attachments": false,
            "rr_pts_for_game_win": "0.0",
            "rr_pts_for_game_tie": "0.0",
            "rr_pts_for_match_win": "1.0",
            "rr_pts_for_match_tie": "0.5",
            "created_by_api": false,
            "credit_capped": false,
            "category": null,
            "hide_seeds": false,
            "prediction_method": 0,
            "predictions_opened_at": null,
            "anonymous_voting": false,
            "max_predictions_per_user": 1,
            "signup_cap": null,
            "game_id": 198501,
            "participants_count": 18,
            "group_stages_enabled": false,
            "allow_participant_match_reporting": true,
            "teams": false,
            "check_in_duration": null,
            "start_at": "2022-03-04T22:00:00.000-03:00",
            "started_checking_in_at": null,
            "tie_breaks": [
                "match wins vs tied",
                "game wins",
                "points scored"
            ],
            "locked_at": null,
            "event_id": null,
            "public_predictions_before_start_time": false,
            "ranked": false,
            "grand_finals_modifier": null,
            "predict_the_losers_bracket": false,
            "spam": null,
            "ham": null,
            "rr_iterations": null,
            "tournament_registration_id": null,
            "donation_contest_enabled": null,
            "mandatory_donation": null,
            "non_elimination_tournament_data": {
                "participants_per_match": ""
            },
            "auto_assign_stations": null,
            "only_start_matches_with_stations": null,
            "registration_fee": "0.0",
            "registration_type": "free",
            "split_participants": false,
            "allowed_regions": [],
            "show_participant_country": null,
            "program_id": null,
            "program_classification_ids_allowed": null,
            "team_size_range": null,
            "toxic": null,
            "use_new_style": null,
            "optional_display_data": {
                "show_standings": "1",
                "show_announcements": true
            },
            "processing": false,
            "oauth_application_id": null,
            "description_source": "",
            "subdomain": null,
            "full_challonge_url": "https://challonge.com/4p7r215v",
            "live_image_url": "https://challonge.com/4p7r215v.svg",
            "sign_up_url": null,
            "review_before_finalizing": true,
            "accepting_predictions": false,
            "participants_locked": true,
            "game_name": "Guilty Gear -Strive-",
            "participants_swappable": false,
            "team_convertable": false,
            "group_stages_were_started": false
        }
    }
]
	`
	// fmt.Println("reading test data", string(testData))
	fmt.Fprint(w, string(testData))
}

func TestMain(m *testing.M) {
	fmt.Println("Starting mock server")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		mockEndPoint(w, r)
	}))

	fetch = New(server.URL, http.DefaultClient, time.Second)
	fmt.Println("Mock Server Running, Start Tests")
	m.Run()
}

// func TestFetchData(t *testing.T) {
// 	tt := []struct {
// 		testName  string
// 		wantData  []Tournaments
// 		wantErr   error
// 		expectErr bool
// 	}{
// 		{
// 			testName: "fetch tournaments test",
// 			wantData: []Tournaments{
// 				{
// 					Tournament: struct {
// 						Id   int    `json:"id"`
// 						Name string `json:"name"`
// 					}{
// 						Id:   10878303,
// 						Name: "BP GGST 3/4 ",
// 					},
// 				},
// 			},
// 			wantErr:   nil,
// 			expectErr: false,
// 		},
// 		{
// 			testName: "fetch participants test",
// 		},
// 	}

// 	for _, test := range tt {
// 		t.Run(test.testName, func(t *testing.T) {
// 			t.Parallel()
// 			gotData, gotErr := fetch.FetchTournamentData(context.Background(), time.Now().Local().Format("2006-01-02"))
// 			assert.Equal(t, test.wantData, gotData)
// 			// assert.Equal(t, test.wantData, gotData)
// 			if test.expectErr {
// 				assert.EqualError(t, gotErr, test.wantErr.Error(), "expected %v but got %v", test.wantErr.Error(), gotErr.Error())
// 			} else {
// 				assert.NoError(t, gotErr)
// 			}
// 		})
// 	}
// }

func TestFetchTournaments(t *testing.T) {
	testStruct := []struct {
		wantData []Tournaments
		wantErr  error
	}{
		{wantData: []Tournaments{
			{
				Tournament: struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				}{
					Id:   10878303,
					Name: "BP GGST 3/4 ",
				},
			},
		},
			wantErr: nil},
	}

	for _, test := range testStruct {
		t.Parallel()
		gotData, gotErr := fetch.FetchTournamentData(context.Background(), time.Now().Local().Format("2006-01-02"))
		assert.Equal(t, test.wantData, gotData)
		assert.NoError(t, gotErr)
	}
	//
}

func TestFetchParticipantData(t *testing.T) {
    testStruct := {
    }
}
