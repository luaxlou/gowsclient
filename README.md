# gowsclient
golang 使用 gorilla/websocket 封装client，使得websocket客户端更易用，更强壮。


这是对websocket的轻量级封装，解决一些典型问题。包含keepalive，autoReconnect,loop read message，write message串行化。以及onConnected onReceive 事件的封装，使用方法见client_test.go.
