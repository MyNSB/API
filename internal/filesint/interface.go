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

func GetDirs() map[string]string {
	gopath := util.GetGOPATH()

	dirs := make(map[string]string)
	dirs["assets"]    = gopath + "/mynsb-api/assets"
	dirs["sensitive"] = gopath + "/mynsb-api/sensitive"
	dirs["database"]  = gopath + "/mynsb-api/database"
	dirs["assets"]    = gopath + "/mynsb-api/assets"

	return dirs
}

func LoadFile(srcDir string, pathQual string) (*os.File, error) {
	requestedDir := GetDirs()[srcDir]
	fullPath := filepath.FromSlash(requestedDir + pathQual)


	// Determine if the file exists
	if _, err := os.Stat(fullPath); err != nil {
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

func CreateFile(srcDir string, parentDir string, newFile string) (*os.File, error) {
	// Get the src dir location
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
