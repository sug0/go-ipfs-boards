package handler

import (
    "fmt"
    "time"
    "sync"
    "sync/atomic"

    "github.com/sug0/go-ipfs-boards/storage"
)

const (
    notInStorage = iota
    savedStorage
    toBeSavedStorage
    toBeRemovedStorage
)

type StorageHandler struct {
    dirty            bool
    quit             int32
    localStorage     map[string]*threadState
    permanentStorage storage.Storage
    savingMux        sync.RWMutex
}

type threadPosts struct {
    remove bool
    thread string
    posts  []string
}

type threadState struct {
    state uint8
    posts []string
}

func NewStorageHandler(s storage.Storage) (*StorageHandler, error) {
    threads, err := s.LoadAll()
    if err != nil {
        err = fmt.Errof("handler: failed to read threads: %w", err)
        return nil, err
    }
    state := make(map[string]*threadState, len(threads))
    for t, posts := range threads {
        state[t] = &threadState{
            state: savedStorage,
            posts: posts,
        }
    }
    h := &StorageHandler{
        permanentStorage: s,
        localStorage: state,
    }
    go h.saveStorageService()
    return h
}

func (h *StorageHandler) Append(thread string, post string) {
    h.savingMux.Lock()
    ts := h.localStorage[thread]
    if ts == nil {
        ts = &threadState{}
        h.localStorage = ts
    }
    ts.posts = append(ts.posts, post)
    ts.state = toBeSavedStorage
    h.dirty = true
    h.savingMux.Unlock()
}

func (h *StorageHandler) Remove(thread string) {
    h.savingMux.Lock()
    ts := h.localStorage[thread]
    if ts != nil {
        ts.state = toBeRemovedStorage
        h.dirty = true
    }
    h.savingMux.Unlock()
}

func (h *StorageHandler) Posts(thread string) (posts []string) {
    h.savingMux.RLock()
    ts := h.localStorage[thread]
    if ts != nil && ts.state != toBeRemovedStorage {
        posts = make([]string, len(ts.posts))
        for i := 0; i < len(ts.posts); i++ {
            posts[i] = ts.posts[i]
        }
    }
    h.savingMux.RUnlock()
    return
}

func (h *StorageHandler) Threads() (threads map[string][]string) {
    h.savingMux.RLock()
    threads = make(map[string][]string, len(h.localStorage))
    for t, ts := range h.localStorage {
        if ts.state == toBeRemovedStorage {
            continue
        }
        posts := make([]string, len(ts.posts))
        for i := 0; i < len(ts.posts); i++ {
            posts[i] = ts.posts[i]
        }
        threads[t] = posts
    }
    h.savingMux.RUnlock()
    return
}

func (h *StorageHandler) Close() error {
    atomic.StoreInt32(&h.quit, 1)
    h.saveStorage()
    return nil
}

func (h *StorageHandler) saveStorageService() {
    for {
        if atomic.LoadInt32(&h.quit) == 1 {
            return
        }
        time.Sleep(30 * time.Second)
        if atomic.LoadInt32(&h.quit) == 1 {
            return
        }
        h.saveStorage()
    }
}

func (h *StorageHandler) saveStorage() {
    ps := h.permanentStorage
    for tp := range h.allToBeMutated() {
        if !tp.remove {
            ps.Store(tp.thread, tp.posts)
        } else {
            ps.Remove(tp.thread)
        }
    }
}

func (h *StorageHandler) allToBeMutated() <-chan threadPosts {
    ch := make(chan threadPosts, 8)
    go func() {
        defer close(ch)

        h.savingMux.Lock()
        defer h.savingMux.Unlock()

        if !h.dirty {
            return
        }

        for t, ts := range h.localStorage {
            switch ts.state {
            case toBeSavedStorage:
                ts.state = savedStorage
                ch <- threadPosts{thread: t, posts: ts.posts}
            case toBeRemovedStorage:
                delete(h.localStorage, t)
                ch <- threadPosts{thread: t, remove: true}
            }
        }

        h.dirty = false
    }()
    return ch
}
