# nikel-cache

[![GoDoc](https://pkg.go.dev/badge/github.com/nikel-api/nikel-cache?status.svg)](https://pkg.go.dev/github.com/nikel-api/nikel-cache?tab=doc) [![Go Report Card](https://goreportcard.com/badge/github.com/nikel-api/nikel-cache)](https://goreportcard.com/report/github.com/nikel-api/nikel-cache)

A simple and performant cache middleware for Gin based on [olebedev's gin-cache](https://github.com/olebedev/gin-cache).

Go version +1.13 required.

### Available Database Backends

* In-Memory
* LevelDB (currently used by Nikel API)
* BadgerDB

### Usage

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nikel-api/nikel-cache"
)

func main() {
	r := gin.New()

	r.Use(cache.New(cache.Options{
		// set zero to make cache never expire
		Expire: 0,

		// set store
		Store: func() *cache.LevelDB {
			store, err := cache.NewLevelDB("cache")
			if err != nil {
				panic(err)
			}
			return store
		}(),

		// uses the header fields to calculate key
		Headers: []string{},

		// strips header fields
		StripHeaders: []string{},

		// bypass cache by response code
		BypassCodes: map[int]bool{},

		// *gin.Context.Abort() will be invoked immediately after cache has been served
		DoNotUseAbort: false,
	}))

	r.Run()
}

```
