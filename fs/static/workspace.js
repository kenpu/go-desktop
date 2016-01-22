// resolves a workspace ID to a websocket URL
// for monitoring the server-side filesystem events
function wsURL(id) {
    var port = window.location.port;
    var host = window.location.hostname;
    var url = "ws://" + host + ":" + port + "/watch/" + id;

    // don't forget the trailing "/"
    if(! url.endsWith("/")) {
        url += "/";
    }

    return url;
}

// resolves a resource relative path to
// the resource URL under the workspace.
// The server will always return the file content
// verbatim.
function resourceURL(resource) {
    return "/workspace/" + 
           window.WORKSPACE_ID + "/" 
           + resource;
}

// this is an (inefficient) way of refreshing
// the #main placeholder with index.html.
// We assume that each workspace has an `index.html`
function reload() {
    $.ajax({
        url: resourceURL("index.html"), // get file
        dataType: "html",               // as html
        success: function(html) {
            $("#main").empty().html(html); // update html
        },
        error: function(jqxhr, textStatus, err) {
            console.error(err);
        },
    });
}

// this is the websocket
var conn;

// reconnect make the connection, and
// registers the event handlers for the incoming messages.
function reconnect() {
    conn = new WebSocket(wsURL(WORKSPACE_ID));
    if(conn) {
        conn.onopen = function() {
            console.debug("socket connection");
        };

        // the incoming message is a fs-event
        // in the form of [filename, event-type]
        // where event-type is "Create", "Write", ...
        // But we ignore the event, and blindly
        // refresh `index.html` using reload()
        conn.onmessage = function(e) {
            reload();
        };
    } else {
        // if the server is not connecting,
        // try again in 2 seconds.
        setTimeout(reconnect, 2000);
    }
}

// Ta-da...
function main() {
    reconnect();
    reload();
}

main();
