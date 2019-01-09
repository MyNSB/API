package filesint

// Interface package for core files within application

import (
	"os"
	"mynsb-api/internal/util"
	"errors"
	"bytes"
	"io"
	"regexp"
)

func GetDirs() map[string]string {
	gopath := util.GetGOPATH()

	dirs := make(map[string]string)
	dirs["assets"] = gopath + "/src/mynsb-api/assets"
	dirs["sensitive"] = gopath + "/src/mynsb-api/sensitive"
	dirs["database"] = gopath + "/src/mynsb-api/sensitive"
	dirs["assets"] = gopath + "/src/mynsb-api/assets"

	return dirs
}

func LoadFile(srcDir string, pathQual string) (*os.File, error) {
	requestedDir := GetDirs()[srcDir]

	// Determine if the file exists
	if _, err := os.Stat(requestedDir + pathQual); err != nil {
		return nil, errors.New("could not locate file")
	}

	// Get a pointer to that file
	f, _ := os.Open(requestedDir + pathQual)

	return f, nil
}

func DataDump(srcDir string, pathQual string) ([]byte, error) {
	f, err := LoadFile(srcDir, pathQual)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	return buf.Bytes(), nil
}

func CreateFile(srcDir string, newFile string) (*os.File, error) {
	// Get the src dir location
	srcDirLoc := GetDirs()[srcDir]

	// Regex for extracting file information
	directoryMatch := regexp.MustCompile(`^(.*/)?(?:$|(.+))`)
	srcInfo := directoryMatch.FindStringSubmatch(newFile)

	// Because os.Create does not create folders for us we need to split the creation into 2 parts, directory creation and file creation
	if _, err := os.Stat(srcDirLoc + srcInfo[1]); os.IsNotExist(err) {
		os.Mkdir(srcDirLoc+srcInfo[1], 0777)
	}

	// Now actually create the file
	return os.Create(srcDirLoc + newFile)
}
