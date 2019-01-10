package main

import (
	"encoding/json"
	"fmt"
	"log"

	// "io/ioutil"

	"github.com/lorands/maven-stage-resource"
	"github.com/lorands/maven-stage-resource/out"
	"os"
)

var trace bool

func main() {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Sprintf("usage: %v <sources directory>", os.Args[0]))
		os.Exit(1)
	}

	tracelog("===> OUT!")

	sourceDir := os.Args[1]

	var request out.Request
	inputRequest(&request)

	trace = request.Source.Verbose

	tracelog("Input directory set. %s\n", sourceDir)
	tracelog("Request params set: %v\n", request)

	//output to stdout...
	response := out.Response{
		Version: resource.Version{
			Version: "n/a",
		},
		Metadata: []resource.MetadataPair {

		},
	}

	outputResponse(response)
}


func fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}

func inputRequest(request *out.Request) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		log.Fatal("[OUT] reading request from stdin", err)
	}
}

func outputResponse(response out.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatal("[OUT] writing response to stdout", err)
	}
}

func tracelog(message string, args ...interface{}) {
	if trace {
		fmt.Fprintf(os.Stderr, message, args...)
	}
}
