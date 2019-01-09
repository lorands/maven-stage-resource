package resource

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

// ArtifactDef defines an artifact as a struct.
type ArtifactDef struct {
	GroupID    string
	ArtifactID string
	AType      string
	Classifier string
}

// ArtifactStrToArtifactDef converts artifact string to struct
func ArtifactStrToArtifactDef(artifact string) (ArtifactDef, error) {
	var def ArtifactDef

	splits := strings.Split(artifact, ":")

	if len(splits) < 3 {
		err := errors.New("you must specify at least GroupID:ArtifactID:type")
		return def, err
	}

	def.GroupID = splits[0]
	def.ArtifactID = splits[1]
	def.AType = splits[2]

	if len(splits) > 3 {
		def.Classifier = splits[3]
	}

	return def, nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
