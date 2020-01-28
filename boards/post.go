package boards

import "time"

const (
    protocol = "IPFS-TXT"
    version  = "0.1"
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

func timeNow() string {
    return time.Now().Format(time.RFC3339Nano)
}

//func timeParse(t string) (time.Time, error) {
//    return time.Parse(time.RFC3339Nano, t)
//}
