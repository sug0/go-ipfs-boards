package main

import (
    "io"
    "os"
    "net/http"
    "os/signal"

    //"nhooyr.io/websocket"
    //"nhooyr.io/websocket/wsjson"

    //"github.com/sug0/go-ipfs-boards/gossip"
    //"github.com/sug0/go-ipfs-boards/boards"
)

const indexPage = `<html>
    <head>
        <title>damn</title>
        <meta charset="utf8"/>
    </head>
    <body>
        ganda jarda mano é isso ai
        <script type="application/javascript">
            console.log('sweet');
        </script>
        <style>
            body {
                color: #ffffff;
                background: #020202;
            }
        </style>
    </body>
</html>`

func main() {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt)

    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/ws", wsHandler)

    go http.ListenAndServe(":8989", nil)
    <-sig
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, indexPage)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "for now I do nothing")
}
