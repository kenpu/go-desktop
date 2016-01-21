File system monitoring
======================

- Monitor a file system.

- Server side:
    - a watchdog of a given directory
    - a websocket for informing the updates of files
    - a restful API for serving resources
- Client side
    - requests an index.html at a workspace.
    - subscribes to workspace watchdog
    - refreshes individual parts of the index.html using React.js

Sitemap:

    GET /workspace/:id/:path => serve static content
    WS /watch/:id/

