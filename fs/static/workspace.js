function wsURL(id) {
    var port = window.location.port;
    var host = window.location.hostname;
    var url = "ws://" + host + ":" + port + "/watch/" + id;
    if(! url.endsWith("/")) {
        url += "/";
    }

    return url;
}

function resourceURL(resource) {
    return "/workspace/" + window.WORKSPACE_ID + "/" + resource;
}

function reload() {
    $.ajax({
        url: resourceURL("index.html"),
        dataType: "html",
        success: function(html) {
            $("#main").empty().html(html);
        },
        error: function(jqxhr, textStatus, err) {
            console.error(err);
        },
    });
}

var conn;

function reconnect() {
    conn = new WebSocket(wsURL(WORKSPACE_ID));
    if(conn) {
        conn.onopen = function() {
            console.debug("socket connection");
        };

        conn.onmessage = function(e) {
            reload();
        };
    } else {
        setTimeout(reconnect, 2000);
    }
}

function main() {
    reconnect();
    reload();
}

main();
