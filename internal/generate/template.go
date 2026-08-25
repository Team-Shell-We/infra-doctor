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
	tmpl, err := template.New("generated"). //generated 이름의 템플릿 생성
						Option("missingkey=error").
						ParseFS(files, templatePath)

	if err != nil {
		return nil, err
	}

	var output bytes.Buffer //출력 버퍼 생성(결과를 메모리에 저장할 버퍼)

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
