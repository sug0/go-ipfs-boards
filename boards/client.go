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

func drainReader(r io.ReadCloser) {
    io.Copy(ioutil.Discard, r)
    r.Close()
}

func (c *Client) GetPost(ref string) (*Post, error) {
    r, err := c.shell.Cat(ref)
    if err != nil {
        err = fmt.Errorf("boards: failed to cat ipfs post: %w", err)
        return nil, err
    }
    defer drainReader(r)
    var p Post
    err = json.NewDecoder(r).Decode(&p)
    if err != nil {
        err = fmt.Errorf("boards: failed to decode post from json: %w", err)
        return nil, err
    }
    r, err = c.shell.Cat(p.Content)
    if err != nil {
        err = fmt.Errorf("boards: failed to cat ipfs post content: %w", err)
        return nil, err
    }
    defer drainReader(r)
    content, err := ioutil.ReadAll(r)
    if err != nil {
        err = fmt.Errorf("boards: failed to read ipfs post content: %w", err)
        return nil, err
    }
    p.Content = string(content)
    return &p, nil
}

func (c *Client) PutPost(p Post) (string, error) {
    err := p.validate()
    if err != nil {
        return "", err
    }
    if p.Thread != "" {
        _, _, err = c.shell.BlockStat(p.Thread)
        if err != nil {
            err = fmt.Errorf("boards: failed to get thread: %w", err)
            return "", err
        }
    }
    _, _, err = c.shell.BlockStat(p.Content)
    if err != nil {
        ref, err := c.putContent(p.Content)
        if err != nil {
            return "", err
        }
        p.Content = ref
    }
    var buf bytes.Buffer
    err = json.NewEncoder(&buf).Encode(&p)
    if err != nil {
        err = fmt.Errorf("boards: failed to encode post into json: %w", err)
        return "", err
    }
    ref, err := c.shell.Add(&buf)
    if err != nil {
        err = fmt.Errorf("boards: failed to add post to ipfs: %w", err)
        return "", err
    }
    c.advertise(p.Topic, p.Thread, ref)
    return ref, nil
}

func (c *Client) advertise(topic, thread, ref string) {
    if thread == "" {
        c.advertiseNewThread(topic, ref)
    } else {
        c.advertiseNewPost(thread, ref)
    }
}

func (c *Client) advertiseNewThread(topic, ref string) {
    subTopics := append([]string{PubsubThreadsPrefix}, strings.Split(topic, "/")...)
    for n := 1; n <= len(subTopics); n++ {
        payload := fmt.Sprintf(`{"Topic":"%s","Ref":"%s"}`, topic, ref)
        c.shell.PubSubPublish(strings.Join(subTopics[:n], "/"), payload)
    }
}

func (c *Client) advertiseNewPost(thread, ref string) {
    subTopics := []string{PubsubPostsPrefix, thread}
    for n := 1; n <= len(subTopics); n++ {
        payload := fmt.Sprintf(`{"Thread":"%s","Ref":"%s"}`, thread, ref)
        c.shell.PubSubPublish(strings.Join(subTopics[:n], "/"), payload)
    }
}

func (c *Client) putContent(content string) (string, error) {
    ref, err := c.shell.Add(strings.NewReader(content))
    if err != nil {
        err = fmt.Errorf("boards: failed to add content to ipfs: %w", err)
        return "", err
    }
    return ref, nil
}
