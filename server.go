package main

import (
    "net"
    "fmt"
    "os"
    "bufio"
)

const (
    CONNECTION_TYPE = "tcp"
    CONNECTION_HOST = "localhost"
    CONNECTION_PORT = "7010"
)

func main() {
    fmt.Println("Server initializing...")
    ln, err := net.Listen(CONNECTION_TYPE, CONNECTION_HOST + ":" + CONNECTION_PORT)
    if (err != nil) {
        fmt.Println("Error (LISTEN): ", err.Error())
        os.Exit(1)
    }
    defer ln.Close()

    fmt.Println("Server listening " + CONNECTION_HOST + ":" + CONNECTION_PORT + "...")

    conn, err := ln.Accept()
    if (err != nil) {
        fmt.Println("Error (ACCEPT): ", err.Error())
        os.Exit(1)
    }

    reader := bufio.NewReader(conn)
    writerConn := bufio.NewWriter(conn)
    writerStdo := bufio.NewWriter(os.Stdout)

    for {
        message, err := reader.ReadString('\n')
        if (err != nil) {
            fmt.Println("Error (READ): ", err.Error())
            os.Exit(1)
        }

        writerStdo.WriteString("Message Received: " + message)
        writerStdo.Flush()

        writerConn.WriteString(message)
        writerConn.Flush()
    }
}
