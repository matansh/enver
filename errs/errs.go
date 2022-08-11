// Package errs is a single repository for all the errors this lib may return
package errs

import "errors"

var (
	ErrNotPtr    = errors.New("the provided configuration needs to be a pointer")
	ErrNotStruct = errors.New("the provided configuration needs to be a struct")
	// ErrUnsupportedType means that you need to add the missing type to the switch in ScanStructFromEnv.
	ErrUnsupportedType = errors.New("unsupported data type")
	ErrEnvVarNotFound  = errors.New("environment did not contain wanted var")
)
