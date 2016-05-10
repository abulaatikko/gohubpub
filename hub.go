package main

import (
    "net"
    "bufio"
    "fmt"
    "os"
    "time"
    "strings"
)

const (
    CONNECTION_TYPE = "tcp"
    CONNECTION_HOST = "localhost"
    CONNECTION_PORT = "7010"
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

func (client *Client) Listen() {
    go client.Read()
    go client.Write()
}

func (client *Client) Read() {
    for {
        line, _ := client.reader.ReadString('\n')
        client.in <- line
    }
}

func (client *Client) Write() {
    for data := range client.out {
        client.writer.WriteString(data)
        client.writer.Flush()
    }
}

func (hub *Hub) Join(connection net.Conn) {
    client := CreateClient(connection)
    hub.clients = append(hub.clients, client)

    go hub.ListenClient(client)
}

func (hub *Hub) ListenClient(client *Client) {
    for {
        in := <-client.in
        if (strings.HasPrefix(in, "/whoami")) {
            client.out <- "hub> " + fmt.Sprintf("%d", client.user_id) + "\n"
        } else if (strings.HasPrefix(in, "/list")) {
            hub.ListClients(client)
        } else if (strings.HasPrefix(in, "/msg")) {
            hub.SendMessage(client, in)
        } else if (strings.HasPrefix(in, "/quit")) {
            hub.QuitClient(client)
        }
    }
}

func (hub *Hub) Write(data string) {
    hub.writer.WriteString(data)
    hub.writer.Flush()
}

func (hub *Hub) Listen() {
    for {
        select {
        case data := <-hub.in:
            hub.Write(data)
        case conn := <-hub.connections:
            hub.Join(conn)
        }
    }
}

func (hub *Hub) SendMessage(forClient *Client, message string) {
    s := strings.Split(message, " ");
    receivers, body := s[1], s[2]
    r := strings.Split(receivers, ",")
    for _, client := range hub.clients {
        for _, receiver := range r {
            if (fmt.Sprintf("%d", client.user_id) == receiver) {
                client.out <- fmt.Sprintf("%d", forClient.user_id) + "> " + body
            }
        }
    }
}

func (hub *Hub) ListClients(forClient *Client) {
    onlyMe := true
    for _, client := range hub.clients {
        if (forClient.user_id != client.user_id) {
            forClient.out <- "hub> " + fmt.Sprintf("%d", client.user_id) + "\n"
            onlyMe = false
        }
    }
    if (onlyMe == true) {
        forClient.out <- "hub> No one else here :(\n"
    }
}

func (hub *Hub) QuitClient(client *Client) {
    var tmpClients = make([]*Client, 0)
    for _, c := range hub.clients {
        if (c.user_id != client.user_id) {
            tmpClients = append(tmpClients, c)
        }
    }
    hub.clients = tmpClients
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

    go hub.Listen()

    return hub
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

func HandleError(err error, message string) {
    if (err != nil) {
        fmt.Println("ERROR (" + message + "): ", err.Error())
        os.Exit(1)
    }
}

func main() {
    fmt.Println("Server initializing...")
    hub := CreateHub()

    listener, err := net.Listen(CONNECTION_TYPE, CONNECTION_HOST + ":" + CONNECTION_PORT)
    HandleError(err, "LISTEN")

    for {
        conn, err := listener.Accept()
        HandleError(err, "ACCEPT")
        hub.connections <- conn
    }
}

