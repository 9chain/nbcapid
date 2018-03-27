// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package primitives

/*
The error codes from and including -32768 to -32000 are reserved for pre-defined errors.
Any code within this range, but not defined explicitly below is reserved for future use.
The error codes are nearly the same as those suggested for XML-RPC at the following url:
http://xmlrpc-epi.sourceforge.net/specs/rfc.fault_codes.php

code				message						meaning
-32700				Parse error					Invalid JSON was received by the server.
												An error occurred on the server while parsing the JSON text.
-32600				Invalid Request				The JSON sent is not a valid Request object.
-32601				Method not found			The method does not exist / is not available.
-32602				Invalid params				Invalid method parameter(s).
-32603				Internal error				Internal JSON-RPC error.
-32000 to -32099	Server error				Reserved for implementation-defined server-errors.
*/

func NewParseError() *JSONError {
	return NewJSONError(-32700, "Parse error", nil)
}
func NewInvalidRequestError() *JSONError {
	return NewJSONError(-32600, "Invalid Request", nil)
}
func NewMethodNotFoundError() *JSONError {
	return NewJSONError(-32601, "Method not found", nil)
}
func NewInvalidParamsError() *JSONError {
	return NewJSONError(-32602, "Invalid params", nil)
}
func NewInternalError() *JSONError {
	return NewJSONError(-32603, "Internal error", nil)
}

/*******************************************************************/
func NewCustomInternalError(data interface{}) *JSONError {
	return NewJSONError(-32603, "Internal error", data)
}
func NewCustomInvalidParamsError(data interface{}) *JSONError {
	return NewJSONError(-32602, "Invalid params", data)
}
