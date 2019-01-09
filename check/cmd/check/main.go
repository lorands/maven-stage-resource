package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/lorands/maven-stage-resource"
	"github.com/lorands/maven-stage-resource/check"
)

var trace bool

func main() {
	var request check.Request
	inputRequest(&request)

	trace = request.Source.Verbose

	adef, err := resource.ArtifactStrToArtifactDef(request.Source.Artifact)

	if err != nil {
		fatal("Fail to process artfiact from resource source", err)
	}

	r := regexp.MustCompile("\\/$")
	srcURL := r.ReplaceAllString(request.Source.Src, "")
	metaURL := strings.Join([]string{srcURL, adef.GroupID, adef.ArtifactID, "maven-metadata.xml"}, "/")
	tracelog("MetaURL: %v\n", metaURL)

	var client http.Client

	req, err := http.NewRequest("GET", metaURL, nil)
	if err != nil {
		fatal("Fail to create request object to maven-metadata.xml", err)
	}
	if request.Source.Username != "" {
		tracelog("Setting basic authorization as requested for user: %v\n", request.Source.Username)
		req.SetBasicAuth(request.Source.Username, request.Source.Password)
	}
	resp, err := client.Do(req)
	if err != nil {
		fatal(fmt.Sprintf("Error response from http. %v\n", resp), err)
	}
	defer resp.Body.Close()

	var response check.Response

	data, err := ioutil.ReadAll(resp.Body)

	tracelog("Meta: %s\n", data)

	if err != nil {
		fatal("Fail to read maven-metadata.xml", err)
	} else {
		var result check.MavenMetadata
		err = xml.Unmarshal(data, &result)
		if err != nil {
			fatal("Fail to process xml", err)
		}

		// tracelog("Unmarshalled meta: %v\n", result)
		// var greaterVersions []string
		tracelog("Latest seen version: %s\n", request.Version.Version)
		tracelog("Release version in meta: %s\n", result.Versioning.Release)

		if result.Versioning.Release != request.Version.Version {
			// for _, version := range result.Versioning.Versions {
			// 	if version.Version < request.Version.Version {
			// 		greaterVersions = append(greaterVersions, version.Version)
			// 	}
			// }
			//releaseVersion := check.Response.Version{Version: "1.1"}
			item := resource.Version{Version: result.Versioning.Release}
			response = append(response, item)
		}
	}

	outputResponse(response)
}

func fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}

func inputRequest(request *check.Request) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		log.Fatal("reading request from stdin", err)
	}
}

func outputResponse(response check.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatal("writing response to stdout", err)
	}
}

func tracelog(message string, args ...interface{}) {
	if trace {
		fmt.Fprintf(os.Stderr, message, args...)
	}
}
