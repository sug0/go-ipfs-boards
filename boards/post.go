package boards

import (
    "fmt"
    "time"
    "strings"
    "unicode"
    "unicode/utf8"
)

const (
    protocol = "IPFS-TXT"
    version  = "0.1.2"

    pubsubPrefix = "/" + protocol + "/" + version + "/boards"
)

const (
    topicMaxLen   = 64
    titleMaxLen   = 256
    contentMaxLen = 1024
)

type Post struct {
    Extensions map[string]interface{} `json:",omitempty"`
    Topic      string                 `json:",omitempty"`
    Title      string                 `json:",omitempty"`
    Thread     string                 `json:",omitempty"`
    Protocol   string
    Version    string
    Content    string
    Posted     string
}

func (p *Post) validate() error {
    // validate default values
    if p.Protocol != "" {
        err := fmt.Errorf("boards: protocol field should be filled automatically")
        return err
    }
    if p.Version != "" {
        err := fmt.Errorf("boards: version field should be filled automatically")
        return err
    }
    if p.Posted != "" {
        err := fmt.Errorf("boards: posted field should be filled automatically")
        return err
    }

    // validate post kinds
    if p.Title == "" && (p.Thread == "" || p.Topic != "") {
        err := fmt.Errorf("boards: malformed reply post")
        return err
    }
    if p.Title != "" && (p.Thread != "" || p.Topic == "") {
        err := fmt.Errorf("boards: malformed op post")
        return err
    }

    // validate content length and other stuff
    p.Topic = strings.Trim(p.Topic, "/")
    var size int
    for i := 0; i < len(p.Topic); i++ {
        r, n := utf8.DecodeRuneInString(p.Topic[i:])
        if unicode.IsSpace(r) {
            err := fmt.Errorf("boards: topic mustn't contain spaces")
            return err
        }
        i += n
        size++
    }
    if size > topicMaxLen {
        err := fmt.Errorf("boards: topic length exceeded %d", topicMaxLen)
        return err
    }
    if utf8.RuneCountInString(p.Title) > titleMaxLen {
        err := fmt.Errorf("boards: title length exceeded %d", titleMaxLen)
        return err
    }
    if p.Content == "" {
        err := fmt.Errorf("boards: can't post empty content")
        return err
    }
    if utf8.RuneCountInString(p.Content) > contentMaxLen {
        err := fmt.Errorf("boards: content length exceeded %d", contentMaxLen)
        return err
    }

    // fill remaining fields
    p.Protocol = protocol
    p.Version = version
    p.Posted = timeNow()

    return nil
}

func timeNow() string {
    return time.Now().Format(time.RFC3339Nano)
}

//func timeParse(t string) (time.Time, error) {
//    return time.Parse(time.RFC3339Nano, t)
//}
