package util

import (
    "bytes"
    "fmt"
    "os"
)

const (
    COMMAND_IDENTITY = "whoami"
    COMMAND_LIST = "list"
    COMMAND_SEND_MESSAGE = "msg"
    COMMAND_QUIT = "quit"
    COMMAND_PREFIX = "/"
)

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
 * The function tells if the asked command is a supported command by the protocol.
 * @param []byte command
 * @return bool
 */
func IsSupportedCommand(command []byte) bool {
    return IsIdentityCommand(command) || IsListCommand(command) || IsSendMessageCommand(command) || IsQuitCommand(command)
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
