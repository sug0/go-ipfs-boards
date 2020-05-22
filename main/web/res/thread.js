let header = document.getElementById('header');
let replyboxContent = document.getElementById('replybox-content');
let replyboxSubmit = document.getElementById('replybox-submit');
let replies = document.getElementById('replies');
let thread = window.location.pathname.split('/')[2];
let ws = new WebSocket('ws://localhost:8989/ws/threads/' + thread);

header.innerText = `Thread - ${thread}`;
document.title = `Thread - ${thread}`;

ws.onmessage = e => {
    let postOp = JSON.parse(e.data);
    if (!postOp) return;
    let op = postOp.Op;
    let ref = postOp.Ref;
    let post = postOp.Post;
    let postDiv = document.createElement('div');
    postDiv.className = op ? 'op post' : 'post';
    postDiv.id = ref;
    let pHeader = document.createElement('a');
    let pDate = new Date(post.Posted);
    pHeader.className = 'post-header';
    if (post.Title) {
        pHeader.innerHTML = `<span class="post-title">${post.Title}</span><span class="post-date">${pDate.toLocaleString()}</span>`;
    } else {
        pHeader.innerHTML = `<span class="post-title">Anonymous</span><span class="post-date">${pDate.toLocaleString()}</span>`;
    }
    pHeader.href = `/threads/${thread}#${ref}`;
    let pContent = document.createElement('p');
    pContent.className = 'post-content';
    pContent.innerText = post.Content;
    postDiv.appendChild(pHeader);
    postDiv.appendChild(pContent);
    replies.appendChild(postDiv);
};

replyboxSubmit.onclick = e => {
    ws.send(JSON.stringify({
        Thread: thread,
        Content: replyboxContent.value
    }));
};

