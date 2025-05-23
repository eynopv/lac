package errorsx

import "errors"

var ErrAuthUnknown = errors.New("unknown auth type")
var ErrAuthParse = errors.New("failed to parse auth")
var ErrApiAuthParse = errors.New("failed to parse api auth")
var ErrApiAuthInvalid = errors.New("header and key are required")
var ErrBasicAuthParse = errors.New("failed to parse basic auth")
var ErrBasicAuthInvalid = errors.New("username and password are required")
var ErrBearerAuthParse = errors.New("failed to parse bearer auth")
var ErrBearerAuthInvalid = errors.New("token is required")
var ErrUnsupportedVariablesValue = errors.New("unsupported variables type")
