package gossip

import (
    "fmt"
    "encoding/json"

    ipfs "github.com/ipfs/go-ipfs-api"
    "github.com/sug0/go-ipfs-boards/boards"
)

type Advertisement struct {
    Thread string
    Topic  string
    Ref    string
}

type Gossip struct {
    bs      *ipfs.PubSubSubscription
    bt      *ipfs.PubSubSubscription
    boards  chan Advertisement
    threads chan Advertisement
}

func (g *Gossip) NextThread() <-chan Advertisement {
    return g.boards
}

func (g *Gossip) NextPost() <-chan Advertisement {
    return g.threads
}

func (g *Gossip) Next() Advertisement {
    select {
    case adv := <-g.boards:
        return adv
    case adv := <-g.threads:
        return adv
    }
}

func (g *Gossip) Close() error {
    go g.bs.Cancel()
    go g.bt.Cancel()
    return nil
}

func NewGossip() (*Gossip, error) {
    sh := ipfs.NewLocalShell()
    if sh == nil {
        err := fmt.Errorf("gossip: ipfs daemon is offline")
        return nil, err
    }
    bs, err := sh.PubSubSubscribe(boards.PubsubThreadsPrefix)
    if err != nil {
        err := fmt.Errorf("gossip: failed to sub to boards: %w", err)
        return nil, err
    }
    bt, err := sh.PubSubSubscribe(boards.PubsubPostsPrefix)
    if err != nil {
        go bs.Cancel()
        err := fmt.Errorf("gossip: failed to sub to threads: %w", err)
        return nil, err
    }
    g := &Gossip{
        bs: bs,
        bt: bt,
        boards: make(chan Advertisement, 4),
        threads: make(chan Advertisement, 4),
    }
    advertise := func(sub *ipfs.PubSubSubscription, advertChan chan<- Advertisement) {
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
    go advertise(g.bs, g.boards)
    go advertise(g.bt, g.threads)
    return g, nil
}
