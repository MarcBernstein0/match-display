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
				Player1Name:        "foo",
				Player2Name:        "bar",
				Round:              1,
				TournamentGameName: "Melty Blood: Type Lumina",
			},
			{
				Player1Name:        "test2",
				Player2Name:        "test3",
				Round:              1,
				TournamentGameName: "Melty Blood: Type Lumina",
			},
			{
				Player1Name:        "Foo",
				Player2Name:        "Bar",
				Round:              1,
				TournamentGameName: "Guilty Gear Xrd Rev 2",
			},
			{
				Player1Name:        "Test2",
				Player2Name:        "Test3",
				Round:              1,
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
