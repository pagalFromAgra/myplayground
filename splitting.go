package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// name := "1|2|3|"
	// splits := strings.Split(name, "|")
	// fmt.Println(strings.Split(name, "|"))
	// fmt.Printf("The length of 4th value is %d\n", len(splits[3]))

	// To set a key/value pair, use `os.Setenv`. To get a
	// value for a key, use `os.Getenv`. This will return
	// an empty string if the key isn't present in the
	// environment.

	if os.Getenv("proxy_test") == "" {
		os.Setenv("FOO", "1")
		fmt.Println("FOO:", os.Getenv("FOO"))
	}
	fmt.Println("proxy_test:", os.Getenv("proxy_test"))

	fmt.Println("TERM_PROGRAM:", os.Getenv("TERM_PROGRAM"))

	// Use `os.Environ` to list all key/value pairs in the
	// environment. This returns a slice of strings in the
	// form `KEY=value`. You can `strings.Split` them to
	// get the key and value. Here we print all the keys.
	fmt.Println()
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0])
	}
}
