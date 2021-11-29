package createhtml

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/MarcBernstein0/match-display/businesslogic"
)

func TestCreateHtml(t *testing.T) {
	mockMatches := businesslogic.Matches{
		MatchList: []businesslogic.Match{
			{
				Player1ID:          159704745,
				Player1Name:        "foo",
				Player2ID:          159704746,
				Player2Name:        "bar",
				TournamentID:       10536790,
				TournamentGameName: "Melty Blood: Type Lumina",
			},
			{
				Player1ID:          159704743,
				Player1Name:        "test2",
				Player2ID:          159704744,
				Player2Name:        "test3",
				TournamentID:       10536790,
				TournamentGameName: "Melty Blood: Type Lumina",
			},
			{
				Player1ID:          159681865,
				Player1Name:        "Foo",
				Player2ID:          159681866,
				Player2Name:        "Bar",
				TournamentID:       10535537,
				TournamentGameName: "Guilty Gear Xrd Rev 2",
			},
			{
				Player1ID:          159681863,
				Player1Name:        "Test2",
				Player2ID:          159681864,
				Player2Name:        "Test3",
				TournamentID:       10535537,
				TournamentGameName: "Guilty Gear Xrd Rev 2",
			},
		},
	}

	result, err := CreateHtml(&mockMatches)
	if err != nil {
		t.Fatalf("Error occured when calling CreateHtml\n%v", err)
	}
	tmp, _ := template.ParseFiles("template.html")
	var mockRes bytes.Buffer
	tmp.Execute(&mockRes, mockMatches)
	if result != mockRes.String() {
		t.Errorf("expected output does not match resulting output. Expected=%v\ngot=%v", mockRes.String(), result)
	}
	// fmt.Println(expectedResult)
	// t.Error("fail")
}
