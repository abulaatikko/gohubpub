package main

import (
    "net"
    "fmt"
    hublib "../src"
    util "../../util"
)

const (
    CONNECTION_TYPE = "tcp"
    CONNECTION_HOST = "localhost"
    CONNECTION_PORT = "7010"
)

func main() {
    fmt.Println("Server initializing...")
    hub := hublib.InitHub()

    listener, err := net.Listen(CONNECTION_TYPE, CONNECTION_HOST + ":" + CONNECTION_PORT)
    util.HandleError(err, "LISTEN")

    for {
        conn, err := listener.Accept()
        util.HandleError(err, "ACCEPT")
        hub.AttachConnection(conn)
    }
}
