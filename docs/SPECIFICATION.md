## 🏦 Project: `gin_mcp` – Go Micro Control Plane

---

## 1. 💡 Product Summary

| Key Feature     | Description                                                  |
| --------------- | ------------------------------------------------------------ |
| Language        | Go                                                           |
| Framework       | Gin                                                          |
| Purpose         | Serve dynamically discovered model functions over HTTP       |
| Model Format    | Executable Go plugins, Python scripts, or gRPC endpoints     |
| Dynamic Routing | Auto-generates REST endpoints for models on file change      |
| Reload          | File watchers re-register routes on change                   |
| Deployment      | Local or containerised (Docker)                              |
| Constraints     | No native dynamic loading (solved via subprocess or plugins) |

---

## 2. 🧱 Architecture Overview

```
.
🔹️ /models/             # Folder for model files
🔹️ /registry/           # Internal registry of model metadata
🔹️ /handlers/           # Model execution logic
🔹️ main.go              # Gin server + route registrar
🔹️ watcher.go           # File watcher
```

---

## 3. ⚙️ Core Modules

### `main.go`

- Initializes:
  - Gin server
  - Model registry
  - Watcher (watch `./models`)
- Registers routes under `/models/:modelname`
- Proxy to handlers

---

### `watcher.go`

- Uses [`fsnotify`](https://github.com/fsnotify/fsnotify)
- On file add/change/delete:
  - Re-registers model route
  - Removes outdated ones

---

### `registry/registry.go`

- Keeps track of:
  - Model name
  - File path
  - Type 
- Maps model name to handler function

---

### `handlers/handler.go`

- Executes model:
  - If Go plugin: use `plugin.Open`
  - If Python: shell out using `exec.Command`
  - Input: JSON payload
  - Output: JSON result

---

## 4. 📅 Sample Model Formats

### `models/add.go` (compiled as `.so`)

```go
package main
import "encoding/json"

func Predict(input []byte) ([]byte, error) {
	var data map[string]int
	json.Unmarshal(input, &data)
	result := map[string]int{"sum": data["a"] + data["b"]}
	return json.Marshal(result)
}
```

### `models/predict.py`

```python
import sys, json
data = json.load(sys.stdin)
# dummy logic
json.dump({"output": data["x"] * 2}, sys.stdout)
```

---

## 5. 🔀 Routing Behavior

| Action                     | Route              |
| -------------------------- | ------------------ |
| File `models/add.go` found | POST `/models/add` |
| File removed               | Route removed      |
| File updated               | Handler reloaded   |

---

## 6. 🔁 Hot Reload Logic

- Triggered by file events (`fsnotify`)
- For `.go` files:
  - Compile to `.so` (`go build -buildmode=plugin`)
  - Load with `plugin.Open`
- For `.py`:
  - Register subprocess-based handler

---

## 7. 📄 API Contract

### Request

```json
POST /models/predict
Content-Type: application/json
{
  "x": 10
}
```

### Response

```json
{
  "output": 20
}
```

---

## 8. 📄 Sample CLI for Model Compilation

```bash
# Build Go plugin
go build -buildmode=plugin -o models/add.so models/add.go
```

---

## 9. 📦 Dependencies

| Package                        | Use                    |
| ------------------------------ | ---------------------- |
| `github.com/gin-gonic/gin`     | HTTP server            |
| `github.com/fsnotify/fsnotify` | File system monitoring |
| `os/exec`                      | Run Python models      |
| `plugin`                       | Load Go plugins        |

---

## 10. 🛠️ Tasks Breakdown for Cursor

### Initial Setup

-

### Optional Enhancements

-

