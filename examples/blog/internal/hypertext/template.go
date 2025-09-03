package hypertext

import (
	"embed"
	"text/template"
)

//go:generate muxt generate --receiver-type=App --template-data-type=templateData

//go:generate counterfeiter -generate
//counterfeiter:generate -o ../fake/app.go --fake-name App . RoutesReceiver

//go:embed *.gohtml
var source embed.FS

var templates = template.Must(template.ParseFS(source, "*.gohtml"))
