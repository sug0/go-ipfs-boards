package main

import (
    "log"
    "flag"

    "github.com/sug0/go-ipfs-boards/boards"
)

func main() {
    var topic string
    var title string
    var thread string
    var content string

    flag.StringVar(&topic, "topic", "", "The topic of the post; equivalent to the board.")
    flag.StringVar(&title, "title", "", "The title of the post.")
    flag.StringVar(&thread, "thread", "", "The thread CID, in case of a reply post.")
    flag.StringVar(&content, "content", "", "The actual content of the post.")
    flag.Parse()

    c, err := boards.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ref, err := c.PutPost(boards.Post{
        Topic: topic,
        Title: title,
        Thread: thread,
        Content: content,
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Posted to:", ref)
}
