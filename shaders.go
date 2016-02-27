package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func CreateShaderProgram(vertexShader, fragShader string) uint32 {
	vertex := CompileShader(vertexShader, gl.VERTEX_SHADER)
	defer gl.DeleteShader(vertex)

	fragment := CompileShader(fragShader, gl.FRAGMENT_SHADER)
	defer gl.DeleteShader(fragment)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.BindFragDataLocation(program, 0, gl.Str("outColor\x00"))

	LinkShaderProgram(program)
	return program
}

func CompileShader(source string, shaderType uint32) (shader uint32) {
	shader = gl.CreateShader(shaderType)

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLen)

		logText := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(shader, logLen, nil, gl.Str(logText))

		panic(fmt.Sprintf("Shader compilation error:\n%v\n%v",
			logText, source))
	}

	return shader
}

func LinkShaderProgram(program uint32) {
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLen)

		logText := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(program, logLen, nil, gl.Str(logText))

		panic(fmt.Sprint("Shader program linking error:\n", logText))
	}
}
