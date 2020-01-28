package boards

import (
    "fmt"
    "time"
    "unicode/utf8"
)

const (
    protocol = "IPFS-TXT"
    version  = "0.1.1"
)

const (
    topicMaxLen   = 50
    titleMaxLen   = 250
    contentMaxLen = 1500
)

type Post struct {
    Protocol string
    Version  string
    Topic    string
    Title    string `json:",omitempty"`
    Thread   string `json:",omitempty"`
    Content  string
    Posted   string
}

func newPost(topic, title, thread, content string) (*Post, error) {
    if utf8.RuneCountInString(topic) > topicMaxLen {
        err := fmt.Errorf("boards: topic length exceeded %d", topicMaxLen)
        return nil, err
    }
    if utf8.RuneCountInString(title) > titleMaxLen {
        err := fmt.Errorf("boards: title length exceeded %d", titleMaxLen)
        return nil, err
    }
    if utf8.RuneCountInString(content) > contentMaxLen {
        err := fmt.Errorf("boards: content length exceeded %d", contentMaxLen)
        return nil, err
    }
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

func timeNow() string {
    return time.Now().Format(time.RFC3339Nano)
}

//func timeParse(t string) (time.Time, error) {
//    return time.Parse(time.RFC3339Nano, t)
//}
