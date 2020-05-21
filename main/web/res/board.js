let postTitle = document.getElementById('post-title');
let postContent = document.getElementById('post-content');
let postSubmit = document.getElementById('post-submit');
let threads = document.getElementById('threads');
let boardPath = window.location.pathname;
let board = boardPath.split('/').slice(2).join('/');
let ws = new WebSocket('ws://localhost:8989/ws' + boardPath);

ws.onmessage = e => {
    let postRef = JSON.parse(e.data);
    if (!postRef) return;
    let ref = postRef.Ref;
    let post = postRef.Post;
    let postDiv = document.createElement('div');
    postDiv.id = 'post';
    let pHeader = document.createElement('a');
    let pDate = new Date(post.Posted);
    pHeader.id = 'post-header';
    pHeader.innerHTML = `<span id="post-title">${post.Title}</span><span id="post-date">${pDate.toLocaleString()}</span>`;
    pHeader.href = '/threads/' + ref;
    let pContent = document.createElement('p');
    pContent.id = 'post-content';
    pContent.innerText = post.Content;
    postDiv.appendChild(pHeader);
    postDiv.appendChild(pContent);
    threads.appendChild(postDiv);
};

postSubmit.onclick = e => {
    ws.send(JSON.stringify({
        Topic: board,
        Title: postTitle.value,
        Content: postContent.value
    }));
};
