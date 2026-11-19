package rpc

import (
	"errors"
	"fmt"

	"yukumo-script/core/api"
)

const (
	codeParseError     = -32700
	codeInvalidRequest = -32600
	codeMethodNotFound = -32601
	codeInvalidParams  = -32602
	codeInternalError  = -32603
	codeServerError    = -32000
)

const (
	KindInvalid  = "invalid"
	KindNotFound = "not_found"
	KindExists   = "exists"
	KindInternal = "internal"
)

type paramsError struct {
	err error
}

func (e paramsError) Error() string { return e.err.Error() }
func (e paramsError) Unwrap() error { return e.err }

func invalidParams(err error) error {
	return paramsError{err: err}
}

func errData(kind, detail string) map[string]string {
	out := map[string]string{"kind": kind}
	if detail != "" {
		out["detail"] = detail
	}
	return out
}

func mapError(err error) (code int, message string, data any) {
	var pe paramsError
	if errors.As(err, &pe) {
		return codeInvalidParams, "Invalid params", errData(KindInvalid, pe.Error())
	}
	switch {
	case errors.Is(err, api.ErrNotFound):
		return codeServerError, err.Error(), errData(KindNotFound, err.Error())
	case errors.Is(err, api.ErrExists):
		return codeServerError, err.Error(), errData(KindExists, err.Error())
	case errors.Is(err, api.ErrInvalid):
		return codeServerError, err.Error(), errData(KindInvalid, err.Error())
	default:
		return codeInternalError, "Internal error", errData(KindInternal, fmt.Sprint(err))
	}
}
