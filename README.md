# TON Framework - Guía para Agentes IA

## Instalación

Clona el repositorio y ejecuta el script de instalación:

```bash
curl -sSL https://raw.githubusercontent.com/elbader17/TON/main/install.sh | bash
```

Esto descargará el proyecto, instalará las dependencias con `go mod tidy` y te indicará cómo levantar el servidor.

### Levantar el servidor

```bash
go run ./cmd/ton
```

El servidor estará disponible en `http://localhost:8080` con el endpoint `POST /tools`.

---

## Propósito
TON (Tool Orchestration Node) es un framework AI-First en Go diseñado como motor de ejecución para Singularity (servidor MCP).

## Principios Fundamentales
- **Cero Magia**: Sin reflection intensivo ni inyección dinámica en runtime
- **Interfaces Mínimas**: Un solo método por interfaz (`PingExecutor`, `EventPublisher`)
- **Errores Tipados**: Siempre usar `*errors.TONError` con código y mensaje LLM
- **Sin Maps Genéricos**: Comunicación entre capas mediante structs explícitos
- **Herramientas Como Primitivas**: Cada funcionalidad debe ser una Tool consumible por MCP

---

## Estructura de Directorios

```
/cmd/ton              → Punto de entrada. Servidor HTTP principal.
                      └── main.go         → Inicializa Orchestrator y levanta servidor.

/internal/domain      → Entidades core e interfaces de dominio.
                      └── tool.go        → PingRequest, PingResponse, PingExecutor interface.
                      └── sandbox.go     → SandboxRequest, SandboxResponse, SandboxExecutor.
                      └── *.go           → Agregar aquí entities e interfaces pequeñas.

/internal/tools        → Implementación de herramientas (Tools).
                      └── tool.go        → Interface Tool (Name, Description, Execute).
                      └── ping.go        → Ejemplo: PingTool + PingExecutor.
                      └── sandbox.go     → SandboxTool para ejecución aislada.
                      └── *.go           → Agregar nuevas tools aquí.

/internal/sandbox      → Implementación de ejecutores aislados.
                      └── linux_executor.go → Executor nativo con CLONE_NEWNET.

/internal/orchestrator → Routing y handlers HTTP que conectan Singularity con Tools.
                      └── orchestrator.go → ToolRegistry + Orchestrator.

/pkg/errors           → Paquete unificado de errores deterministas.
                      └── error.go       → TONError struct, códigos predefinidos.

/pkg/telemetry        → Logging y trazabilidad.
                      └── telemetry.go   → LogInfo, LogError, LogDebug.
```

---

## Reglas de Codificación

### Go Version
Go 1.21 o superior.

### Dependencias
Priorizar librería estándar. Solo agregar terceros si justifican ahorro de tokens.

### Comentarios
Todo exportable debe tener comment en formato Go estándar.

### Testing
Código en `/internal/tools` y `/internal/domain` debe ser 100% testeable con mocks explícitos.

---

## Cómo Agregar una Nueva Tool

### 1. Definir Request/Response en `/internal/domain`
```go
type MyRequest struct {
    Input string
}

type MyResponse struct {
    Output string
}
```

### 2. Crear Executor en `/internal/domain`
```go
type MyExecutor interface {
    Execute(ctx context.Context, req MyRequest) (*MyResponse, error)
}
```

### 3. Implementar Tool en `/internal/tools/my.go`
```go
type MyTool struct {
    executor domain.MyExecutor
}

func NewMyTool() *MyTool {
    return &MyTool{executor: NewMyExecutor()}
}

func (t *MyTool) Name() string        { return "my_tool" }
func (t *MyTool) Description() string { return "Description for LLM" }

func (t *MyTool) Execute(ctx context.Context, req interface{}) (interface{}, *errors.TONError) {
    // cast req a domain.MyRequest, ejecutar, devolver
}
```

### 4. Registrar en Orchestrator
```go
// internal/orchestrator/orchestrator.go
func (o *Orchestrator) registerDefaultTools() {
    o.registry.Register(tools.NewPingTool())
    o.registry.Register(tools.NewMyTool())  // AGREGAR AQUÍ
}
```

---

## Convenciones de Errores

```go
// Códigos disponibles
ErrCodeNotFound   // recurso no encontrado
ErrCodeValidation // input inválido
ErrCodeInternal   // error interno
ErrCodeExecution  // error en ejecución de tool

// Uso
errors.NewNotFound("user not found")
errors.NewValidation("invalid input")
errors.NewInternal(err)
errors.NewExecution("tool failed")
errors.Wrap(originalErr, errors.ErrCodeValidation, "context")
```

---

## Sandbox Tool

La tool `execute_in_sandbox` permite ejecutar comandos (builds, tests) en un entorno aislado con network namespace (`CLONE_NEWNET`) para prevenir exfiltración de datos.

**Nombre:** `execute_in_sandbox`

**Request:**
```json
{
    "tool": "execute_in_sandbox",
    "params": {
        "ProjectPath": "/path/to/project",
        "Command": "go build ./...",
        "TimeoutSecs": 120,
        "AllowedEnvs": ["PATH", "HOME"]
    }
}
```

**Response:**
```json
{
    "result": {
        "ExitCode": 0,
        "Stdout": "...",
        "Stderr": "",
        "Duration": "2.5s"
    },
    "error": null
}
```

** Aislamiento:**
- `CLONE_NEWNET`: Red aislada (solo localhost disponible)
- `context.WithTimeout`: Timeout configurable
- Permiso de variables de entorno controlable via `AllowedEnvs`

---

## API HTTP del Orchestrator

**Endpoint:** `POST /tools`

**Request:**
```json
{
    "tool": "ping",
    "params": {}
}
```

**Response:**
```json
{
    "result": { "Result": "pong" },
    "error": null
}
```

**Con error:**
```json
{
    "result": null,
    "error": { "code": "NOT_FOUND", "message": "tool not found: unknown" }
}
```

---

## Commands Útiles

```bash
go build ./...          # Compilar todo
go test ./...           # Ejecutar tests
go run ./cmd/ton        # Levantar servidor en :8080
```

---

## Patrón de Diseño de Tools

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│ Singularity │───▶│ Orchestrator │───▶│   Tool      │
│  (MCP)      │     │ (Router)     │     │ (Executor)  │
└─────────────┘     └──────────────┘     └─────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │   Domain    │
                                        │ (Entities)  │
                                        └─────────────┘
```

El Tool recibe el request genérico, lo castea al tipo específico del dominio, ejecuta via Executor, y devuelve resultado o error tipado.
