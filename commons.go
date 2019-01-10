package resource

import (
	"errors"
	"regexp"
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

// GetUrls will create urls to access pom and archive file
// returns: pomUrl, pomFileName, archiveUrl, archiveFileName
func GetUrls(src string, adef ArtifactDef, version string) (string, string, string, string) {
	r := regexp.MustCompile("\\/$") //to cut off last if slash if prsent
	srcURL := r.ReplaceAllString(src, "")
	gpath := strings.Replace(adef.GroupID, ".", "/", -1)
	srcBase := strings.Join([]string{srcURL, gpath, adef.ArtifactID, version}, "/")
	cordSlice := []string{adef.ArtifactID, version}
	if len(adef.Classifier) > 0 {
		cordSlice = append(cordSlice, adef.Classifier)
	}
	fileBaseCords := strings.Join(cordSlice, "-")
	pomFileName := fileBaseCords + ".pom"
	srcPom := strings.Join([]string{srcBase, pomFileName}, "/")
	//tracelog("srcPom: %v\n", srcPom)

	arFileName := fileBaseCords + "." + adef.AType

	srcArchive := strings.Join([]string{srcBase, arFileName}, "/")
	//tracelog("srcArchive: %v\n", srcArchive)

	return srcPom, pomFileName, srcArchive, arFileName
}

