// Package gz provides an http.Handler, which will gzip-Encode the responses
//
// This is the first time we write a library package. Therefore we should give the package a nice name.
// It lies in its own directory and can be imported with
// import github.com/gomicroprojects/yxorp/gz
//
// We didn't necessarily need a library here, but hey, why not.
package gz

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// We create our own http.Responsewriter, which we will pass to all following http.Handlers
// This chaining of handlers makes the http package so elegant. We can easily define middleware, which handles
// gzip, authentication, language detection etc.
//
// we don't need to export the struct (that's why it starts with a lowercase letter;
// see: http://golang.org/ref/spec#Exported_identifiers) since we will just wrap the handlers and replace
// the response writer (see below)
type gzResponseWriter struct {
	// this is the gzip writer, it also implements the io.Writer interface
	w *gzip.Writer
	// we embed an http.ResponseWriter (see http://golang.org/ref/spec#Struct_types )
	// and compose a new type gzResponseWriter, which is also a http.ResponseWriter
	http.ResponseWriter
}

// Write to override the http.ResponseWriter.Write() method
//
// We write to our gzip writer instead, which in turn will write to the http.ResponseWriter
func (gz gzResponseWriter) Write(p []byte) (int, error) {
	return gz.w.Write(p)
}

// GzHandler takes a handler and will return a new handler, which will gzip-encode the responses
func GzHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if the client accepts gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// if not, just serve the handler
			h.ServeHTTP(w, r)
			return
		}
		// set the required Content-Encoding header
		w.Header().Set("Content-Encoding", "gzip")

		// and now to the magic. this is one of the biggest strengths of the Go language: orthogonality
		// we can just put all the building blocks together, and they will fit. just like LEGO :D

		// *gzip.Writer is a writer, which will gzip encode everything before writing to the underlying writer
		// so we take the http.ResponseWriter, which is also a writer, put it into gzip
		gz := gzip.NewWriter(w)
		defer gz.Close()
		// create our own responsewriter
		gzr := gzResponseWriter{
			// with the gzip writer, which will write to the first http.ResponseWriter
			w: gz,
			// embed the first http.ResponseWriter for the other methods Header(), WriteHeader()
			ResponseWriter: w,
		}
		// and serve the handler
		h.ServeHTTP(gzr, r)
		// phew, pretty complicated stuff. but once you wrap your head around this, you will see the elegance
	})
}
