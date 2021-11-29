package createhtml

import (
	"bytes"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/MarcBernstein0/match-display/businesslogic"
	"github.com/MarcBernstein0/match-display/ulits/errorhandling"
)

func CreateHtml(matches *businesslogic.Matches) (string, error) {
	_, f, _, _ := runtime.Caller(0)
	path := filepath.Dir(f)
	t, err := template.ParseFiles(path + "/template.html")
	if ok, err := errorhandling.HandleError("could not create template", err); ok {
		return "", err
	}
	var res bytes.Buffer
	err = t.Execute(&res, matches)
	if ok, err := errorhandling.HandleError("could not parse template file", err); ok {
		return "", err
	}
	return res.String(), nil
}
