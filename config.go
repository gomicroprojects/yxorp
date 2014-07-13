package main

// This is the first time we have another file inside our main package.
// See http://golang.org/doc/code.html#Organization for information on how packages are organized.
//
// A directory can only belong to one package. Try to change the package in this file to "foo" or whatever:
// can't load package: package github.com/gomicroprojects/yxorp: found packages foo (config.go) and main (doc.go) in...

import (
	"encoding/json"
	// The fmt package (http://golang.org/pkg/fmt/)
	// Output/error formatting
	"fmt"
	// The os package (http://golang.org/pkg/os/)
	// We will use this to check, open and read the config file
	"os"
)

// ProxyConfig is a configuration for one proxy instance
// We need a target URL for httputil.NewSingleHostReverseProxy()
// GZ will indicate whether to gzip encode this proxy
//
// We will use the json package (http://golang.org/pkg/encoding/json/) to parse the config
type ProxyConfig struct {
	TargetURL string
	// The third part here is a field tag http://golang.org/ref/spec#Struct_types
	// The tag is used to instruct the json package to omit empty (zero value) fields
	//
	// By convention the tag is a concatenation of key:"value" pairs. I added an XML for
	// demonstration purposes
	// See http://golang.org/pkg/reflect/#StructTag
	GZ bool `json:",omitempty" xml:",attr"`
}

// Config is the representation of the config file
// It will be a JSON-Object with the name being the host.
//
// So it will look like this:
//     {
//         "www.example.com": {
//             "TargetURL": "http://localhost:8080/example"
//         },
//         "www2.example.com": {
//             "TargetURL": "http://localhost:8081/",
//             "GZ": true
//         }
//     }
type Config map[string]ProxyConfig

func loadConfig() (Config, error) {
	// Since we are in the same package as yxorp.go, we also have access to the variable "configFileName"
	// Let's get some information on the file
	cfgFileInfo, err := os.Stat(configFileName)
	// Always check for errors
	if err != nil {
		// fmt.Errorf() is like fmt.Sprintf, but returns an error type instead of a string
		// you will see this kind of "error chaining" quite often where the received error is
		// enrichted with local information.
		// In this case we add the config file name to the error message
		return nil, fmt.Errorf("could not open config file %s: %s", configFileName, err)
	}
	if cfgFileInfo.IsDir() {
		// the config file should be a file
		return nil, fmt.Errorf("config %s is a directory", configFileName)
	}
	// open the file for reading
	cfgFile, err := os.Open(configFileName)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %s", configFileName, err)
	}
	// A defer statement (http://golang.org/ref/spec#Defer_statements)
	// We want to make sure that the file is closed when we return from this function
	// Always close files, so they don't leak
	defer cfgFile.Close()

	// initialize a Config
	// The Config type has a map as an underlying type, so we can make() it
	cfg := make(Config)
	// This is one of the best features of the Go language: orthogonality
	// *os.File is an io.Reader as well, so we can use it directly with the json.Decoder
	dec := json.NewDecoder(cfgFile)
	err = dec.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("could not read config in %s: %s", configFileName, err)
	}
	return cfg, nil
}
