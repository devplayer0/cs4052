package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-gl/mathgl/mgl32"
)

// LineBuffer is a list of vertices for a line
var LineBuffer = []mgl32.Vec3{
	{-1, -1, 0},
	{1, 1, 0},
}

// CubeVertices is a list of vertices for a cube
var CubeVertices = []mgl32.Vec3{
	{-1, 1, -1},
	{-1, -1, -1},
	{1, -1, -1},
	{1, -1, -1},
	{1, 1, -1},
	{-1, 1, -1},

	{-1, -1, 1},
	{-1, -1, -1},
	{-1, 1, -1},
	{-1, 1, -1},
	{-1, 1, 1},
	{-1, -1, 1},

	{1, -1, -1},
	{1, -1, 1},
	{1, 1, 1},
	{1, 1, 1},
	{1, 1, -1},
	{1, -1, -1},

	{-1, -1, 1},
	{-1, 1, 1},
	{1, 1, 1},
	{1, 1, 1},
	{1, -1, 1},
	{-1, -1, 1},

	{-1, 1, -1},
	{1, 1, -1},
	{1, 1, 1},
	{1, 1, 1},
	{-1, 1, 1},
	{-1, 1, -1},

	{-1, -1, -1},
	{-1, -1, 1},
	{1, -1, -1},
	{1, -1, -1},
	{-1, -1, 1},
	{1, -1, 1},
}

// TemplateFile loads a template from file and executes it
func TemplateFile(tplFile string, params interface{}) (string, error) {
	tplData, err := ioutil.ReadFile(tplFile)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %v: %w", tplFile, err)
	}
	tpl, err := template.New(tplFile).Funcs(sprig.TxtFuncMap()).Parse(string(tplData))
	if err != nil {
		return "", fmt.Errorf("failed to parse fragment shader template: %w", err)
	}
	dataBuf := &bytes.Buffer{}
	if err := tpl.Execute(dataBuf, params); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return dataBuf.String(), nil
}
