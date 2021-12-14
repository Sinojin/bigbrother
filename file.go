package bigbrother

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var FolderOrFileNameCannotBeEmpty = errors.New("Folder or file name can not be empty")
var FolderOrFileNotFound = errors.New("Folder or file not found!")

type FileInfo struct {
	ID   string //fixme:uuid maybe ??
	Name string

	Dir bool

	Files []*FileInfo

	CreateAt time.Time
	Version  int

	IsRoot bool
}

func (f *FileInfo) IsDir() bool {
	return f.Dir
}

func (f *FileInfo) ChangeName(name string) error {
	if name == "" {
		return FolderOrFileNameCannotBeEmpty
	}
	f.Name = name
	f.Version++
	return nil
}
func (f *FileInfo) getFullPathIfItsRoot(name string) (string, error) {
	if f.IsRoot == true {
		name, err := filepath.Abs(name)
		if err != nil {
			return "", err
		}
		return name, nil
	}
	return name, nil
}
func (f *FileInfo) Get(name string) (*FileInfo, error) {
	name, err := f.getFullPathIfItsRoot(name)
	if err != nil {
		return nil, err
	}
	//os.PathSeparator automatically chooses separator for os type
	//for instance for windows uses backslashes or linux uses slash
	osPathSeparator := string(os.PathSeparator)
	pList := strings.Split(name, osPathSeparator)
	lenthOfPList := len(pList)
	//if it's end of path, It should be found searched file or directory ,otherwise
	if lenthOfPList == 1 {
		if pList[0] == f.Name {
			return f, nil
		}
		return nil, FolderOrFileNotFound
	} else if lenthOfPList > 1 {
		fName := pList[0]
		if fName == f.Name {
			//it gets second name of path because we need to be sure that it exists or not
			SecondFName := pList[1]
			secondF, err := f.findNextFByName(SecondFName)
			if err != nil {
				return nil, err
			}
			//removing checked path
			pList = pList[1:]
			//Creating new path without first element
			newPath := strings.Join(pList, osPathSeparator)
			return secondF.Get(newPath)

		}
		return nil, FolderOrFileNotFound
	}
	return nil, FolderOrFileNotFound
}

//findNextFByName finds next file or directory if it's exists
func (f *FileInfo) findNextFByName(name string) (*FileInfo, error) {
	//it checks files exists or not
	if f.Files == nil {
		return nil, FolderOrFileNotFound
	}
	//it tries to find right file or directory
	for _, file := range f.Files {
		if file.Name == name {
			return file, nil
		}
	}
	return nil, FolderOrFileNotFound
}
