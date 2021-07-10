function getdomain() {
    var url = "config.json"
    var request = new XMLHttpRequest();
    request.open("get", url, false);
    request.send(null);
    if (request.readyState == 4) {
        if (request.status == 200) {
            return JSON.parse(request.responseText)["backend"];
        }
        else {
            return "api/add";
        }
    }
}

function getQueryString(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    var res = window.location.search.substr(1).match(reg);
    if (res != null) {
        return decodeURIComponent(res[2]);
    };
    return null;
}

function gettext(url, arg) {
    var request = new XMLHttpRequest();
    request.open("get", url + "/get?k=" + arg, false);
    request.send(null);
    if (request.readyState == 4) {
        if (request.status == 200) {
            return request.responseText;
        }
        else {
            return "查询失败";
        }
    }
}
