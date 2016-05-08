package main

import (
    "net"
    "bufio"
    "fmt"
    "os"
    "time"
)

type Client struct {
    in chan string
    out chan string
    user_id uint64
    reader *bufio.Reader
    writer *bufio.Writer
}

type Hub struct {
    clients []*Client
    connections chan net.Conn
    in chan string
    out chan string
    writer *bufio.Writer
}

const (
    CONNECTION_TYPE = "tcp"
    CONNECTION_HOST = "localhost"
    CONNECTION_PORT = "7010"
)

func (client *Client) Read() {
    for {
        line, err := client.reader.ReadString('\n')
        if (err != nil) {
            fmt.Println("Error (READ): ", err.Error())
            os.Exit(1)
        }
        client.in <- line
    }
}

func (client *Client) Write() {
    for data := range client.out {
        client.writer.WriteString(data)
        client.writer.Flush()
    }
}

func (client *Client) Listen() {
    go client.Read()
    go client.Write()
}

func CreateClient(connection net.Conn) *Client {
    writer := bufio.NewWriter(connection)
    reader := bufio.NewReader(connection)

    user_id := uint64(time.Now().UnixNano())

    client := &Client{
        in: make(chan string),
        out: make(chan string),
        reader: reader,
        writer: writer,
        user_id: user_id,
    }

    client.Listen()

    return client
}

func (hub *Hub) Broadcast(data string) {
    for _, client := range hub.clients {
        client.out <- data
    }
}

func (hub *Hub) Join(connection net.Conn) {
    client := CreateClient(connection)
    hub.clients = append(hub.clients, client)
    go func() {
        for {
            hub.in <- <-client.in
        }
    }()
    client.writer.WriteString("Your user_id: " + fmt.Sprintf("%d", client.user_id))
    client.writer.Flush()
}

func (hub *Hub) Write(data string) {
    hub.writer.WriteString(data)
    hub.writer.Flush()
}

func (hub *Hub) Listen() {
    go func() {
        for {
            select {
            case data := <-hub.in:
                hub.Broadcast(data)
                hub.Write(data)
            case conn := <-hub.connections:
                hub.Join(conn)
            }
        }
    }()
}

func CreateHub() *Hub {
    writer := bufio.NewWriter(os.Stdout)

    hub := &Hub{
        clients: make([]*Client, 0),
        connections: make(chan net.Conn),
        in: make(chan string),
        out: make(chan string),
        writer: writer,
    }

    hub.Listen()

    return hub
}

func main() {
    fmt.Println("Server initializing...")
    hub := CreateHub()

    listener, err := net.Listen(CONNECTION_TYPE, CONNECTION_HOST + ":" + CONNECTION_PORT)
    if (err != nil) {
        fmt.Println("Error (LISTEN): ", err.Error())
        os.Exit(1)
    }

    for {
        conn, err := listener.Accept()
        if (err != nil) {
            fmt.Println("Error (LISTEN): ", err.Error())
            os.Exit(1)
        }
        hub.connections <- conn
    }
}

