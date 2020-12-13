package util

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Shader represents an OpenGL shader
type Shader struct {
	Type uint32
	ID   uint32
}

// NewShader creates a new OpenGL shader of the specified type
func NewShader(t uint32, source string) *Shader {
	s := &Shader{
		ID:   gl.CreateShader(t),
		Type: t,
	}

	csources, free := gl.Strs(source + "\x00")
	defer free()
	gl.ShaderSource(s.ID, 1, csources, nil)

	return s
}

// NewShaderTemplate creates a new OpenGL shader from source pre-processed as a
// Go template
func NewShaderTemplate(t uint32, source string, tplData interface{}) (*Shader, error) {
	tpl, err := template.New("anonymous").Parse(source)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	sourceBuf := &bytes.Buffer{}
	if err := tpl.Execute(sourceBuf, tplData); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return NewShader(t, sourceBuf.String()), nil
}

// NewShaderFile creates a new OpenGL shader from a source file
func NewShaderFile(t uint32, sourceFile string) (*Shader, error) {
	data, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file %v: %w", sourceFile, err)
	}

	return NewShader(t, string(data)), nil
}

// NewShaderTemplateFile creates a new OpenGL shader from a source file
// pre-processed as a Go template
func NewShaderTemplateFile(t uint32, sourceFile string, tplData interface{}) (*Shader, error) {
	sourceData, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file %v: %w", sourceFile, err)
	}

	return NewShaderTemplate(t, string(sourceData), tplData)
}

// Compile compiles the OpenGL shader
func (s *Shader) Compile() error {
	gl.CompileShader(s.ID)

	var status int32
	gl.GetShaderiv(s.ID, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var len int32
		gl.GetShaderiv(s.ID, gl.INFO_LOG_LENGTH, &len)

		log := strings.Repeat("\x00", int(len+1))
		gl.GetShaderInfoLog(s.ID, len, nil, gl.Str(log))

		return errors.New(log)
	}

	return nil
}
