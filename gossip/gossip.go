package gossip

import (
    "fmt"
    "sync"
    "encoding/json"

    ipfs "github.com/ipfs/go-ipfs-api"
    "github.com/sug0/go-ipfs-boards/boards"
)

type subscription struct {
    count int
    sub   *ipfs.PubSubSubscription
}

type Gossip struct {
    shell *ipfs.Shell

    boards  chan Advertisement
    threads chan Advertisement

    boardsMap  map[string]*subscription
    threadsMap map[string]*subscription

    boardsMux  sync.Mutex
    threadsMux sync.Mutex
}

func advertise(sub *ipfs.PubSubSubscription, advertChan chan<- Advertisement) {
    for {
        m, err := sub.Next()
        if err != nil {
            return
        }
        var adv Advertisement
        err = json.Unmarshal(m.Data, &adv)
        if err != nil {
            continue
        }
        advertChan <- adv
    }
}

func (g *Gossip) AddBoardWhitelist(board string) error {
    g.boardsMux.Lock()
    defer g.boardsMux.Unlock()
    if s := g.boardsMap[board]; s != nil {
        s.count++
        return nil
    }
    var topic string
    if board == "*" {
        topic = boards.PubsubThreadsPrefix
    } else {
        topic = boards.PubsubThreadsPrefix + "/" + board
    }
    sub, err := g.shell.PubSubSubscribe(topic)
    if err != nil {
        err = fmt.Errorf("gossip: failed to sub to board: %w", err)
        return err
    }
    g.boardsMap[board] = &subscription{sub: sub, count: 1}
    go advertise(sub, g.boards)
    return nil
}

func (g *Gossip) AddThreadWhitelist(thread string) error {
    g.threadsMux.Lock()
    defer g.threadsMux.Unlock()
    if s := g.threadsMap[thread]; s != nil {
        s.count++
        return nil
    }
    var topic string
    if thread == "*" {
        topic = boards.PubsubPostsPrefix
    } else {
        topic = boards.PubsubPostsPrefix + "/" + thread
    }
    sub, err := g.shell.PubSubSubscribe(topic)
    if err != nil {
        err = fmt.Errorf("gossip: failed to sub to thread: %w", err)
        return err
    }
    g.threadsMap[thread] = &subscription{sub: sub, count: 1}
    go advertise(sub, g.threads)
    return nil
}

func (g *Gossip) DelBoardWhitelist(board string) {
    g.boardsMux.Lock()
    s := g.boardsMap[board]
    s.count--
    if s.count == 0 {
        s.sub.Cancel()
        delete(g.boardsMap, board)
    }
    g.boardsMux.Unlock()
}

func (g *Gossip) DelThreadWhitelist(thread string) {
    g.threadsMux.Lock()
    s := g.threadsMap[thread]
    s.count--
    if s.count == 0 {
        s.sub.Cancel()
        delete(g.threadsMap, thread)
    }
    g.threadsMux.Unlock()
}

func (g *Gossip) Threads() <-chan Advertisement {
    return g.boards
}

func (g *Gossip) Posts() <-chan Advertisement {
    return g.threads
}

func (g *Gossip) Close() error {
    g.boardsMux.Lock()
    for k, s := range g.boardsMap {
        delete(g.boardsMap, k)
        s.sub.Cancel()
    }
    g.boardsMux.Unlock()

    g.threadsMux.Lock()
    for k, s := range g.threadsMap {
        delete(g.threadsMap, k)
        s.sub.Cancel()
    }
    g.threadsMux.Unlock()

    return nil
}

func NewGossip() (*Gossip, error) {
    sh := ipfs.NewLocalShell()
    if sh == nil {
        err := fmt.Errorf("gossip: ipfs daemon is offline")
        return nil, err
    }
    g := &Gossip{
        shell: sh,
        boards: make(chan Advertisement, 8),
        threads: make(chan Advertisement, 8),
        boardsMap: make(map[string]*subscription),
        threadsMap: make(map[string]*subscription),
    }
    return g, nil
}
