package main

import (
    "net"
    "sync"
    client "../src"
    util "../../util"
)

const (
    CONNECTION_TYPE = "tcp"
    CONNECTION_HOST = "localhost"
    CONNECTION_PORT = "7010"
)

func main() {
    conn, err := net.Dial(CONNECTION_TYPE, CONNECTION_HOST + ":" + CONNECTION_PORT)
    util.HandleError(err, "DIAL")

    // close the connection when the main() returns
    defer conn.Close()

    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        client.Send(conn)
        wg.Done()
    }()
    go func() {
        client.Read(conn)
        wg.Done()
    }()
    wg.Wait()
}
