package web

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/Cyan903/c-share/pkg/log"
)

func ParseTemplate(w http.ResponseWriter, tmpl string, data any) {
	short := fmt.Sprintf("%s.tmpl.html", tmpl)
	t, err := template.New(short).ParseFiles(
		"internal/templates/components/base.tmpl.html",
		fmt.Sprintf("internal/templates/%s", short),
	)

	if err != nil {
		log.Error.Println("Could not read template", tmpl, err)
		os.Exit(1)
	}

	t.Execute(w, data)
}
