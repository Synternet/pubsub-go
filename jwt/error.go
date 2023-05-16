package jwt

import "errors"

var ErrNoImportAndExport = errors.New("permission violation: no imports or exports claimed")
