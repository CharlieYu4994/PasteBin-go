function getdomain() {
    var url = "config.json"
    var request = new XMLHttpRequest();
    request.open("get", url, false);
    request.send(null);
    if (request.readyState == 4) {
        return JSON.parse(request.responseText)["backend"];
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

function gettext(url) {
    var request = new XMLHttpRequest();
    var arg = getQueryString("k")
    request.open("get", url+"/get?k="+arg, false)
    request.send(null)
    if (request.readyState == 4) {
        return request.responseText
    }
}
