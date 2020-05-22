let header = document.getElementById('header');
let replyboxTitle = document.getElementById('replybox-title');
let replyboxContent = document.getElementById('replybox-content');
let replyboxSubmit = document.getElementById('replybox-submit');
let threads = document.getElementById('threads');
let boardPath = window.location.pathname;
let board = boardPath.split('/').slice(2).join('/');
let ws = new WebSocket('ws://localhost:8989/ws' + boardPath);

header.innerText = `Board - ${board}`;
document.title = `Board - ${board}`;

ws.onmessage = e => {
    let postRef = JSON.parse(e.data);
    if (!postRef) return;
    let ref = postRef.Ref;
    let post = postRef.Post;
    let postDiv = document.createElement('div');
    postDiv.className = 'post';
    let pHeader = document.createElement('a');
    let pDate = new Date(post.Posted);
    pHeader.className = 'post-header';
    pHeader.innerHTML = `<span class="post-title">${post.Title}</span><span class="post-date">${pDate.toLocaleString()}</span>`;
    pHeader.href = '/threads/' + ref;
    let pContent = document.createElement('p');
    pContent.className = 'post-content';
    pContent.innerText = post.Content;
    postDiv.appendChild(pHeader);
    postDiv.appendChild(pContent);
    threads.appendChild(postDiv);
};

replyboxSubmit.onclick = e => {
    ws.send(JSON.stringify({
        Topic: board,
        Title: replyboxTitle.value,
        Content: replyboxContent.value
    }));
    replyboxTitle.value = '';
    replyboxContent.value = '';
};
