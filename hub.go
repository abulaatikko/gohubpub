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
    in chan string
    out chan string
    clients []*Client
    connections chan net.Conn
    writer *bufio.Writer
}

/**
 * The function initializes communication from and to the client.
 */
func (client *Client) InitCommunication() {
    go client.Read()
    go client.Send()
}

/**
 * The function reads data from the client.
 */
func (client *Client) Read() {
    for {
        line, _ := client.reader.ReadString('\n')
        client.in <- line
    }
}

/**
 * The function sends data to the client.
 */
func (client *Client) Send() {
    for data := range client.out {
        client.writer.WriteString(data)
        client.writer.Flush()
    }
}

/**
 * The function joins client to the hub.
 *
 * @param net.Conn connection
 */
func (hub *Hub) Join(connection net.Conn) {
    client := InitClient(connection)
    hub.clients = append(hub.clients, client)

    go hub.ListenClient(client)
}

/**
 * The function decides what to do for the client requests.
 *
 * @param Client client
 */
func (hub *Hub) ListenClient(client *Client) {
    for {
        in := <-client.in
        if (strings.HasPrefix(in, "/whoami")) {
            hub.TellIdentity(client)
        } else if (strings.HasPrefix(in, "/list")) {
            hub.ListClients(client)
        } else if (strings.HasPrefix(in, "/msg")) {
            hub.SendMessage(client, in)
        } else if (strings.HasPrefix(in, "/quit")) {
            hub.UnjoinClient(client)
        }
    }
}

/**
 * The function prints data to the hub.
 *
 * @param string message
 */
func (hub *Hub) Write(message string) {
    hub.writer.WriteString(message)
    hub.writer.Flush()
}

/**
 * The function listens to the new connections and incoming data
 */
func (hub *Hub) ListenChannels() {
    for {
        select {
        case data := <-hub.in:
            hub.Write(data)
        case conn := <-hub.connections:
            hub.Join(conn)
        }
    }
}

/**
 * The function implements the /msg command.
 *
 * @param Client fromClient
 * @param string message
 */
func (hub *Hub) SendMessage(fromClient *Client, message string) {
    if (strings.Count(message, " ") <= 1) {
        fromClient.out <- "hub> Invalid /msg command parameters. Use /msg [user_id1,user_id2,...] [msg]\n"
        return
    }
    s := strings.Split(message, " ");
    receivers, body := s[1], s[2]
    r := strings.Split(receivers, ",")
    for _, client := range hub.clients {
        for _, receiver := range r {
            if (fmt.Sprintf("%d", client.user_id) == receiver) {
                client.out <- fmt.Sprintf("%d", fromClient.user_id) + "> " + body
            }
        }
    }
}

/**
 * The function implements the /list command.
 *
 * @param Client forClient
 */
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

/**
 * The function implements the /quit command.
 *
 * @param Client client
 */
func (hub *Hub) UnjoinClient(client *Client) {
    var tmpClients = make([]*Client, 0)
    for _, c := range hub.clients {
        if (c.user_id != client.user_id) {
            tmpClients = append(tmpClients, c)
        }
    }
    hub.clients = tmpClients
}

/**
 * The function implements the /whoami command.
 *
 * @param Client client
 */
func (hub *Hub) TellIdentity(client *Client) {
    client.out <- "hub> " + fmt.Sprintf("%d", client.user_id) + "\n"
}

func InitHub() *Hub {
    writer := bufio.NewWriter(os.Stdout)

    hub := &Hub{
        clients: make([]*Client, 0),
        connections: make(chan net.Conn),
        in: make(chan string),
        out: make(chan string),
        writer: writer,
    }

    go hub.ListenChannels()

    return hub
}

/**
 * The function initialize a new client.
 *
 * @return Client
 */
func InitClient(connection net.Conn) *Client {
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

    client.InitCommunication()

    return client
}

/**
 * The function handles errors.
 *
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
    fmt.Println("Server initializing...")
    hub := InitHub()

    listener, err := net.Listen(CONNECTION_TYPE, CONNECTION_HOST + ":" + CONNECTION_PORT)
    HandleError(err, "LISTEN")

    for {
        conn, err := listener.Accept()
        HandleError(err, "ACCEPT")
        hub.connections <- conn
    }
}

