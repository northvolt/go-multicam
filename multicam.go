package go-multicam

// OpenDriver starts up the Multicam drivers.
func OpenDriver() error {
	return ErrCannotOpenDriver
}

// CloseDriver closes the Multicam drivers. Call before exiting.
func CloseDriver() error {
	return nil
}
