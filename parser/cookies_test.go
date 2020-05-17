package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCookies(t *testing.T) {
	cookies, err := ParseCookies("./testdata/http-client.cookies")
	assert.Nil(t, err)
	assert.Len(t, cookies, 3)
}
