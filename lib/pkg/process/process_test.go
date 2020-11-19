package process

import (
	"errors"
	"os"
	"testing"
	"time"
)

// 0: Not started
// 1: Not directory, filename not a number
// 2: Is directory, filename is a number
// 3: Is directory, filename not a number
var stateTestCheckProcessLinux int

type testFileInfo struct {
}

func (fi testFileInfo) Name() string {
	var name string

	switch stateTestCheckProcessLinux {
	case 1:
		name = "testFilename"
	case 2:
		name = "123"
	case 3:
		name = "notanumber"
	}

	return name
}

func (fi testFileInfo) Size() int64 {
	return 0
}

func (fi testFileInfo) Mode() os.FileMode {
	return os.ModePerm
}

func (fi testFileInfo) ModTime() time.Time {
	return time.Now()
}

func (fi testFileInfo) IsDir() bool {
	var isDir bool

	stateTestCheckProcessLinux = stateTestCheckProcessLinux + 1

	switch stateTestCheckProcessLinux {
	case 2, 3:
		isDir = true
	default:
		isDir = false
	}

	return isDir
}

func (fi testFileInfo) Sys() interface{} {
	return nil
}

func TestCheckProcessLinux(t *testing.T) {
	goodOutput := func(string) ([]byte, error) {
		return []byte("123 (bash) 1 1 1"), nil
	}

	errorReturn := func(string) ([]byte, error) {
		return []byte("123 (bash) 1 1 1"), errors.New("read data error")
	}

	badOutput := func(string) ([]byte, error) {
		return []byte("123 bash 1 1 1"), nil
	}

	procName, err := getPidNameWithHandler(goodOutput, 123)

	if err != nil {
		t.Error("getPidNameWithHandler returned an error on valid data")
	}

	if procName == "" {
		t.Error("getPidNameWithHandler did not return a valid name on valid data")
	}

	procName, err = getPidNameWithHandler(errorReturn, 123)

	if err == nil {
		t.Error("getPidNameWithHandler did not return an error on a read data error")
	}

	if procName != "" {
		t.Error("getPidNameWithHandler returned a name on a read data error")
	}

	procName, err = getPidNameWithHandler(badOutput, 123)

	if err == nil {
		t.Error("getPidNameWithHandler did not return an error when data parse should fail")
	}

	if procName != "" {
		t.Error("getPidNameWithHandler returned a name when data parse should fail")
	}

	svc := processByNameHandlers{
		open: func(n string) (*os.File, error) {

			return nil, nil
		},
		close: func(f *os.File) error {
			return nil
		},
		readDir: func(f *os.File, entries int) ([]os.FileInfo, error) {
			fi := testFileInfo{}
			fiSlice := []os.FileInfo{fi, fi, fi}

			return fiSlice, nil
		},
		getPidName: getPidNameWithHandler,
		readFile: func(string) ([]byte, error) {
			return []byte("123 (bash) 1 1 1"), nil
		},
	}

	fileList, err := getProcessesByNameWithHandlers(svc, "bash")

	if err != nil {
		t.Error("getProcessesByNameWithHandlers returned an error when given valid input")
	}
	if fileList == nil && len(fileList) != 1 {
		t.Error("getProcessesByNameWithHandlers file list not correct when given valid input")
	}

	errString := "read directory error"
	svc.readDir = func(f *os.File, entries int) ([]os.FileInfo, error) {
		return nil, errors.New(errString)
	}

	fileList, err = getProcessesByNameWithHandlers(svc, "bash")

	if err == nil || err.Error() != errString {
		t.Error("getProcessesByNameWithHandlers should have returned a read directory error")
	}

	if fileList != nil {
		t.Error("getProcessesByNameWithHandlers returned a file list but should have returned an error")
	}

	errString = "file open error"
	svc.open = func(n string) (*os.File, error) {
		return nil, errors.New(errString)
	}

	fileList, err = getProcessesByNameWithHandlers(svc, "bash")

	if err == nil || err.Error() != errString {
		t.Error("getProcessesByNameWithHandlers should have returned a file open error")
	}

	if fileList != nil {
		t.Error("getProcessesByNameWithHandlers returned a file list but should have returned an error")
	}
}
