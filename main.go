package main

import (
	"bufio"
	"bytes"
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
	ctx := ClientCtx{
		cookies:      []parser.Cookie{},
		env:          map[string]map[string]interface{}{},
		privateEnv:   map[string]map[string]interface{}{},
		requestFiles: []string{},
	}

	err := filepath.Walk("./testdata/", func(path string, info os.FileInfo, err error) error {
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

			if err = json.Unmarshal(raw, &ctx.env); err != nil {
				return errors.Wrap(err, "unable to parse cookies")
			}

		case ClientPrivateEnv:
			log.Println("found private client env...")
			raw, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "unable to read file")
			}

			if err = json.Unmarshal(raw, &ctx.privateEnv); err != nil {
				return errors.Wrap(err, "unable to unmarshal private env")
			}

		case ClientCookies:
			log.Println("found cookies..")
			ctx.cookies, err = parser.ParseCookies(path)
			if err != nil {
				return errors.Wrap(err, "unable to parse cookies")
			}
		default:
			if ext := filepath.Ext(path); ext == ".http" {
				log.Println("found http requests set", info.Name())
				ctx.requestFiles = append(ctx.requestFiles, path)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal("directory analise filed : ", err.Error())
	}

	ctx.mergeEnv()
	ctx.parseRequestFiles()
	fmt.Printf("%+v\n", ctx)
}

type ClientCtx struct {
	cookies      []parser.Cookie
	env          map[string]map[string]interface{}
	privateEnv   map[string]map[string]interface{}
	requestFiles []string
}

func (c *ClientCtx) mergeEnv() {
	for env := range c.privateEnv {
		if _, ok := c.env[env]; !ok {
			c.env[env] = map[string]interface{}{}
		}

		for key, val := range c.privateEnv[env] {
			c.env[env][key] = val
		}
	}
}

func (c *ClientCtx) parseRequestFiles() error {
	for _, file := range c.requestFiles {
		raw, _ := ioutil.ReadFile(file)

		var metrics bytes.Buffer
		metrics.WriteString(string(raw))
		scanner := bufio.NewScanner(&metrics)

		httpParser := httpParser{env: c.env["test"]}
		requests, _ := httpParser.ParseRequests(scanner)

		for i, request := range requests {
			fmt.Printf("%d %+v\n", i, request)
		}

	}

	return nil
}
