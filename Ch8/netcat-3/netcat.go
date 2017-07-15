package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // Игнорируем ошибки
		log.Println("done")
		done <- struct{}{} // Сигнал главной go-подпрограмме
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done
	// Ожидание завершения фонового горутина
}

func mustCopy(iow io.Writer, src io.Reader) {
	if _, err := io.Copy(iow, src); err != nil {
		log.Fatal(err)
	}
}
