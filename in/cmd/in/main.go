package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/lorands/maven-stage-resource"
	"github.com/lorands/maven-stage-resource/in"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var trace bool

func main() {
	var request in.Request

	destinationDir := os.Args[1]

	tracelog("===> IN!")

	inputRequest(&request)

	trace = request.Source.Verbose

	version := request.Version.Version
	if len(version) < 0 {
		fatal("Version is empty!", nil)
	}

	execute(request, version, destinationDir)

	//http put to target repo

	response := in.Response{
		Version: resource.Version{
			Version: version,
		},
	}

	outputResponse(response)

}

func execute(request in.Request, version string, destinationDir string) {
	adef, err := resource.ArtifactStrToArtifactDef(request.Source.Artifact)
	if err != nil {
		fatal("Fail to process artfiact from resource source", err)
	}

	if len(request.Params.Version) > 0 {
		version = request.Params.Version
	}

	srcPom, srcPomFile, srcAr, srcArFile := resource.GetUrls(request.Source.Src, adef, version)
	tPom, _, tAr, _ := resource.GetUrls(request.Source.Target, adef, version)
	//donwload to destFolder
	download(request, destinationDir, srcPom, srcPomFile)
	download(request, destinationDir, srcAr, srcArFile)
	//write version file
	ioutil.WriteFile(filepath.Join(destinationDir, "version"), []byte(version), 0644)

	if ! request.Params.DownloadOnly {
		upload(request, destinationDir, srcPomFile, tPom)
		upload(request, destinationDir, srcArFile, tAr)
	}
}

func download(request in.Request, destDir string, src string, fileName string) {
	var client http.Client

	tracelog("To download from url: %s\n", src)

	req, err := http.NewRequest("GET", src, nil)
	if err != nil {
		fatal("Fail to create request object for download", err)
	}
	if request.Source.Username != "" {
		tracelog("Setting basic authorization as requested for user: %v\n", request.Source.Username)
		req.SetBasicAuth(request.Source.Username, request.Source.Password)
	}
	resp, err := client.Do(req)

	if err != nil {
		fatal(fmt.Sprintf("Error response from http. %v\n", resp), err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		tracelog("Status code for download: %d\n", resp.StatusCode)
	} else {
		fatal(fmt.Sprintf("Fail to download artifact. Status code %s", resp.Status), nil)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath.Join(destDir, fileName))
	defer out.Close()

	n, err := io.Copy(out, resp.Body)

	if err != nil || n < 1 {
		fatal(fmt.Sprintf("Faild to download file. %v\n", src), err)
	}

}

func upload(request in.Request, srcDir string, fileName string, dest string) error {

	path := filepath.Join(srcDir, fileName)
	tracelog("Upload file: %s to url: %s", path, dest )
	client := &http.Client{}
	var reader io.Reader

	file, err := os.Open(path)
	defer file.Close()
	reader = bufio.NewReader(file)
	req, err := http.NewRequest("PUT", dest, reader)
	if err != nil {
		return err
	}
	if request.Source.Username != "" {
		tracelog("Setting basic authorization as requested for user: %v\n", request.Source.Username)
		req.SetBasicAuth(request.Source.Username, request.Source.Password)
	}
	resp, err := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		tracelog("Status code for download: %d\n", resp.StatusCode)
	} else {
		fatal(fmt.Sprintf("Fail to upload artifact. Status code %s\n", resp.Status), nil)
	}
	defer resp.Body.Close()

	if err != nil {
		fatal(fmt.Sprintf("Error response from http. %v\n", resp), err)
	}
	return nil
}

func fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}

func inputRequest(request *in.Request) {
	//reader := bufio.NewReader(os.Stdin)
	//text, _ := reader.ReadString('\n')
	//tracelog("IN: stdin: %s\n", text)
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
	//if err := json.Unmarshal([]byte(text), request); err != nil {
		log.Fatal("[IN] reading request from stdin: ", err)
	}
}

func outputResponse(response in.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatal("[IN] writing response to stdout: ", err)
	}
}

func tracelog(message string, args ...interface{}) {
	if trace {
		fmt.Fprintf(os.Stderr, message, args...)
	}
}
