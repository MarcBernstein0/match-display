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
	if err = errorhandling.HandleError("could not create template", err); err != nil {
		return "", err
	}
	var res bytes.Buffer
	err = t.Execute(&res, matches)
	if err = errorhandling.HandleError("could not parse template file", err); err != nil {
		return "", err
	}
	return res.String(), nil
}
