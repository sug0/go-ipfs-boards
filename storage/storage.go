package storage

// Allows one to implement a way to store posts
// gossiped on IPFS.
type Storage interface {
    LoadAll() (map[string][]string, error)
    Store(thread string, posts []string) error
    Load(thread string) ([]string, error)
    Remove(thread string) error
}
