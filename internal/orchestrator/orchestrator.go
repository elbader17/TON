package orchestrator

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ton/framework/internal/tools"
	"github.com/ton/framework/pkg/errors"
)

type ToolRegistry struct {
	tools map[string]tools.Tool
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{tools: make(map[string]tools.Tool)}
}

func (r *ToolRegistry) Register(t tools.Tool) {
	r.tools[t.Name()] = t
}

func (r *ToolRegistry) Execute(name string, req interface{}) (interface{}, *errors.TONError) {
	tool, ok := r.tools[name]
	if !ok {
		return nil, errors.NewNotFound("tool not found: " + name)
	}
	return tool.Execute(context.Background(), req)
}

type Request struct {
	Tool   string          `json:"tool"`
	Params json.RawMessage `json:"params"`
}

type Response struct {
	Result interface{}       `json:"result"`
	Error  *errors.TONError `json:"error,omitempty"`
}

type Orchestrator struct {
	registry *ToolRegistry
}

func NewOrchestrator() *Orchestrator {
	o := &Orchestrator{registry: NewToolRegistry()}
	o.registerDefaultTools()
	return o
}

func (o *Orchestrator) registerDefaultTools() {
	o.registry.Register(tools.NewPingTool())
	o.registry.Register(tools.NewSandboxTool())
}

func (o *Orchestrator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Error: errors.NewValidation("invalid request body")})
		return
	}

	result, err := o.registry.Execute(req.Tool, req.Params)
	response := Response{Result: result}
	if err != nil {
		response.Error = err
	}
	json.NewEncoder(w).Encode(response)
}