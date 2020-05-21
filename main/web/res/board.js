let postTitle = document.getElementById('post-title');
let postContent = document.getElementById('post-content');
let postSubmit = document.getElementById('post-submit');
let threads = document.getElementById('threads');
let boardPath = window.location.pathname;
let board = boardPath.split('/').slice(2).join('/');
let ws = new WebSocket('ws://localhost:8989/ws' + boardPath);

ws.onmessage = e => {
    let post = JSON.parse(e.data);
    if (!post) return;
    let postDiv = document.createElement('div');
    postDiv.id = 'post';
    let pHeader = document.createElement('h1');
    pHeader.id = 'post-header';
    pHeader.innerText = post.Title + ' | ' + post.Posted;
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
