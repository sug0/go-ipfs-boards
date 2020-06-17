package storage

import "github.com/sug0/go-ipfs-boards/boards"

// Stores the reply list for a particular thread.
type Threads = map[string][]string

// Allows one to implement a way to store posts
// gossiped on IPFS.
type Storage interface {
    Store(threads Threads) error
    Load() (Threads, error)
}
