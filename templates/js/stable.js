function queryLoginStatus() {
    fetch("/loginStatus", {
        method:"GET"
    }).then(function(response) {
        response.json().then(function(data) {
          addStableItem(data["isLogin"])
        })
    })
};

function addStableItem(loginStatus) {
    var divElement = document.createElement("div");
    divElement.className = "user_status";

    // 创建 <button> 元素
    var buttonElement = document.createElement("button");
    buttonElement.className = "button"

    // 将 <button> 元素添加到 <div> 元素中
    divElement.appendChild(buttonElement);
    buttonElement.innerText = "聊天"
    buttonElement.addEventListener("click", function() {
        window.location.href="/chat"
    })

    if (loginStatus == "true") {
        // 创建 <button> 元素
        var buttonElement = document.createElement("button");
        buttonElement.className = "button"

        // 将 <button> 元素添加到 <div> 元素中
        divElement.appendChild(buttonElement);
        buttonElement.innerText = "登出"
        buttonElement.addEventListener("click", function() {
            window.location.href="/logout"
        })

        buttonElement = document.createElement("button");
        buttonElement.className = "button"

        // 将 <button> 元素添加到 <div> 元素中
        divElement.appendChild(buttonElement);
        buttonElement.innerText = "个人信息"
        buttonElement.addEventListener("click", function() {
            window.location.href="/update"
        })
    } else {
        var currentURL = window.location.href;
        var parts = currentURL.split('/');
        var lastPart = parts[parts.length - 2];

        if (lastPart == "login") {
            var buttonElement = document.createElement("button");
            buttonElement.className = "button"

            // 将 <button> 元素添加到 <div> 元素中
            divElement.appendChild(buttonElement);
            buttonElement.innerText = "注册"
            buttonElement.addEventListener("click", function() {
                window.location.href="/register"
            })
        } else if (lastPart == "register") {
            var buttonElement = document.createElement("button");
            buttonElement.className = "button"

            // 将 <button> 元素添加到 <div> 元素中
            divElement.appendChild(buttonElement);
            buttonElement.innerText = "登陆"
            buttonElement.addEventListener("click", function() {
                window.location.href="/login"
            })
        } else {
            // 创建 <button> 元素
            var buttonElement = document.createElement("button");
            buttonElement.className = "button"
            // 将 <button> 元素添加到 <div> 元素中
            divElement.appendChild(buttonElement);
            buttonElement.innerText = "登陆"
            buttonElement.addEventListener("click", function() {
                window.location.href="/login"
            })
            buttonElement = document.createElement("button");
            buttonElement.className = "button"
            divElement.appendChild(buttonElement)
            buttonElement.innerText = "注册"
            buttonElement.addEventListener("click", function() {
                window.location.href="/register"
            })
        }
    }
    document.body.appendChild(divElement);
}

document.addEventListener("DOMContentLoaded", function(event){
    queryLoginStatus()
})

document.addEventListener('DOMContentLoaded', function() {
    var link = document.createElement('link');
    link.rel = 'icon';
    link.type = 'image/x-icon';
    link.href = '/static/favicon.ico';
    document.head.appendChild(link);
  });