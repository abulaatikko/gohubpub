# gohubpub

TCP hub and client chat application experiment in go programming language.

## Requirements

* go

## Installation

1. git clone https://github.com/lassiheikkinen/gohubpub.git

## Commands

### /whoami

### /list

### /msg [user_id1,user_id2,...] [message]

### /quit

## Use case

Start a hub (shell A) and three clients (B, C, D):

    go run hub/bin/main.go
    go run client/bin/main.go
    go run client/bin/main.go
    go run client/bin/main.go

Shell A:

    Server initializing...

Shell B:

    /whoami
    hub> 1462980334513306605
    /list
    hub> 1462980880307246656,1463330283978967334
    /msg 1462980880307246656 Hi!

Shell C:

    1462980334513306605> Hi!
    /msg 1462980334513306605,1463330283978967334 What's up guys?

Shell B:

    1462980880307246656> What's up guys?

Shell D:

    1463330283978967334> What's up guys?

Shell B, C, D:

    /quit

Shell A:

    [CTRL-C]

## Tests

    cd util; go test -v -cover; cd -
    cd client/src; go test -v -cover; cd -

## Todo

* config
* e2e testing (run hub and client virtually to test everything)
* verbose
* optimize resource usage
* protocol specification
