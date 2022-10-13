package multicam

import "os"

// GoMulticamVersion of this package, for display purposes.
// Change this variable on a new package release.
const GoMulticamVersion = "0.3.0"

// Version returns the current Golang package version.
func Version() string {
	return GoMulticamVersion
}

// SDKVersion returns the current SDK version using an ENV variable, since not supported by Euresys directly.
// It can be set automatically within a Docker container built using ARG/ENV variables, as is done by the
// Dockerfile located in this package.
// If not set by the host or by container, it will return empty string.
func SDKVersion() string {
	return os.Getenv("MULTICAM_SDK_VERSION")
}
