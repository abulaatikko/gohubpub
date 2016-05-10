package main

import (
    "net"
    "fmt"
    "bufio"
    "os"
    "time"
    "strings"
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
        input, err := reader.ReadString('\n')
        if (err != nil) {
            fmt.Println("Error (STDIN READ): ", err.Error())
            os.Exit(1)
        }
        if (strings.HasPrefix(input, "/list")) {
        } else if (strings.HasPrefix(input, "/msg")) {
        } else if (strings.HasPrefix(input, "/whoami")) {
        } else if (strings.HasPrefix(input, "/quit")) {
            running = false
        } else {
            fmt.Println("Supported commands:")
            fmt.Println("  /whoami")
            fmt.Println("  /list")
            fmt.Println("  /msg [user_id1,user_id2,...] [message]")
            fmt.Print("> ")
            continue
        }
        writer.WriteString(input)
        writer.Flush()
    }
}

func Read(conn net.Conn) {
    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(os.Stdout)

    for ;running; {
        input, err := reader.ReadString('\n')
        if (err != nil) {
            fmt.Println("ERROR (CONNECTION READ): ", err.Error())
            os.Exit(1)
        }
        writer.WriteString(input)
        writer.WriteString("> ")
        writer.Flush()
    }
}

func main() {
    running = true
    conn, err := net.Dial(CONNECTION_TYPE, CONNECTION_HOST + ":" + CONNECTION_PORT)
    if (err != nil) {
        fmt.Println("Error (DIAL): ", err.Error())
        os.Exit(1)
    }

    // close the connection when main() returns
    defer conn.Close()

    fmt.Print("> ")

    go Send(conn)
    go Read(conn)

    for ;running; {
        time.Sleep(3600 * 24 * 7 * 365);
    }
}

