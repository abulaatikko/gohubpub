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
    if err != nil {
        fmt.Println("Error (LISTEN): ", err.Error())
        os.Exit(1)
    }

    defer ln.Close()
    fmt.Println("Server listening " + CONNECTION_HOST + ":" + CONNECTION_PORT + "...")

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error (ACCEPT): ", err.Error())
            os.Exit(1)
        }

        message, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            fmt.Println("Error (READ): ", err.Error())
        }

        fmt.Print("Message received: ", string(message))
        conn.Write([]byte("Message is back: " + message + "\n"))

        conn.Close()
    }

}
