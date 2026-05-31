package app

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
)

type Env struct {
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	HTTPClient *http.Client
	Now        func() time.Time
	GOOS       string
	GOARCH     string
}

func defaultEnv() Env {
	return Env{
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		HTTPClient: http.DefaultClient,
		Now:        time.Now,
		GOOS:       runtime.GOOS,
		GOARCH:     runtime.GOARCH,
	}
}

func (e Env) httpClient() *http.Client {
	if e.HTTPClient != nil {
		return e.HTTPClient
	}
	return http.DefaultClient
}

func (e Env) now() time.Time {
	if e.Now != nil {
		return e.Now()
	}
	return time.Now()
}
