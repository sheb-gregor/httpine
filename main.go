package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sheb-gregor/go-httpine/parser"
)

const (
	ClientEnv        = "http-client.env.json"
	ClientPrivateEnv = "http-client.private.env.json"
	ClientCookies    = "http-client.cookies"
)

func main() {
	ctx := NewClientCtx()
	err := ctx.AnalyzeDirectory("./testdata/")
	if err != nil {
		log.Fatal("directory analyze filed: ", err.Error())
	}

	err = ctx.ParseRequests("todo")
	if err != nil {
		log.Fatal("requests parsing failed: ", err.Error())
	}

	fmt.Printf("%+v\n", ctx)
}

type ClientCtx struct {
	cookies    []parser.Cookie
	env        map[string]map[string]interface{}
	privateEnv map[string]map[string]interface{}

	requestFiles []string
	requests     map[string][]Request
}

func NewClientCtx() ClientCtx {
	return ClientCtx{
		cookies:      []parser.Cookie{},
		env:          map[string]map[string]interface{}{},
		privateEnv:   map[string]map[string]interface{}{},
		requestFiles: []string{},
		requests:     map[string][]Request{},
	}

}

func (client *ClientCtx) AnalyzeDirectory(dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		switch info.Name() {
		case ClientEnv:
			log.Println("found client env...")
			raw, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "unable to read file")
			}

			if err = json.Unmarshal(raw, &client.env); err != nil {
				return errors.Wrap(err, "unable to parse cookies")
			}

		case ClientPrivateEnv:
			log.Println("found private client env...")
			raw, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "unable to read file")
			}

			if err = json.Unmarshal(raw, &client.privateEnv); err != nil {
				return errors.Wrap(err, "unable to unmarshal private env")
			}

		case ClientCookies:
			log.Println("found cookies..")
			client.cookies, err = parser.ParseCookies(path)
			if err != nil {
				return errors.Wrap(err, "unable to parse cookies")
			}
		default:
			if ext := filepath.Ext(path); ext == ".http" {
				log.Println("found http requests set", info.Name())
				client.requestFiles = append(client.requestFiles, path)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	client.mergeEnv()
	return nil
}

func (client *ClientCtx) mergeEnv() {
	for env := range client.privateEnv {
		if _, ok := client.env[env]; !ok {
			client.env[env] = map[string]interface{}{}
		}

		for key, val := range client.privateEnv[env] {
			client.env[env][key] = val
		}
	}
}

func (client *ClientCtx) ParseRequests(envName string) error {
	for _, file := range client.requestFiles {
		readFile, err := os.Open(file)

		if err != nil {
			log.Fatalf("failed to open file: %s", err)
		}

		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)
		var fileContent []string

		for fileScanner.Scan() {
			fileContent = append(fileContent, fileScanner.Text())
		}

		if err = readFile.Close(); err != nil {
			return err
		}

		// httpParser := HTTPFileParser{env: client.env[envName]}
		httpParser := NewHTTPParser(client.env[envName])
		requests, err := httpParser.ParseFile(fileContent)
		if err != nil {
			return errors.Wrap(err, "error while parsing "+file)
		}

		client.requests[file] = requests
	}

	return nil
}
