function getdomain() {
    var url = "config.json"
    var request = new XMLHttpRequest();
    request.open("get", url, false);
    request.send(null);
    if (request.status == 200) {
        return JSON.parse(request.responseText)["backend"];
    }
}
