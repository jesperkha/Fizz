package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	ErrNotALibrary = errors.New("'%s' is not a known library")
)

// Prints functions of the given library to the terminal. Returns
// an error if not a known library name.
func PrintDocs(libname string) error {
	filename := fmt.Sprintf("./lib/%s/%s_docs.md", libname, libname)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(ErrNotALibrary.Error(), libname)
	}

	section := strings.Split(string(file), "\n-->")[0]
	section = strings.TrimLeft(section, "<!-")
	// Todo: add color formatting to output
	fmt.Println(section)
	return nil
}
