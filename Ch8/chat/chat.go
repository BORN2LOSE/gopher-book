// Чат-сервер, позволяющий нескольким пользователям обмениваться сообщениями.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// !+ broadcoster
type client chan<- string // Канал исходящих сообщений

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // Все входящие сообщения клиента
)

// !+ broadcaster
func broadcaster() {
	clients := make(map[client]bool) // Все подключенные клиенты
	for {
		select {
		case msg := <-messages:
			// Широковещательное входящее сообщение во все
			// каналы исходящих сообщений для клиентов.
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

// !- broadcaster

// !+ serverConnect
func serverConnect(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// Заметка: игнорируем потенциальные ошибки input.Err()

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

// !- serverConnect

// !+ clientWriter
func clientWriter(connect net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(connect, msg) // Примечание: игнорируем ошибки сети
	}
}

// !- clientWriter

// !+ main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go serverConnect(conn)
	}
}

// !- main
