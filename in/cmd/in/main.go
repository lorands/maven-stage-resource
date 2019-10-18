package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lorands/maven-stage-resource"
	"github.com/lorands/maven-stage-resource/in"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var trace bool
var resourceDir string

func main() {
	var request in.Request

	destinationDir := os.Args[1]

	var err error
	resourceDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		tracelog("Fail to get resource directory! Bailing out.")
		os.Exit(2)
	}

	tracelog("===> IN!")

	inputRequest(&request)

	trace = request.Source.Verbose

	var version string
	if len(request.Params.Version) > 0 {
		version = request.Params.Version
	} else {
		version = request.Version.Version
		if len(version) < 0 {
			fatal("Version is empty!", nil)
		}
	}

	version = readIfFile(version)

	if err := execute(request, version, destinationDir, resourceDir); err != nil {
		fatal("Fail to process.", err)
	}

	response := in.Response{
		Version: resource.Version{
			Version: version,
		},
	}

	outputResponse(response)

}

func readIfFile(version string) string {
	file, err := os.Open(version)
	if err != nil {
		return version
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return version
	}
	return line
}

func execute(request in.Request, version string, destinationDir string, resourceDir string) error {
	adef, err := resource.ArtifactStrToArtifactDef(request.Source.Artifact)
	if err != nil {
		fatal("Fail to process artfiact from resource source", err)
	}



	srcPom, srcPomFile, srcAr, srcArFile := resource.GetUrls(request.Source.Src, adef, version)
	//tPom, _, tAr, _ := resource.GetUrls(request.Source.Target, adef, version)
	//donwload to destFolder
	download(request, destinationDir, srcPom, srcPomFile)
	download(request, destinationDir, srcAr, srcArFile)
	//write version file
	ioutil.WriteFile(filepath.Join(destinationDir, "version"), []byte(version), 0644)

	if ! request.Params.DownloadOnly {
		//upload(request, destinationDir, srcPomFile, tPom)
		//upload(request, destinationDir, srcArFile, tAr)
		if err := mvnDeploy(resourceDir, request, destinationDir, srcArFile, srcPomFile); err != nil {
			return err
		}
	}

	return nil
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

func mvnDeploy(resourceDir string, request in.Request, srcDir string, archiveFileName string, pomFileName string) error {
	archivePath := filepath.Join(srcDir, archiveFileName)
	pomPath := filepath.Join(srcDir, pomFileName)

	params := [] string {
		"-s",
		"settings.xml",
		"deploy:deploy-file",
		fmt.Sprintf("-Durl=%s", request.Source.Target),
		fmt.Sprintf("-Drepository.username=%s", request.Source.Username),
		fmt.Sprintf("-Drepository.password=%s", request.Source.Password),
		fmt.Sprintf("-DpomFile=%s", pomPath),
		fmt.Sprintf("-Dfile=%s", archivePath),
	}

	return runCmd(resourceDir, "./mvnw", params)
}

func runCmd(workDir string, cmdStr string, strParams []string) error {

	tracelog("About to execute %s with params: %v", cmdStr, strParams)
	cmd := exec.Command(cmdStr, strParams...)
	cmd.Dir = workDir
	var sout bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				tracelog("Exit Status: %d", status.ExitStatus())
				if status.ExitStatus() != 0 {
					return fmt.Errorf("Non zero exit code form cli: %d\n%s", status.ExitStatus(), sout.String())
				}
			}
		} else {
			//log.Fatalf("cmd.Wait: %v", err)
			return err
		}
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
