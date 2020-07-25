package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const keyPrefix = "cache:"

var (
	errNotFound      = errors.New("not found")
	errAlreadyExists = errors.New("already exists")
)

// Cached is a cached item
type Cached struct {
	Status   int
	Body     []byte
	Header   http.Header
	ExpireAt time.Time
}

// Store interface for filesystems to implement
type Store interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Remove(string) error
}

// Options for cache
type Options struct {
	Store         Store
	Expire        time.Duration
	Headers       []string
	StripHeaders  []string
	DoNotUseAbort bool
}

// Cache struct implements Store interface
type Cache struct {
	Store
	options Options
	expires map[string]time.Time
}

// Get cache value
func (c *Cache) Get(key string) (*Cached, error) {
	data, err := c.Store.Get(key)
	if err != nil {
		return nil, err
	}
	var cch *Cached
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err = dec.Decode(&cch)
	if err != nil {
		return nil, err
	}
	if cch.ExpireAt.Nanosecond() != 0 && cch.ExpireAt.Before(time.Now()) {
		err := c.Store.Remove(key)
		return nil, err
	}
	return cch, nil
}

// Set cache value
func (c *Cache) Set(key string, cch *Cached) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(*cch)
	if err != nil {
		return err
	}
	return c.Store.Set(key, b.Bytes())
}

type wrappedWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

// Write response
func (rw *wrappedWriter) Write(body []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(body)
	if err == nil {
		rw.body.Write(body)
	}
	return n, err
}

// New cache
func New(o ...Options) gin.HandlerFunc {
	opts := Options{
		Store:  NewInMemory(),
		Expire: 0,
	}

	for _, i := range o {
		opts = i
		break
	}

	cache := Cache{
		Store:   opts.Store,
		options: opts,
		expires: make(map[string]time.Time),
	}

	return func(c *gin.Context) {

		// only GET method available for caching
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		toHash := c.Request.URL.RequestURI()
		for _, k := range cache.options.Headers {
			if v, ok := c.Request.Header[k]; ok {
				toHash += k
				toHash += strings.Join(v, "")
			}
		}

		key := keyPrefix + md5String(toHash)

		if cch, _ := cache.Get(key); cch == nil {
			// cache miss
			writer := c.Writer
			rw := wrappedWriter{ResponseWriter: c.Writer}
			c.Writer = &rw
			c.Next()
			c.Writer = writer

			header := rw.Header()

			for _, k := range cache.options.StripHeaders {
				header.Del(k)
			}

			cache.Set(key, &Cached{
				Status: rw.Status(),
				Body:   rw.body.Bytes(),
				Header: rw.Header(),
				ExpireAt: func() time.Time {
					if cache.options.Expire == 0 {
						return time.Time{}
					}
					return time.Now().Add(cache.options.Expire)
				}(),
			})

		} else {
			// cache found
			start := time.Now()
			c.Writer.WriteHeader(cch.Status)
			for k, val := range cch.Header {
				for _, v := range val {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.Header().Add("X-Gin-Cache", fmt.Sprintf("%f ms", time.Now().Sub(start).Seconds()*1000))
			c.Writer.Write(cch.Body)

			if !cache.options.DoNotUseAbort {
				c.Abort()
			}
		}
	}
}

func md5String(url string) string {
	h := md5.New()
	io.WriteString(h, url)
	return hex.EncodeToString(h.Sum(nil))
}

func init() {
	gob.Register(Cached{})
}
