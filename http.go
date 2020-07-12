package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lancer-kit/sam"
)

const (
	KeyMethodGet     = "GET"
	KeyMethodHead    = "HEAD"
	KeyMethodPost    = "POST"
	KeyMethodPut     = "PUT"
	KeyMethodPatch   = "PATCH"
	KeyMethodDelete  = "DELETE"
	KeyMethodOptions = "OPTIONS"
	KeyMethodConnect = "CONNECT"
	KeyMethodTrace   = "TRACE"
)

const (
	TComment          = '#'
	TRequestSeparator = "###"
	THandlerStart     = "> {%"
	THandlerEnd       = "%}"
)

const (
	PSWaitRequest    sam.State = "wait_request"
	PSWaitHeaders    sam.State = "wait_headers"
	PSCollectHeaders sam.State = "collect_headers"
	PSWaitBody       sam.State = "wait_body"
	PSCollectBody    sam.State = "collect_body"
	PSWaitHandler    sam.State = "wait_handler"
	PSCollectHandler sam.State = "collect_handler"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string

	rawHandler []string
}

type HTTPFileParser struct {
	env            map[string]interface{}
	requests       []Request
	currentRequest *Request

	state sam.State
	sam.StateMachine
}

func NewHTTPParser(env map[string]interface{}) HTTPFileParser {
	machine := sam.NewStateMachine()
	workerSM, err := machine.
		AddTransitions(PSWaitRequest, PSCollectHeaders, PSCollectBody, PSCollectHandler).
		AddTransitions(PSCollectHeaders, PSCollectBody, PSCollectHandler, PSWaitRequest).
		AddTransitions(PSCollectBody, PSCollectHandler, PSWaitRequest).
		Finalize(PSWaitRequest)
	if err != nil || workerSM == nil {
		log.Fatal("init failed: ", err)
	}
	return HTTPFileParser{
		env: env,
	}
}

func (parser *HTTPFileParser) ParseFile(fileContent []string) ([]Request, error) {
	var look4Body, look4Headers, look4Handler bool
	_ = look4Body

	for lineNumber, line := range fileContent {
		println(lineNumber)

		if len(line) < 1 {
			continue
		}
		if line[0] == '#' {
			continue
		}

		switch parser.state {
		case PSWaitRequest:
			words := strings.Fields(line)
			if len(words) < 2 {
				continue
			}

			switch words[0] {
			case KeyMethodGet, KeyMethodDelete, KeyMethodHead, KeyMethodOptions:
				look4Body = false
				look4Headers = true
				parser.currentRequest = &Request{
					Method:  words[0],
					URL:     parser.fillSubstitutions(words[1]),
					Headers: map[string]string{},
				}
				continue

			case KeyMethodPost, KeyMethodPatch, KeyMethodPut:
				look4Body = true
				look4Headers = true
				parser.currentRequest = &Request{
					Method:  words[0],
					URL:     parser.fillSubstitutions(words[1]),
					Headers: map[string]string{},
				}
				continue
			}

			// case PSLook4Headers:
			// 	words := strings.Fields(line)
			// 	if len(words) < 2 {
			// 		continue
			// 	}

			if parser.currentRequest != nil && look4Headers && strings.HasSuffix(words[0], ":") {
				header := strings.TrimSuffix(words[0], ":")
				parser.currentRequest.Headers[header] = parser.fillSubstitutions(strings.Join(words[1:], " "))
				continue
			}
			// case PSLook4Body:
			// case PSLook4Handler:
			parser.currentRequest.rawHandler = append(parser.currentRequest.rawHandler, line)
			if strings.HasSuffix(line, THandlerEnd) {

			}
		case PSWaitHandler:
			if strings.HasPrefix(line, THandlerStart) {

			}
		}

		if look4Handler || (!look4Handler && strings.HasPrefix(line, THandlerStart)) {

			look4Handler = !strings.HasSuffix(line, THandlerEnd)
		}

		if strings.HasPrefix(line, TRequestSeparator) {
			if parser.currentRequest != nil {
				parser.requests = append(parser.requests, *parser.currentRequest)
				parser.currentRequest = nil
				look4Body = false
				look4Headers = false
			}
			continue
		}

	}

	return parser.requests, nil
}

func (parser *HTTPFileParser) fillSubstitutions(line string) string {
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

func (parser *HTTPFileParser) dynamicVariables() map[string]func() string {
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
