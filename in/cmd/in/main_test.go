package main

import (
	"github.com/lorands/maven-stage-resource"
	"github.com/lorands/maven-stage-resource/in"
	"io/ioutil"
	"testing"
)


func TestExecute(t *testing.T) {
	trace = true

	request := &in.Request {
		Source: resource.Source{
			Artifact: "commons-lang:commons-lang:jar",
			Src: "https://repo1.maven.org/maven2",
			Target: "https://boggusogus.com/artifact/",
			Verbose: true,
			//Username: "",
			//Password: "",
		},
		Version: resource.Version{
			Version: "2.4",
		},
		Params: in.Params {
			Version: "",
			DownloadOnly: true,
		},
	}

	version := "2.5"

	destDir, _ := ioutil.TempDir( "","maven-stage")

	execute(request, version, destDir)

	//check if files are there...
	files, _ := ioutil.ReadDir(destDir)

	if len(files) != 3 {
		t.Errorf("Excpected to have 3 files, insted we have %d files", len(files))
	}

	cntr := 0
	for _, file := range files {
		if file.Name() == "commons-lang-2.5.jar" {
			cntr++
		} else if file.Name() == "commons-lang-2.5.pom" {
				cntr++
		} else if file.Name() == "version" {
					cntr++
		} else {
						t.Errorf("Unknow file found: %s", file.Name())
		}
	}
	if cntr != 3 {
		t.Errorf("Not all required files found, expecting 3 well defined, found: %d", cntr)
	}

}

//func TestExecute2(t *testing.T) {
//	trace = true
//	request := &in.Request {
//		Source: resource.Source{
//			Artifact: "my.project.test:test-server:jar",
//			Src: "https://foo.bar/repository/releases",
//			Target: "https://foo.bar/repository/mohosz-uat",
//			Verbose: true,
//			Username: "ci",
//			Password: "XxxXxxXxxX",
//		},
//		Version: resource.Version{
//			Version: "1.0.191",
//		},
//		Params: in.Params {
//			//Version: "",
//			DownloadOnly: false,
//		},
//	}
//
//	version := "1.0.393"
//
//	destDir, _ := ioutil.TempDir( "","maven-stage")
//	//destDir := "/tmp/oo"
//
//	execute(request, version, destDir)
//
//	//check if files are there...
//	files, _ := ioutil.ReadDir(destDir)
//
//	if len(files) != 3 {
//		t.Errorf("Excpected to have 3 files, insted we have %d files", len(files))
//	}
//
//}