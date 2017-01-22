# chash
Go library for consistent hashing

## Usage
```go
package main

import "fmt"
import "github.com/mayankz/chash"

func main() {
    h := chash.New()
    
    h.AddNode("a", 1)
    h.AddNode("b", 1)
    h.AddNode("c", 1)
    fmt.Println(h.GetNode("test"))
    fmt.Println(h.GetNode("test1"))
    fmt.Println(h.GetNode("test3"))

    h.RemoveNode("b")
    fmt.Println(h.GetNode("test"))
    fmt.Println(h.GetNode("test1"))
    fmt.Println(h.GetNode("test3"))
}
```
