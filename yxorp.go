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
	// flag.Parse() will parse the flags and do its magic
	//
	// As an added bonus we have a basic help message with the -h flag built-in. Try it out.
	flag.Parse()

	// load the config file
	// see the config.go file for the implementation of loadConfig()
	err := loadConfig()
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
}
