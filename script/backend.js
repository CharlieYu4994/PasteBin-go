function getdomain() {
    var url = "/config.json"
    var request = new XMLHttpRequest();
    request.open("get", url);
    request.send(null);
    request.onload = function () {
        if (request.status == 200) {
            console.log(request.responseText)
            var backend = JSON.parse(request.responseText)["backend"];
            return backend
        }
    }
}
