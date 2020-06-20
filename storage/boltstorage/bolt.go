package boltstorage

import (
    "os"
    "fmt"
    "bytes"
    "encoding/gob"

    bolt "go.etcd.io/bbolt"
)

type Storage struct {
    db *bolt.DB
}

func Open(path string, mode os.FileMode, options *bolt.Options) (*Storage, error) {
    db, err := bolt.Open(path, mode, options)
    if err != nil {
        err = fmt.Errorf("boltstorage: failed to open bolt db: %w", err)
        return nil, err
    }
    return &Storage{db}, nil
}

func (s *Storage) Close() error {
    err := s.db.Close()
    if err != nil {
        err = fmt.Errorf("boltstorage: failed to close bolt db: %w", err)
    }
    return err
}

func (s *Storage) LoadAll() (map[string][]string, error) {
    threads := make(map[string][]string)
    err := s.db.View(func(tx *bolt.Tx) error {
        bkt := tx.Bucket([]byte("threads"))
        if bkt == nil {
            return nil
        }
        return bkt.ForEach(func(thread, posts []byte) error {
            p, err := gobDecode(posts)
            if err != nil {
                return err
            }
            threads[string(thread)] = p
            return nil
        })
    })
    if err != nil {
        err = fmt.Errorf("boltstorage: failed to load threads: %w", err)
        return nil, err
    }
    return threads, nil
}

func (s *Storage) Store(thread string, posts []string) error {
    err := s.db.Update(func(tx *bolt.Tx) error {
        bkt, err := tx.CreateBucketIfNotExists([]byte("threads"))
        if err != nil {
            return err
        }
        data, err := gobEncode(posts)
        if err != nil {
            return err
        }
        return bkt.Put([]byte(thread), data)
    })
    if err != nil {
        err = fmt.Errorf("boltstorage: failed to save posts: %w", err)
    }
    return err
}

func (s *Storage) Load(thread string) (posts []string, err error) {
    err = s.db.View(func(tx *bolt.Tx) error {
        bkt, err := tx.CreateBucketIfNotExists([]byte("threads"))
        if err != nil {
            return err
        }
        data := bkt.Get([]byte(thread))
        if data == nil {
            return nil
        }
        posts, err = gobDecode(data)
        return err
    })
    if err != nil {
        err = fmt.Errorf("boltstorage: failed to load posts: %w", err)
    }
    return
}

func (s *Storage) Remove(thread string) error {
    err := s.db.Update(func(tx *bolt.Tx) error {
        bkt := tx.Bucket([]byte("threads"))
        if bkt == nil {
            return nil
        }
        return bkt.Delete([]byte(thread))
    })
    if err != nil {
        err = fmt.Errorf("boltstorage: failed to save posts: %w", err)
    }
    return err
}

func gobEncode(p []string) ([]byte, error) {
    var buf bytes.Buffer
    err := gob.NewEncoder(&buf).Encode(p)
    return buf.Bytes(), err
}

func gobDecode(data []byte) (p []string, err error) {
    err = gob.NewDecoder(bytes.NewReader(data)).Decode(&p)
    return
}
