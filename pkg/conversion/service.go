package conversion

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bluenviron/goroslib/v2/pkg/msgconv"
)

var tplService = template.Must(template.New("").Parse(
	`//autogenerated:yes
//nolint:revive,lll
package {{ .GoPkgName }}

import (
{{- range $k, $v := .Imports }}
    "{{ $k }}"
{{- end }}
)
{{ .Request }}
{{ .Response }}
type {{ .Name }} struct {
{{- if .RosPkgName }}
    msg.Package ` + "`" + `ros:"{{ .RosPkgName }}"` + "`" + `
{{- end }}
    {{ .Name }}Req
    {{ .Name }}Res
}
`))

// ImportService generates Go file from a .srv file and writes to the io.Writer.
func ImportService(path string, goPkgName string, rosPkgName string, w io.Writer) error {
	name := strings.TrimSuffix(filepath.Base(path), ".srv")

	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(buf)

	parts := strings.Split(content, "---")
	if len(parts) != 2 {
		return fmt.Errorf("definition must contain a request and a response")
	}

	reqDef, err := msgconv.ParseMessageDefinition(goPkgName, name+"Req", parts[0])
	if err != nil {
		return err
	}

	resDef, err := msgconv.ParseMessageDefinition(goPkgName, name+"Res", parts[1])
	if err != nil {
		return err
	}

	imports := make(map[string]struct{})
	for i := range reqDef.Imports {
		imports[i] = struct{}{}
	}
	for i := range resDef.Imports {
		imports[i] = struct{}{}
	}

	request, err := reqDef.Write()
	if err != nil {
		return err
	}

	response, err := resDef.Write()
	if err != nil {
		return err
	}

	return tplService.Execute(w, map[string]interface{}{
		"GoPkgName":  goPkgName,
		"RosPkgName": rosPkgName,
		"Imports":    imports,
		"Request":    request,
		"Response":   response,
		"Name":       name,
	})
}