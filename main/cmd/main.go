package main

import (
    "os"
    "log"
    "flag"
    "os/signal"

    "github.com/sug0/go-ipfs-boards/gossip"
    "github.com/sug0/go-ipfs-boards/boards"
)

func main() {
    var snoop bool
    var topic string
    var title string
    var thread string
    var content string

    flag.BoolVar(&snoop, "gossip", false, "Read the post gossip in IPFS.")
    flag.StringVar(&topic, "topic", "", "The topic of the post; equivalent to the board.")
    flag.StringVar(&title, "title", "", "The title of the post.")
    flag.StringVar(&thread, "thread", "", "The thread CID, in case of a reply post.")
    flag.StringVar(&content, "content", "", "The actual content of the post.")
    flag.Parse()

    if snoop {
        snoopPosts()
        return
    }

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

func snoopPosts() {
    g, err := gossip.NewGossip()
    if err != nil {
        log.Fatal(err)
    }
    defer g.Close()
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    for {
        select {
        case <-quit:
            return
        case a := <-g.NextThread():
            log.Printf("New thread: %s: On topic: %s\n", a.Ref, a.Topic)
        case a := <-g.NextPost():
            log.Printf("New post: %s: On thread: %s\n", a.Ref, a.Thread)
        }
    }
}
