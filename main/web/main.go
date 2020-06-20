package main

import (
    "os"
    "log"
    "sync"
    "time"
    "math/rand"
    "net/http"
    "os/signal"

    "nhooyr.io/websocket"
    "nhooyr.io/websocket/wsjson"

    "github.com/julienschmidt/httprouter"

    "github.com/sug0/go-ipfs-boards/boards"
    "github.com/sug0/go-ipfs-boards/gossip"
    "github.com/sug0/go-ipfs-boards/storage/boltstorage"
    "github.com/sug0/go-ipfs-boards/storage/storagehandler"
)

var (
    postGossip *gossip.Gossip
    client     *boards.Client

    storageH *storagehandler.StorageHandler
)

//go:generate go run generate/resources.go

func init() {
    rand.Seed(time.Now().UnixNano())
}

func main() {
    var err error

    postGossip, err = gossip.NewGossip()
    if err != nil {
        panic(err)
    }
    defer postGossip.Close()

    client, err = boards.NewClient()
    if err != nil {
        panic(err)
    }

    st, err := boltstorage.Open(
    if err != nil {
        panic(err)
    }
    defer st.Close()

    storageH, err = storagehandler.NewStorageHandler(st)
    if err != nil {
        panic(err)
    }
    defer storageH.Close()

    go func() {
        delta := time.Duration(rand.Int() & 0xf)
        time.Sleep((30 + delta) * time.Second)
        client.AdvertiseThreads(storageH.Threads())
    }()

    sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt)

    router := httprouter.New()
    router.GET("/", indexHandler)
    router.GET("/board.js", boardScriptHandler)
    router.GET("/board.css", boardStyleHandler)
    router.GET("/thread.js", threadScriptHandler)
    router.GET("/thread.css", threadStyleHandler)
    router.GET("/boards/*board", boardsHandler)
    router.GET("/threads/:thread", threadsHandler)
    router.GET("/ws/boards/*board", wsHandlerBoards)
    router.GET("/ws/threads/:thread", wsHandlerThreads)

    go func() {
        panic(http.ListenAndServe(":8989", loggingMiddleware(router)))
    }()
    <-sig
}

func loggingMiddleware(next http.Handler) http.Handler {
    handler := func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s\t%s\t%s\n", r.RemoteAddr, r.Method, r.RequestURI)
        next.ServeHTTP(w, r)
    }
    return http.HandlerFunc(handler)
}

func indexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    http.Redirect(w, r, "/boards/newfag", http.StatusSeeOther)
}

func boardScriptHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    w.Header().Set("Content-Type", "application/javascript")
    w.Write(boardScript)
}

func boardStyleHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    w.Header().Set("Content-Type", "text/css")
    w.Write(boardStyle)
}

func boardsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    w.Write(boardsIndexPage)
}

func threadScriptHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    w.Header().Set("Content-Type", "application/javascript")
    w.Write(threadScript)
}

func threadStyleHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    w.Header().Set("Content-Type", "text/css")
    w.Write(threadStyle)
}

func threadsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    w.Write(threadsIndexPage)
}

func wsHandlerBoards(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    board := ps.ByName("board")
    if board == "/" {
        http.Error(w, "not a board", http.StatusBadRequest)
        return
    }
    c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer c.Close(websocket.StatusNormalClosure, "bye")
    board = board[1:]
    err = postGossip.AddBoardWhitelist(board)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer postGossip.DelBoardWhitelist(board)

    // the parent context
    ctx := r.Context()

    // receive posts
    type newPost struct {
        ok   bool
        post boards.Post
    }
    posts := make(chan newPost)
    go func() {
        for {
            var p newPost
            err := wsjson.Read(ctx, c, &p.post)
            if err != nil {
                p.ok = false
                posts <- p
                return
            }
            p.ok = true
            posts <- p
        }
    }()

    for {
        select {
        case p := <-posts:
            if !p.ok {
                return
            }
            go client.PutPost(p.post)
        case adv := <-postGossip.Threads():
            if adv.Topic != board {
                continue
            }
            p, err := client.GetPost(adv.Ref)
            if err != nil {
                continue
            }
            post := struct{
                Post *boards.Post
                Ref  string
            }{
                Post: p,
                Ref: adv.Ref,
            }
            go wsjson.Write(ctx, c, &post)
        }
    }
}

func wsHandlerThreads(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    thread := ps.ByName("thread")
    c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer c.Close(websocket.StatusNormalClosure, "bye")
    err = postGossip.AddThreadWhitelist(thread)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer postGossip.DelThreadWhitelist(thread)

    // the parent context
    ctx := r.Context()

    // receive posts
    type newPost struct {
        ok   bool
        post boards.Post
    }
    posts := make(chan newPost)
    go func() {
        for {
            var p newPost
            err := wsjson.Read(ctx, c, &p.post)
            if err != nil {
                p.ok = false
                posts <- p
                return
            }
            p.ok = true
            posts <- p
        }
    }()

    // write op post first
    p, err := client.GetPost(thread)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    post := struct{
        Op   bool
        Ref  string
        Post *boards.Post
    }{
        Op: true,
        Ref: thread,
        Post: p,
    }
    go wsjson.Write(ctx, c, &post)

    for {
        select {
        case p := <-posts:
            if !p.ok {
                return
            }
            go client.PutPost(p.post)
        case adv := <-postGossip.Posts():
            if adv.Thread != thread {
                continue
            }
            p, err := client.GetPost(adv.Ref)
            if err != nil {
                continue
            }
            post := struct{
                Op   bool
                Ref  string
                Post *boards.Post
            }{
                Op: false,
                Ref: adv.Ref,
                Post: p,
            }
            go wsjson.Write(ctx, c, &post)
        }
    }
}
