package main

import (
    "io"
    "os"
    "log"
    "net/http"
    "os/signal"

    "nhooyr.io/websocket"
    "nhooyr.io/websocket/wsjson"

    "github.com/julienschmidt/httprouter"

    "github.com/sug0/go-ipfs-boards/boards"
    "github.com/sug0/go-ipfs-boards/gossip"
)

const boardsIndexPage = `<html>
    <head>
        <title>aew mermaum</title>
        <meta charset="utf8"/>
    </head>
    <body>
        <div id="threads"/>
        <script type="application/javascript">
            let threads = document.getElementById('threads');
            let board = window.location.pathname;
            let ws = new WebSocket('ws://localhost:8989/ws' + board);
            ws.onmessage = e => {
                let post = JSON.parse(e.data);
                if (!post) return;
                let postDiv = document.createElement('div');
                postDiv.id = 'post';
                let pHeader = document.createElement('h1');
                pHeader.id = 'post-header';
                pHeader.innerText = post.Title + ' | ' + post.Posted;
                let pContent = document.createElement('p');
                pContent.id = 'post-content';
                pContent.innerText = post.Content;
                postDiv.appendChild(pHeader);
                postDiv.appendChild(pContent);
                threads.appendChild(postDiv);
            };
        </script>
        <style>
        </style>
    </body>
</html>`

var (
    postGossip *gossip.Gossip
    client     *boards.Client
)

func main() {
    var err error

    postGossip, err = gossip.NewGossip()
    if err != nil {
        log.Fatal(err)
    }
    defer postGossip.Close()

    client, err = boards.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt)

    router := httprouter.New()
    router.GET("/boards/*board", boardsHandler)
    router.GET("/threads/:thread", threadsHandler)
    router.GET("/ws/boards/*board", wsHandlerBoards)
    router.GET("/ws/threads/:thread", wsHandlerThreads)

    go func() {
        log.Fatal(http.ListenAndServe(":8989", router))
    }()
    <-sig
}

func boardsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    io.WriteString(w, boardsIndexPage)
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
    closeCtx := c.CloseRead(ctx)

    for {
        select {
        case <-closeCtx.Done():
            return
        case adv := <-postGossip.Threads():
            p, err := client.GetPost(adv.Ref)
            if err != nil {
                continue
            }
            wsjson.Write(ctx, c, p)
        }
    }
}

func threadsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    // nothing
}

func wsHandlerThreads(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    // nothing
}
