package main

import (
    "net"
    "fmt"
    "bufio"
    "os"
    "time"
    "strings"
)

const (
    CONNECTION_TYPE = "tcp"
    CONNECTION_HOST = "localhost"
    CONNECTION_PORT = "7010"
)

var running bool
var commands = [4]string{"whoami", "list", "msg", "quit"}

func Send(conn net.Conn) {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(conn)

    for ;running; {
        input, err := reader.ReadString('\n')
        HandleError(err, "STDIN READ")

        if (IsSupportedCommand(input)) {
            if (strings.HasPrefix(input, "/quit")) {
                running = false
            }
            writer.WriteString(input)
            writer.Flush()
        } else {
            fmt.Println("----------------------")
            fmt.Println("Supported commands:")
            for _, c := range commands {
                fmt.Println("  /" + c)
            }
            fmt.Println("----------------------")
        }
    }
}

func Read(conn net.Conn) {
    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(os.Stdout)

    for ;running; {
        input, err := reader.ReadString('\n')
        HandleError(err, "CONNECTION READ")

        writer.WriteString(input)
        writer.Flush()
    }
}

func IsSupportedCommand(command string) bool {
    for _, c := range commands {
        if (strings.HasPrefix(command, "/" + c)) {
            return true
        }
    }
    return false
}

func HandleError(err error, message string) {
    if (err != nil) {
        fmt.Println("ERROR (" + message + "): ", err.Error())
        os.Exit(1)
    }
}

func main() {
    running = true
    conn, err := net.Dial(CONNECTION_TYPE, CONNECTION_HOST + ":" + CONNECTION_PORT)
    HandleError(err, "DIAL")

    // close the connection when main() returns
    defer conn.Close()

    go Send(conn)
    go Read(conn)

    for ;running; {
        time.Sleep(3600 * 24 * 7 * 365);
    }
}

