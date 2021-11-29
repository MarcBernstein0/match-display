package createhtml

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/MarcBernstein0/match-display/businesslogic"
	"github.com/MarcBernstein0/match-display/ulits/errorhandling"
)

func CreateHtml(matches businesslogic.Matches) (string, error) {
	t, err := template.ParseFiles("template.html")
	if ok, err := errorhandling.HandleError("could not create template", err); ok {
		return "", err
	}
	fmt.Println(t)
	var res bytes.Buffer
	err = t.Execute(&res, matches)
	if ok, err := errorhandling.HandleError("could not parse template file", err); ok {
		return "", err
	}
	return res.String(), nil
}
