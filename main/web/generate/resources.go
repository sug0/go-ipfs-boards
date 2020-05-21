package main

import (
    "os"
    "log"
    "fmt"
    "bufio"
    "path/filepath"
)

func main() {
    resources := map[string]string{
        "boardsIndexPage": "board.html",
        "boardScript": "board.js",
        "boardStyle": "board.css",

        "threadsIndexPage": "thread.html",
        "threadScript": "thread.js",
        "threadStyle": "thread.css",
    }
    f, err := os.Create("resources.go")
    if err != nil {
        log.Printf("error: failed to create resources file: %s\n", err)
        return
    }
    defer f.Close()
    w := bufio.NewWriter(f)
    defer w.Flush()
    fmt.Fprintf(w, "package main;")
    for name, file := range resources {
        writeResource(w, name, file)
    }
}

func writeResource(w *bufio.Writer, name, file string) {
    fmt.Fprintf(w, "var %s=[]byte{", name)
    defer fmt.Fprintf(w, "};")
    f, err := os.Open(filepath.Join("res", file))
    if err != nil {
        log.Printf("warning: failed to read resource %s: %s\n", name, err)
        return
    }
    defer f.Close()
    r := bufio.NewReader(f)
    for {
        x, err := r.ReadByte()
        if err != nil {
            // assume EOF
            return
        }
        fmt.Fprintf(w, "%d,", x)
    }
}
