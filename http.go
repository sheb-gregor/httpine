package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	KeyRequestSeparator = "###"
	KeyMethodGet        = "GET"
	KeyMethodHead       = "HEAD"
	KeyMethodPost       = "POST"
	KeyMethodPut        = "PUT"
	KeyMethodPatch      = "PATCH"
	KeyMethodDelete     = "DELETE"
	KeyMethodOptions    = "OPTIONS"
	KeyMethodConnect    = "CONNECT"
	KeyMethodTrace      = "TRACE"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
}

type httpParser struct {
	env map[string]interface{}
}

func (parser *httpParser) ParseRequests(scanner *bufio.Scanner) ([]Request, error) {
	var r *Request
	var look4Body, look4Headers bool
	var requests []Request
	_ = look4Body

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			continue
		}

		if strings.HasPrefix(line, KeyRequestSeparator) {
			if r != nil {
				requests = append(requests, *r)
				r = nil
				look4Body = false
				look4Headers = false
			}
			continue
		}

		if line[0] == '#' {
			continue
		}

		words := strings.Fields(line)
		if len(words) < 2 {
			continue
		}

		switch words[0] {
		case KeyMethodGet, KeyMethodDelete, KeyMethodHead, KeyMethodOptions:
			look4Body = false
			look4Headers = true
			r = &Request{
				Method:  words[0],
				URL:     parser.fillSubstitutions(words[1]),
				Headers: map[string]string{},
			}
			continue

		case KeyMethodPost, KeyMethodPatch, KeyMethodPut:
			look4Body = true
			look4Headers = true
			r = &Request{
				Method:  words[0],
				URL:     parser.fillSubstitutions(words[1]),
				Headers: map[string]string{},
			}
			continue
		}

		if r != nil && look4Headers && strings.HasSuffix(words[0], ":") {
			header := strings.TrimSuffix(words[0], ":")
			r.Headers[header] = parser.fillSubstitutions(strings.Join(words[1:], " "))
			continue
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (parser *httpParser) fillSubstitutions(line string) string {
	matches := regexp.MustCompile(`\{\{[a-zA-Z0-9\-_]+\}\}`).FindAllString(line, -1)

	for _, match := range matches {
		key := strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}")
		val, ok := parser.env[key]
		if !ok {
			continue
		}

		line = strings.ReplaceAll(line, match, fmt.Sprintf("%+v", val))
	}

	matches = regexp.MustCompile(`\{\{[$uuid|$timestamp|$randomInt]+\}\}`).FindAllString(line, -1)
	for _, match := range matches {
		key := strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}")
		valGen, ok := parser.dynamicVariables()[key]
		if !ok {
			continue
		}

		val := valGen()
		line = strings.ReplaceAll(line, match, fmt.Sprintf("%+v", val))
	}

	return line
}

func (parser *httpParser) dynamicVariables() map[string]func() string {
	return map[string]func() string{
		"$uuid":      func() string { return uuid.New().String() },
		"$timestamp": func() string { return fmt.Sprintf("%d", time.Now().Unix()) },
		"$randomInt": func() string {
			r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
			val := r1.Intn(1000)
			return fmt.Sprintf("%d", val)
		},
	}
}
