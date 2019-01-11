package filesint

// Interface package for core files within application

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"mynsb-api/internal/util"
)


// UTILITY FUNCTIONS

// GetDirs returns all the directories associated with important information within the api, this includes:
// assets, sensitive data and database data
func GetDirs() map[string]string {

	gopath := util.GetGOPATH()
	return map[string]string{
		"assets": gopath + "/mynsb-api/assets",
		"sensitive": gopath + "/mynsb-api/sensitive",
		"database": gopath + "/mynsb-api/database",
	}
}


// LoadFile takes a src directory... e.g assets, sensitive, e.t.c and a path to a specific file within that directory and returns a file pointer to that file
func LoadFile(srcDir string, fileLoc string) (*os.File, error) {

	requestedDir := GetDirs()[srcDir]
	fullPath := filepath.FromSlash(requestedDir + fileLoc)

	// Determine if the file exists
	if _, notExists := os.Stat(fullPath); notExists != nil {
		return nil, errors.New("could not locate file")
	}

	f, _ := os.Open(requestedDir + fileLoc)
	return f, nil
}


// DataDump returns a data dump of all the data stored within a specific file
func DataDump(srcDir string, fileLoc string) ([]byte, error) {

	f, err := LoadFile(srcDir, fileLoc)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	dataBuf := bytes.NewBuffer(nil)
	io.Copy(dataBuf, f)

	return dataBuf.Bytes(), nil
}


// CreateFile takes a src directory, parent directory and a file name and creates a new file at srcDrc/parentDir/newFile.*
func CreateFile(srcDir string, parentDir string, newFile string) (*os.File, error) {

	srcDirLoc := GetDirs()[srcDir]

	// Because os.Create does not create folders for us we need to split the creation into 2 parts, directory creation and file creation
	if _, err := os.Stat(srcDirLoc + parentDir); os.IsNotExist(err) {
		err := os.MkdirAll(srcDirLoc+parentDir, 0777)
		if err != nil {
			return nil, err
		}
	}

	// Now actually create the file
	return os.Create(srcDirLoc + parentDir + "/" + newFile)
}
