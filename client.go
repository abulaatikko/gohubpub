package main

import (
    "net"
    "fmt"
    "bufio"
    "os"
    "time"
    "bytes"
)

const (
    CONNECTION_TYPE = "tcp"
    CONNECTION_HOST = "localhost"
    CONNECTION_PORT = "7010"
)

// tells the application state (whether it's running or not)
var running bool

// supported commands by the protocol
var commands = [4]string{"whoami", "list", "msg", "quit"}

/**
 * The function sends data from the client to the hub.
 * @param net.Conn hub
 */
func Send(hub net.Conn) {
    reader := *bufio.NewReader(os.Stdin)
    writer := *bufio.NewWriter(hub)

    for ;running; {
        input, err := reader.ReadBytes('\n')
        HandleError(err, "STDIN READ")

        if (IsSupportedCommand(input)) {
            if (bytes.HasPrefix(input, []byte("/quit"))) {
                running = false
            }
            for _, b := range input {
                writer.WriteByte(b)
            }
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

/**
 * The function reads data from the hub and prints it to the client.
 * @param net.Conn hub
 */
func Read(hub net.Conn) {
    reader := *bufio.NewReader(hub)
    writer := *bufio.NewWriter(os.Stdout)

    for ;running; {
        input, err := reader.ReadBytes('\n')
        HandleError(err, "CONNECTION READ")

        for _, b := range input {
            writer.WriteByte(b)
        }
        writer.Flush()
    }
}

/**
 * The function tells if the asked command is a supported command by the protocol.
 * @param string command
 * @return bool
 */
func IsSupportedCommand(command []byte) bool {
    for _, c := range commands {
        if (bytes.HasPrefix(command, append([]byte("/"), []byte(c)...))) {
            return true
        }
    }
    return false
}

/**
 * The function handles errors.
 * @param error err
 * @param string message
 */
func HandleError(err error, message string) {
    if (err != nil) {
        fmt.Println("ERROR (" + message + "): ", err.Error())
        os.Exit(1)
    }
}

func main() {
    running = true
    hub, err := net.Dial(CONNECTION_TYPE, CONNECTION_HOST + ":" + CONNECTION_PORT)
    HandleError(err, "DIAL")

    // close the connection when the main() returns
    defer hub.Close()

    go Send(hub)
    go Read(hub)

    for ;running; {
        time.Sleep(3600 * 24 * 7 * 365);
    }
}

