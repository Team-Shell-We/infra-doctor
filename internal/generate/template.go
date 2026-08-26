package generate

import (
	"bytes"
	"io/fs"
	"path"
	"text/template"
)

func RenderTemplate(
	files fs.FS,
	templatePath string,
	data any,
) ([]byte, error) {
	tmpl, err := template.New("generated").
		Option("missingkey=error").
		ParseFS(files, templatePath)

	if err != nil {
		return nil, err
	}

	var output bytes.Buffer

	err = tmpl.ExecuteTemplate(
		&output,
		path.Base(templatePath),
		data,
	)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
