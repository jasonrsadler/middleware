package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	crlf       = "\r\n"
	colonspace = ": "
)

//ChecksumMiddleware function
func ChecksumMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)

		//get recorded headers
		var response []string
		for k, v := range rec.Header() {
			for _, val := range v {
				response = append(response, fmt.Sprintf("%s: %s\r\n", k, val))
			}
		}
		sort.Strings(response) //lexico

		//build 'X-Checksum-Headers'
		var ckH []string
		for k := range rec.Header() {
			ckH = append(ckH, k)
		}
		ckHName := "X-Checksum-Headers"
		sort.Strings(ckH)
		var xCksmH = strings.Join(ckH, ";")

		//build the canonical response
		var fResp []string
		fResp = append(fResp, strconv.Itoa(rec.Code)+crlf+strings.Join(response, "")+ckHName+colonspace+xCksmH+crlf+crlf+string(rec.Body.Bytes()))
		var canonicalResponse = strings.Join(fResp, "")

		//hash the canoncial response
		hash := sha1.New()
		hash.Write([]byte(canonicalResponse))
		xChecksum := hex.EncodeToString(hash.Sum(nil))

		//set headers
		w.Header().Set(ckHName, xCksmH)
		w.Header().Set("X-Checksum", xChecksum)
		for k, v := range rec.Header() {
			response = append(response, k)
			w.Header().Set(k, strings.Join(v, ""))
		}

		//write out headers
		w.WriteHeader(rec.Code)

		//send the body
		w.Write(rec.Body.Bytes())
	})
}

// Do not change this function.
func main() {
	var listenAddr = flag.String("http", "localhost:8080", "address to listen on for HTTP")
	flag.Parse()

	l := log.New(os.Stderr, "", 1)

	http.Handle("/", ChecksumMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Printf("%s - %s", r.Method, r.URL)
		w.Header().Set("X-Foo", "bar")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Date", "Sun, 08 May 2016 14:04:53 GMT")
		msg := "Curiosity is insubordination in its purest form.\n"
		w.Header().Set("Content-Length", strconv.Itoa(len(msg)))
		fmt.Fprintf(w, msg)
	})))

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
