package utils

import (
	"fmt"
	"os"
)

// Unwrap exits, if there is an error occured.
// It must be used only in cases when error cannot be handled
// like in main functions
func Unwrap(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// UnwrapWith it is Unwrap but with note
func UnwrapWith(err error, note string) {
	if err != nil {
		err := fmt.Errorf("%s: %s", note, err.Error())
		fmt.Println(err)
		os.Exit(-1)
	}
}
