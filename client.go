package main

import (
    "net"
    "fmt"
    "bufio"
    "os"
    "time"
)

var running bool

const (
    CONNECTION_TYPE = "tcp"
    CONNECTION_HOST = "localhost"
    CONNECTION_PORT = "7010"
)

func Send(conn net.Conn) {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(conn)

    for ;running; {
        fmt.Print("> ")
        input, err := reader.ReadString('\n')
        if (err != nil) {
            fmt.Println("Error (READ): ", err.Error())
        }
        if (input == "/quit\n") {
            running = false
        }
        writer.WriteString(input)
        writer.Flush();
    }
}

func main() {
    running = true
    conn, err := net.Dial(CONNECTION_TYPE, CONNECTION_HOST + ":" + CONNECTION_PORT)
    if (err != nil) {
        fmt.Println("Error (DIAL): ", err.Error())
        os.Exit(1)
    }
    defer conn.Close()

    go Send(conn);

    for ;running; {
        time.Sleep(3600 * 24 * 7 * 365);
    }
}

