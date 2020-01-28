package main

import (
    "fmt"

    "github.com/sug0/go-ipfs-boards/boards"
)

func main() {
    c, err := boards.NewClient()
    if err != nil {
        panic(err)
    }
    ref, err := c.PutPost(boards.Post{
        Thread: "QmQx3tUXcjd4YK3xLuWQEaoLu753RU7o4JgDYA4JXKRtSS",
        Content: "Hello there, I am a test reply! :D",
    })
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
