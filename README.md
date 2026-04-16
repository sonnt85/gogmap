# gogmap

A thread-safe generic map for Go with a built-in global string map instance.

## Installation

```bash
go get github.com/sonnt85/gogmap
```

## Features

- Generic `GlobalMap[T]` type safe for concurrent access via `sync.RWMutex`
- Get with or without existence check
- Set, delete, and snapshot (copy) operations
- Pre-instantiated global `map[string]string` accessible via package-level functions

## Usage

```go
package main

import (
    "fmt"

    "github.com/sonnt85/gogmap"
)

func main() {
    // Use the built-in global string map
    gogmap.Set("key", "value")
    fmt.Println(gogmap.Get("key"))       // "value"

    v, ok := gogmap.GetVal("key")        // "value", true
    fmt.Println(v, ok)

    gogmap.Del("key")

    snapshot := gogmap.Map()             // returns a copy of the map

    // Use a typed generic map
    m := gogmap.NewGlobalMap[int]()
    m.Set("count", 42)
    count := m.Get("count")             // 42
    _ = count
    _ = snapshot
}
```

## API

### `GlobalMap[T]`

- `NewGlobalMap[T]() *GlobalMap[T]` — creates a new typed concurrent map
- `(*GlobalMap[T]).Get(key string) T` — returns value or zero value
- `(*GlobalMap[T]).GetVal(key string) (T, bool)` — returns value and existence flag
- `(*GlobalMap[T]).Set(key string, value T)` — sets a key
- `(*GlobalMap[T]).Del(key string)` — deletes a key
- `(*GlobalMap[T]).Map() map[string]T` — returns a snapshot copy

### Package-level (global `map[string]string`)

- `Get(key string) string`
- `GetVal(key string) (string, bool)`
- `Set(key, value string)`
- `Del(key string)`
- `Map() map[string]string`

## Author

**sonnt85** — [thanhson.rf@gmail.com](mailto:thanhson.rf@gmail.com)

## License

MIT License - see [LICENSE](LICENSE) for details.
