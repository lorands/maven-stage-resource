package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	resource "github.com/lorands/maven-stage-resource"
	"github.com/lorands/maven-stage-resource/in"
)

var trace bool

func main() {
	var request *in.Request

	destinationDir := os.Args[1]

	inputRequest(request)

	trace = request.Source.Verbose

	version := request.Version.Version
	if len(version) < 0 {
		fatal("Version is empty!", nil)
	}

	//do it babe
	adef, err := resource.ArtifactStrToArtifactDef(request.Source.Artifact)
	if err != nil {
		fatal("Fail to process artfiact from resource source", err)
	}

	srcPom, srcPomFile, srcAr, srcArFile := getUrls(request.Source.Src, adef, version)
	tPom, _, tAr, _ := getUrls(request.Source.Target, adef, version)

	//donwload to destFolder

	download(destinationDir, srcPom, srcPomFile)
	download(destinationDir, srcAr, srcArFile)

	upload(destinationDir, srcPomFile, tPom)
	upload(destinationDir, srcArFile, tAr)

	//http put to target repo

	response := in.Response{
		Version: resource.Version{
			Version: version,
		},
	}

	outputResponse(response)

}

func getUrls(src string, adef resource.ArtifactDef, version string) (string, string, string, string) {
	r := regexp.MustCompile("\\/$")
	srcURL := r.ReplaceAllString(src, "")
	srcBase := strings.Join([]string{srcURL, adef.GroupID, adef.ArtifactID, version}, "/")
	cordSlice := []string{adef.ArtifactID, version}
	if len(adef.Classifier) > 0 {
		cordSlice = append(cordSlice, adef.Classifier)
	}
	fileBaseCords := strings.Join(cordSlice, "-")
	pomFileName := fileBaseCords + ".pom"
	srcPom := strings.Join([]string{srcBase, pomFileName}, "/")
	tracelog("srcPom: %v\n", srcPom)

	arFileName := fileBaseCords + "." + adef.AType

	srcArchive := strings.Join([]string{srcBase, arFileName}, "/")
	tracelog("srcArchive: %v\n", srcArchive)

	return srcPom, pomFileName, srcArchive, arFileName
}

func fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}

func inputRequest(request *in.Request) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		log.Fatal("reading request from stdin", err)
	}
}

func outputResponse(response in.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatal("writing response to stdout", err)
	}
}

func tracelog(message string, args ...interface{}) {
	if trace {
		fmt.Fprintf(os.Stderr, message, args...)
	}
}
