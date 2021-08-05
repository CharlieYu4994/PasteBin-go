function getdomain() {
    var backend

    fetch("config.json").then(function (resp) {
        if (resp.status == 200) {
            backend = resp.json()["backend"];
        } else {
            backend = "api/backend"
        }
    }).catch(function (err) {
		console.error(err);
	})

    return backend
}

function getQueryString(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    var res = window.location.search.substr(1).match(reg);
    if (res != null) {
        return decodeURIComponent(res[2]);
    };
    return null;
}

function gettext(url, key) {
    var request = new XMLHttpRequest();
    request.open("get", url + "/get?k=" + key, false);
    request.send(null);
    if (request.readyState == 4) {
        if (request.status == 200) {
            return request.responseText;
        } else {
            return "查询失败";
        }
    }
}

function sendFormData(url, form) {
    var request = new XMLHttpRequest();
    var fd = new FormData(form);

    request.withCredentials = true;
    request.open("POST", url + "/add");
    request.send(fd);
    request.addEventListener("load", function () {
        if (request.status == 200) {
            window.location.href = "/get?k=" + request.responseText
        } else {
            alert("提交失败");
        }
    });
}

function del(url, key) {
    var request = new XMLHttpRequest();
    request.withCredentials = true;
    request.open("GET", url + "/del?k=" + key);
    request.send(null);
    request.addEventListener("load", function () {
        switch (request.status) {
            case 403:
                alert("这不是你上传的 Paste");
                break;
            case 500:
                alert("Paste 不存在");
                break;
            case 200:
                alert("删除成功");
                window.location.href = "/";
                break;
            default:
                alert("未知错误");
                break;
        }
    });
}
