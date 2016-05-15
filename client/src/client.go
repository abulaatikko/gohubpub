package client

import (
    "net"
    "fmt"
    "bufio"
    "os"
    "bytes"
    "math"
    util "../../util"
)

const (
    MAX_MESSAGE_BODY_SIZE = 1024 * 1024
    MAX_RECEIVERS = 255
)

// supported commands by the protocol
var commands = [4]string{util.COMMAND_IDENTITY, util.COMMAND_LIST, util.COMMAND_SEND_MESSAGE, util.COMMAND_QUIT}

/**
 * The function sends data from the client to the hub.
 * @param net.Conn hub
 */
func Send(hub net.Conn) {
    reader := *bufio.NewReader(os.Stdin)
    writer := *bufio.NewWriter(hub)

    running := true
    for ;running; {
        input, err := reader.ReadBytes('\n')
        util.HandleError(err, "STDIN READ")

        if (!util.IsSupportedCommand(input)) {
            PrintCommandsList()
            continue
        }

        if (util.IsQuitCommand(input)) {
            running = false
        }

        if (util.IsSendMessageCommand(input)) {
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

    for {
        input, err := reader.ReadBytes('\n')
        util.HandleError(err, "CONNECTION READ")

        if (util.IsQuitCommand(input)) {
            break
        }

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
        newBody := []byte{}
        for _, value := range body {
            newBody = append(newBody, []byte{value}...)
            if (len(newBody) >= MAX_MESSAGE_BODY_SIZE) {
                break;
            }
        }
        body = newBody
    }

    receiversJoined := receivers
    if (bytes.Contains(receivers, []byte(","))) {
        receiversParts := bytes.SplitN(receivers, []byte(","), MAX_RECEIVERS + 1) // because SplitN returns the reminder as the last slice
        len := int(math.Min(float64(len(receiversParts)), float64(MAX_RECEIVERS)))
        receiversJoined = bytes.Join(receiversParts[:len], []byte(","))
    }

    return append(command, append([]byte(" "), append(receiversJoined, append([]byte(" "), body...)...)...)...)
}

/**
 * The function prints a command list
 */
func PrintCommandsList() {
    fmt.Println("----------------------")
    fmt.Println("Supported commands:")
    for _, c := range commands {
        fmt.Println("  " + util.COMMAND_PREFIX + c)
    }
    fmt.Println("----------------------")
}
