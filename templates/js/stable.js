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

    if (loginStatus == "true") {
        // 创建 <button> 元素
        var buttonElement = document.createElement("button");

        // 将 <button> 元素添加到 <div> 元素中
        divElement.appendChild(buttonElement);
        buttonElement.innerText = "登出"
        buttonElement.addEventListener("click", function() {
            window.location.href="/logout"
        })
    } else {
        // 创建 <button> 元素
        var buttonElement = document.createElement("button");
        // 将 <button> 元素添加到 <div> 元素中
        divElement.appendChild(buttonElement);
        buttonElement.innerText = "登陆"
        buttonElement.addEventListener("click", function() {
            window.location.href="/login"
        })
        buttonElement = document.createElement("button");
        divElement.appendChild(buttonElement)
        buttonElement.innerText = "注册"
        buttonElement.addEventListener("click", function() {
            window.location.href="/registe"
        })
    }
    document.body.appendChild(divElement);
}

document.addEventListener("DOMContentLoaded", function(event){
    queryLoginStatus()
})