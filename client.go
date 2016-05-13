package main

import (
    "net"
    "fmt"
    "bufio"
    "os"
    "time"
    "bytes"
    "errors"
)

const (
    CONNECTION_TYPE = "tcp"
    CONNECTION_HOST = "localhost"
    CONNECTION_PORT = "7010"
    MAX_MESSAGE_BODY_SIZE = 1024 * 1024
    MAX_RECEIVERS = 255
    COMMAND_IDENTITY = "whoami"
    COMMAND_LIST = "list"
    COMMAND_SEND_MESSAGE = "msg"
    COMMAND_QUIT = "quit"
    COMMAND_PREFIX = "/"
)

// tells the application state (whether it's running or not)
var running bool

// supported commands by the protocol
var commands = [4]string{COMMAND_IDENTITY, COMMAND_LIST, COMMAND_SEND_MESSAGE, COMMAND_QUIT}

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

        if (!IsSupportedCommand(input)) {
            PrintCommandsList()
            continue
        }

        if (IsQuitCommand(input)) {
            running = false
        }

        if (IsSendMessageCommand(input)) {
            input = ValidateSendMessage(input)
        }

        for _, b := range input {
            writer.WriteByte(b)
        }
        writer.Flush()
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
 * The function validates the send message.
 * @param []byte input
 * @return []byte
 */
func ValidateSendMessage(input []byte) []byte {
    inputParts := bytes.SplitN(input, []byte(" "), 3)
    command := inputParts[0]
    receivers := inputParts[1]
    body := inputParts[2]

    if (len(body) > MAX_MESSAGE_BODY_SIZE) {
        errorMsg := "Message body is too large."
        HandleError(errors.New(errorMsg), errorMsg)
    }

    receiversParts := bytes.SplitN(receivers, []byte(","), MAX_RECEIVERS)
    receiversJoined := bytes.Join(receiversParts, []byte(","))

    return append(command, append([]byte(" "), append(receiversJoined, append([]byte(" "), body...)...)...)...)
}

/**
 * The function prints a command list
 */
func PrintCommandsList() {
    fmt.Println("----------------------")
    fmt.Println("Supported commands:")
    for _, c := range commands {
        fmt.Println("  " + COMMAND_PREFIX + c)
    }
    fmt.Println("----------------------")
}

/**
 * The function tells if the asked command is a supported command by the protocol.
 * @param []byte command
 * @return bool
 */
func IsSupportedCommand(command []byte) bool {
    return IsIdentityCommand(command) || IsListCommand(command) || IsSendMessageCommand(command) || IsQuitCommand(command)
}

/**
 * The function tells if the asked command is a IDENTITY command.
 * @param []byte command
 * @return bool
 */
func IsIdentityCommand(command []byte) bool {
    return IsCommand(command, COMMAND_IDENTITY)
}

/**
 * The function tells if the asked command is a LIST command.
 * @param []byte command
 * @return bool
 */
func IsListCommand(command []byte) bool {
    return IsCommand(command, COMMAND_LIST)
}

/**
 * The function tells if the asked command is a SEND_MESSAGE command.
 * @param []byte command
 * @return bool
 */
func IsSendMessageCommand(command []byte) bool {
    return IsCommand(command, COMMAND_SEND_MESSAGE)
}

/**
 * The function tells if the asked command is a QUIT command.
 * @param []byte command
 * @return bool
 */
func IsQuitCommand(command []byte) bool {
    return IsCommand(command, COMMAND_QUIT)
}

/**
 * The function tells if the asked command is a given command.
 * @param []byte command
 * @return bool
 */
func IsCommand(commandCandidate []byte, command string) bool {
    return bytes.HasPrefix(commandCandidate, []byte(COMMAND_PREFIX + command))
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

