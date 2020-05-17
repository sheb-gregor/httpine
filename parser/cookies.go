package parser

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Cookie struct {
	Domain string
	Path   string
	Name   string
	Value  string
	Date   int64
}

func ParseCookies(path string) ([]Cookie, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cookies []Cookie
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == SComment {
			continue
		}

		cookie := Cookie{}
		chunks := strings.Split(line, STabSeparator)
		if len(chunks) < 5 {
			return nil, errors.New("invalid line format")
		}

		cookie.Domain = chunks[0]
		cookie.Path = chunks[1]
		cookie.Name = chunks[2]
		cookie.Value = chunks[3]
		cookie.Date, err = strconv.ParseInt(chunks[4], 10, 64)
		if err != nil {
			return nil, err
		}
		cookies = append(cookies, cookie)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cookies, nil
}
