package lib

import (
	"embed"
	"errors"
	"fmt"
)

var (
	ErrNotALibrary = errors.New("'%s' is not a known library")

	//go:embed _libdump
	embeddedDocs embed.FS
)

// Prints functions of the given library to the terminal. Returns
// an error if not a known library name.
func PrintDocs(libname string) error {
	filename := fmt.Sprintf("_libdump/%s.txt", libname)
	file, err := embeddedDocs.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(ErrNotALibrary.Error(), libname)
	}

	fmt.Println()
	fmt.Print(string(file))
	return nil
}
