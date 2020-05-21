package boards

const (
    protocol = "IPFS-TXT"
    version  = "0.1.2"

    PubsubThreadsPrefix = "/" + protocol + "/" + version + "/boards"
    PubsubPostsPrefix   = "/" + protocol + "/" + version + "/threads"
)
