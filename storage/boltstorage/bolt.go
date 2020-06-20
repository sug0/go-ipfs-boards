package boltstorage

import (
    "fmt"
    "bytes"
    "encoding/gob"

    bolt "github.com/coreos/bbolt"
)

type Storage struct {
    db *bolt.DB
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

func gobEncode(p []string) ([]byte, error) {
    var buf bytes.Buffer
    err := gob.NewEncoder(&buf).Encode(p)
    return buf.Bytes(), err
}

func gobDecode(data []byte) (p []string, err error) {
    err = gob.NewDecoder(bytes.NewReader(data)).Decode(&p)
    return
}
