package main

import (
    "fmt"

    "github.com/sug0/go-ipfs-boards/boards"
)

func main() {
    c := boards.NewClient()
    ref, err := c.PutPost(
        "random",                // topic (analogous to board)
        "this is a test post",   // thread title
        "",                      // empty because this is OP
        "this is an IPFS post!", // post content
    )
    if err != nil {
        panic(err)
    }
    fmt.Printf("Posted to: %s\n", ref)
    post, err := c.GetPost(ref)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%#v\n", post)
}
