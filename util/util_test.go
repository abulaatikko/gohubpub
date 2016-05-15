package util_test

import (
    util "../util"
    "testing"
)

func TestIsQuitCommand(t *testing.T) {
    command := []byte("/quit")
    if (!util.IsQuitCommand(command)) {
        t.Error("/quit IsQuitCommand")
    }

    command = []byte("/list")
    if (util.IsQuitCommand(command)) {
        t.Error("/list isnt IsQuitCommand")
    }

    command = []byte("/ quit")
    if (util.IsQuitCommand(command)) {
        t.Error("/ quit isnt IsQuitCommand")
    }

    command = []byte("/quit something extra doesnt affect")
    if (!util.IsQuitCommand(command)) {
        t.Error("/quit something extra is IsQuitCommand")
    }
}

func TestIsSendMessageCommand(t *testing.T) {
    command := []byte("/msg 1234,2345 What's up?")
    if (!util.IsSendMessageCommand(command)) {
        t.Error("IsSendMessageCommand")
    }
}

func TestIsIdentityCommand(t *testing.T) {
    command := []byte("/whoami")
    if (!util.IsIdentityCommand(command)) {
        t.Error("IsIdentityCommand")
    }
}

func TestIsListCommand(t *testing.T) {
    command := []byte("/list")
    if (!util.IsListCommand(command)) {
        t.Error("IsListCommand")
    }
}

func TestIsSupportedCommand(t *testing.T) {
    command := []byte("/quit")
    if (!util.IsSupportedCommand(command)) {
        t.Error("IsSupportedCommand")
    }
    command = []byte("/list")
    if (!util.IsSupportedCommand(command)) {
        t.Error("IsSupportedCommand")
    }
    command = []byte("/msg 2345 Terve!")
    if (!util.IsSupportedCommand(command)) {
        t.Error("IsSupportedCommand")
    }
    command = []byte("/quit")
    if (!util.IsSupportedCommand(command)) {
        t.Error("IsSupportedCommand")
    }
}