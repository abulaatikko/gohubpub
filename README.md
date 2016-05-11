# gohubpub

TCP hub and client chat application experiment in go programming language.

## Requirements

* go

## Installation

1. git clone https://github.com/lassiheikkinen/gohubpub.git

## Usage

1. shell A: go run hub.go
2. shell B: go run client.go
3. shell C: go run client.go
4. Write commands in the shell B or C.

## Commands

### /whoami

    /whoami
    hub> 1462980334513306605

### /list

    /list
    hub> 1462980880307246656

### /msg [user_id1,user_id2,...] [message]

Shell A:

    /msg 1462980880307246656 Terve!

Shell B:

    1462980334513306605> Terve!

### /quit

    /quit
