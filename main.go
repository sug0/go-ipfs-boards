package main

import (
    "fmt"

    "github.com/sug0/go-ipfs-boards/boards"
)

func main() {
    c := boards.NewClient()
    ref, err := c.PutContent("primeiro post desta tábua de texto no ipfs")
    if err != nil {
        panic(err)
    }
    post, _ := boards.NewPost("b", "isto é um teste", "", ref)
    ref, err = c.PutPost(post)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Posted to: %s\n", ref)
    post, err = c.GetPost(ref)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%#v\n", post)
}
