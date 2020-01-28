package main

import (
    "io"
    "fmt"
    "time"
    "bytes"
    "strings"
    "io/ioutil"
    "encoding/json"

    ipfs "github.com/ipfs/go-ipfs-api"
)

const (
    protocol = "IPFS-TXT"
    version  = "0.1"
)

type Client struct {
    shell *ipfs.Shell
}

type Post struct {
    Protocol string
    Version  string
    Topic    string
    Title    string `json:",omitempty"`
    Thread   string `json:",omitempty"`
    Content  string
    Posted   string
}

func main() {
    c := NewClient()
    ref, err := c.PutContent("primeiro post desta tábua de texto no ipfs")
    if err != nil {
        panic(err)
    }
    post, _ := NewPost("b", "isto é um teste", "", ref)
    ref, err = c.PutPost(post)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Posted to: %s\n", ref)
    post, err = c.GetPost(ref)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%#v\n", post)
}

func NewPost(topic, title, thread, content string) (*Post, error) {
    return &Post{
        Protocol: protocol,
        Version: version,
        Topic: topic,
        Title: title,
        Thread: thread,
        Content: content,
        Posted: timeNow(),
    }, nil
}

func NewClient() *Client {
    return &Client{ipfs.NewLocalShell()}
}

func (c *Client) PutContent(content string) (string, error) {
    ref, err := c.shell.Add(strings.NewReader(content))
    if err != nil {
        err = fmt.Errorf("boards: failed to add content to ipfs: %w", err)
        return "", err
    }
    return ref, nil
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

func (c *Client) PutPost(p *Post) (string, error) {
    var buf bytes.Buffer
    err := json.NewEncoder(&buf).Encode(p)
    if err != nil {
        err = fmt.Errorf("boards: failed to encode post into json: %w", err)
        return "", err
    }
    ref, err := c.shell.Add(&buf)
    if err != nil {
        err = fmt.Errorf("boards: failed to add post to ipfs: %w", err)
        return "", err
    }
    return ref, nil
}

func timeNow() string {
    return time.Now().Format(time.RFC3339Nano)
}

func timeParse(t string) (time.Time, error) {
    return time.Parse(time.RFC3339Nano, t)
}
