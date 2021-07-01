function getdomain() {
    var url = "config.json"
    var request = new XMLHttpRequest();
    request.open("get", url);
    request.send(null);
    request.onload = function () {
        if (request.status == 200) {
            return JSON.parse(request.responseText)["backend"];
        }
    }
}
