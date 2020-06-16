package utils

// Unwrap panics, if there is an error occured.
// It must be used only in cases when error cannot be handled
// like in main functions
func Unwrap(err error) {
	if err != nil {
		panic(err)
	}
}
