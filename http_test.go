package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	url := "https://{{host}}:{{port2-2}}/status/404/{{host}}"
	matches := regexp.MustCompile(`\{\{[a-zA-Z0-9\-_]+\}\}`).FindAllString(url, -1)
	assert.Len(t, matches, 3)
	assert.Equal(t, "{{host}}", matches[0])
	assert.Equal(t, "{{port2-2}}", matches[1])
	assert.Equal(t, "{{host}}", matches[2])

	line := "Test-UUID: {{$uuid}} - Test-time: {{$timestamp}} {{$randomInt}} "
	matches = regexp.MustCompile(`\{\{[$uuid|$timestamp|$randomInt]+\}\}`).FindAllString(line, -1)
	assert.Equal(t, "{{$uuid}}", matches[0])
	assert.Equal(t, "{{$timestamp}}", matches[1])
	assert.Equal(t, "{{$randomInt}}", matches[2])
}

func TestHttpParser_fillSubstitutions(t *testing.T) {
	parser := HTTPFileParser{
		env: map[string]interface{}{
			"host": "localhost",
			"port": 2020,
			"id":   1,
			"name": "alice",
		},
	}

	target := "http://localhost:2020/user/alice?for=1&from=localhost"
	result := parser.fillSubstitutions("http://{{host}}:{{port}}/user/{{name}}?for={{id}}&from={{host}}")
	assert.Equal(t, target, result)

	result = parser.fillSubstitutions("http://{{host}}:{{port}}/user/{{$uuid}}?for={{$randomInt}}&s={{$timestamp}}")
	t.Log(result)
	assert.NotEqual(t, target, result)
}
