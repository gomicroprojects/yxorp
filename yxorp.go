package main

// package main
//
// The package clause (http://golang.org/ref/spec#Package_clause)
// "main" indicates that this file is a command and will generate an executable
// binary.
//
// Take a look at http://golang.org/doc/code.html for an introduction on
// "How to Write Go Code"

// The import declaration (http://golang.org/ref/spec#Import_declarations)
// This is in the multiline format.
// You could also write:
//     import a
//     import b
import (
	// See gz/gz.go for the implementation of our first own library
	"github.com/gomicroprojects/yxorp/gz"
	// The http package (http://golang.org/pkg/net/http/) for our HTTP front servers
	"net/http"
	// The httputil package (http://golang.org/pkg/net/http/httputil/) for the reverse proxy implementation
	"net/http/httputil"
	// The url package (http://golang.org/pkg/net/url/) to parse the URL string in the config
	"net/url"
	// The time package for the timeout constants for our HTTP server
	"time"
	// The os package (http://golang.org/pkg/os/) for os.Exit()
	"os"
	// The fmt package (http://golang.org/pkg/fmt/) for formatted I/O
	"fmt"
	// The flag package (http://golang.org/pkg/flag/)
	// We will us this to parse command line flags
	"flag"
)

// We declare a variable "configFileName", which is a string
// Since it is just declared, it is inialized to its zero value.
// The zero value for a string is an empty string ""
var configFileName string

// The server address to listen to
// This will be set by a flag as well
var serverAddress string

// We declare the map "proxyMap"
// It will map a host to an http.Handler
// In the default case (no gzip) it will be a *httputil.ReverseProxy, which implements the http.Handler
// see (*httputil.ReverseProxy).ServeHTTP(http.ResponseWriter, *http.Request)
var proxyMap map[string]http.Handler

func main() {
	// Setting up the flags using the flag package
	//
	// The StringVar function takes a pointer to a string. Since the flag package will modify the contents of the
	// "configFileName" string for us, it needs to know the address of the string. Simply passing the value of the
	// (now empty) string would not be enough.
	// To get the pointer to our "configFileName" var, we take the address of it with the ampersand operator "&".
	//
	// Our flag name will be "c", the default value an empty string, and a short description.
	flag.StringVar(&configFileName, "c", "", "The config file name to use. Example: /tmp/yxorp.json")
	// and the server address
	flag.StringVar(&serverAddress, "a", ":8080", "The server address to listen to.")
	// flag.Parse() will parse the flags and do its magic
	//
	// As an added bonus we have a basic help message with the -h flag built-in. Try it out.
	flag.Parse()

	// load the config file
	// see the config.go file for the implementation of loadConfig()
	cfg, err := loadConfig()
	if err != nil {
		// print out the error
		fmt.Println(err)
		// os.Args[0] is the command name
		fmt.Printf("\nUsage of %s\n", os.Args[0])
		// this will print out the help for the flags
		flag.PrintDefaults()
		// exit the program with an error code
		os.Exit(1)
	}
	// initalize the proxy map
	proxyMap = make(map[string]http.Handler)
	// we range over the cfg entries (see http://golang.org/ref/spec#RangeClause)
	// the key will be the host name (since it is a host-based reverse proxy)
	for host, proxyCfg := range cfg {
		// parse the url
		targetURL, err := url.Parse(proxyCfg.TargetURL)
		if err != nil {
			// exit on any error
			fmt.Printf("error on config host %s parsing target URL: %s", host, err)
			os.Exit(1)
		}
		// NewSingleHostReverseProxy will return a *httputil.Reversproxy, which in turn is a http.Handler
		// that's why we can assign it to the map
		proxyMap[host] = httputil.NewSingleHostReverseProxy(targetURL)
		// config requested to gzip encode this proxy
		if proxyCfg.GZ {
			// so, we will wrap our default proxy handler with our gzip handler
			proxyMap[host] = gz.GzHandler(proxyMap[host])
		}
	}

	// Create an HTTP server
	// You can read this as: server is a pointer to(take the address of) an
	// http.Server with fields...
	//
	// This is equivalent to:
	// server := new(http.Server)
	// server.Addr = serverAddress
	// server.Handler = proxy()
	// etc.
	server := &http.Server{
		Addr: serverAddress,
		// the implementation is further below in this file, the function proxy() will return an http.Handler
		Handler:      proxy(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// Start serving
	err = server.ListenAndServe()
	if err != nil {
		// if there was an error, it is usually a fatal error and we can't continue serving
		fmt.Println(err)
		os.Exit(1)
	}
}

// will return a basic proxy handling http.Handler
func proxy() http.Handler {
	// http.HandlerFunc will turn our anonymous func(w http.ResponseWriter, r *http.Request) into an http.Handler
	// check out the code for http.HandlerFunc here: http://golang.org/pkg/net/http/#HandlerFunc
	// this is an example of a very elegant and surprising use of the language features
	// http.HandlerFunc is a type with an underlying type func(http.ResponseWriter, *http.Request)
	// the HandlerFunc implements the http.Handler, and the ServeHTTP() implementation will call the underlying func type
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we will check if there is an entry for he request Host
		h, ok := proxyMap[r.Host]
		if !ok {
			// no entry, HTTP status not found (is this the correct status?)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// proceed with the matched handler
		h.ServeHTTP(w, r)
	})
}
