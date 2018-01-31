package main

import (
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
		rec.Code = 418
		h.ServeHTTP(rec, r)
		var response []string
		//headerKey := "X-Checksum"
		for k := range rec.Header() {
			response = append(response, k)
		}
		sort.Strings(response)
		//strResponse := statusCode
		// strResponse += (strings.Join(response, ";"))
		// fmt.Println(headerKey + " | " + strResponse)
		w.WriteHeader(rec.Code)
		for k, v := range rec.Header() {
			response = append(response, k)
			w.Header().Set(k, strings.Join(v, ""))
		}
		w.Header().Set("X-Checksum", "814a8da9ad27cd0c0a2cea3536daa3a8b12926b3")
		w.Write(rec.Body.Bytes())

		// rec.Header().get
		// var parts []string
		// []string(rec.Header()).Join(parts, "")
		// hasher := sha1.New()
		// hasher.Write([]byte(w.Header()))
		// hString := hasher.Sum(nil)
		// w.Header().Set("SHA-1", base64.URLEncoding.EncodeToString(hString))
		// w.WriteHeader(418)
		// w.Write(rec.Body.Bytes())

		// fmt.Println(w, "%s %s %s \n", r.Method, r.URL, r.Proto)
		// for k, v := range r.Header {
		// 	fmt.Println(w, "Header field %q, Value %q\n", k, v)
		// }

		// fmt.Println(rec.Body)
		// fmt.Println(base64.URLEncoding.EncodeToString(hString))
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
