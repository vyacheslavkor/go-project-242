package fs

import "errors"

var (
	// ErrPathNotExist is returned when the requested path does not exist.
	ErrPathNotExist = errors.New("path not exists")

	// ErrPermissionDenied is returned when the process lacks access to the path.
	ErrPermissionDenied = errors.New("permission denied")
)
