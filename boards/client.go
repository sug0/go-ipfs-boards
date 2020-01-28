package boards

import (
    "io"
    "fmt"
    "bytes"
    "strings"
    "io/ioutil"
    "encoding/json"

    ipfs "github.com/ipfs/go-ipfs-api"
)

type Client struct {
    shell *ipfs.Shell
}

func NewClient() (*Client, error) {
    sh := ipfs.NewLocalShell()
    if sh == nil {
        err := fmt.Errorf("boards: ipfs daemon is offline")
        return nil, err
    }
    return &Client{sh}, nil
}

func (c *Client) GetPost(ref string) (*Post, error) {
    r, err := c.shell.Cat(ref)
    if err != nil {
        err = fmt.Errorf("boards: failed to cat ipfs post: %w", err)
        return nil, err
    }
    defer func() {
        io.Copy(ioutil.Discard, r)
        r.Close()
    }()
    var p Post
    err = json.NewDecoder(r).Decode(&p)
    if err != nil {
        err = fmt.Errorf("boards: failed to decode post from json: %w", err)
        return nil, err
    }
    return &p, nil
}

func (c *Client) PutPost(topic, title, thread, content string) (string, error) {
    p, err := newPost(topic, title, thread)
    if err != nil {
        return "", err
    }
    ref, err := c.putContent(content)
    if err != nil {
        return "", err
    }
    p.Content = ref
    var buf bytes.Buffer
    err = json.NewEncoder(&buf).Encode(p)
    if err != nil {
        err = fmt.Errorf("boards: failed to encode post into json: %w", err)
        return "", err
    }
    ref, err = c.shell.Add(&buf)
    if err != nil {
        err = fmt.Errorf("boards: failed to add post to ipfs: %w", err)
        return "", err
    }
    return ref, nil
}

func (c *Client) putContent(content string) (string, error) {
    ref, err := c.shell.Add(strings.NewReader(content))
    if err != nil {
        err = fmt.Errorf("boards: failed to add content to ipfs: %w", err)
        return "", err
    }
    return ref, nil
}
