package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

type Method struct {
	Name   string
	Req    reflect.Type
	Res    reflect.Type
	Handle func(ctx context.Context, params json.RawMessage) (any, error)
}

func Def[Req, Res any](name string, f func(context.Context, Req) (Res, error)) Method {
	return Method{
		Name: name,
		Req:  reflect.TypeFor[Req](),
		Res:  reflect.TypeFor[Res](),
		Handle: func(ctx context.Context, raw json.RawMessage) (any, error) {
			var req Req
			if len(raw) > 0 && string(raw) != "null" {
				if err := json.Unmarshal(raw, &req); err != nil {
					return nil, invalidParams(err)
				}
			}
			return f(ctx, req)
		},
	}
}

type Registry struct {
	byName  map[string]Method
	methods []Method
}

func NewRegistry() *Registry {
	return &Registry{byName: make(map[string]Method)}
}

func (r *Registry) Register(methods ...Method) {
	if r.byName == nil {
		r.byName = make(map[string]Method)
	}
	for _, m := range methods {
		if m.Name == "" {
			panic("rpc: empty method name")
		}
		if _, ok := r.byName[m.Name]; ok {
			panic(fmt.Sprintf("rpc: duplicate method %q", m.Name))
		}
		r.byName[m.Name] = m
		r.methods = append(r.methods, m)
	}
}

func (r *Registry) Methods() []Method {
	out := make([]Method, len(r.methods))
	copy(out, r.methods)
	return out
}

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      json.RawMessage `json:"id"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (r *Registry) Handle(ctx context.Context, raw []byte) []byte {
	var req rpcRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return mustMarshal(rpcResponse{
			JSONRPC: "2.0",
			Error:   &rpcError{Code: codeParseError, Message: "Parse error", Data: errData("internal", err.Error())},
			ID:      nil,
		})
	}
	if req.JSONRPC != "" && req.JSONRPC != "2.0" {
		return errorResponse(req.ID, codeInvalidRequest, "Invalid Request", errData("invalid", "jsonrpc must be \"2.0\""))
	}
	if req.Method == "" {
		return errorResponse(req.ID, codeInvalidRequest, "Invalid Request", errData("invalid", "method is required"))
	}
	m, ok := r.byName[req.Method]
	if !ok {
		if isNotification(req.ID) {
			return nil
		}
		return errorResponse(req.ID, codeMethodNotFound, "Method not found", errData("invalid", req.Method))
	}
	params := req.Params
	if len(params) == 0 {
		params = json.RawMessage(`{}`)
	}
	result, err := m.Handle(ctx, params)
	if err != nil {
		if isNotification(req.ID) {
			return nil
		}
		code, message, data := mapError(err)
		return errorResponse(req.ID, code, message, data)
	}
	if isNotification(req.ID) {
		return nil
	}
	return mustMarshal(rpcResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      req.ID,
	})
}

func isNotification(id json.RawMessage) bool {
	return len(id) == 0
}

func errorResponse(id json.RawMessage, code int, message string, data any) []byte {
	return mustMarshal(rpcResponse{
		JSONRPC: "2.0",
		Error:   &rpcError{Code: code, Message: message, Data: data},
		ID:      id,
	})
}

func mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		fallback, _ := json.Marshal(rpcResponse{
			JSONRPC: "2.0",
			Error:   &rpcError{Code: codeInternalError, Message: "Internal error"},
			ID:      nil,
		})
		return fallback
	}
	return b
}
