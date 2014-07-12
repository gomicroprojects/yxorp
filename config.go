package main

// This is the first time we have another file inside our main package.
// See http://golang.org/doc/code.html#Organization for information on how packages are organized.
//
// A directory can only belong to one package. Try to change the package in this file to "foo" or whatever:
// can't load package: package github.com/gomicroprojects/yxorp: found packages foo (config.go) and main (doc.go) in...

import (
	// The fmt package (http://golang.org/pkg/fmt/)
	// Output/error formatting
	"fmt"
	// The os package (http://golang.org/pkg/os/)
	// We will use this to check, open and read the config file
	"os"
)

func loadConfig() error {
	// Since we are in the same package as yxorp.go, we also have access to the variable "configFileName"
	// Let's get some information on the file
	cfgFileInfo, err := os.Stat(configFileName)
	// Always check for errors
	if err != nil {
		// fmt.Errorf() is like fmt.Sprintf, but returns an error type instead of a string
		// you will see this kind of "error chaining" quite often where the received error is
		// enrichted with local information.
		// In this case we add the config file name to the error message
		return fmt.Errorf("could not open config file %s: %s", configFileName, err)
	}
	if cfgFileInfo.IsDir() {
		// the config file should be a file
		return fmt.Errorf("config %s is a directory.", configFileName)
	}
	return nil
}
