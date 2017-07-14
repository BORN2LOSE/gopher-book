// Данная программа осуществляет обратный
// отсчет для запуска ракеты.

package main

import (
	"fmt"
	"os"
	"time"
)

// !+
func main() {
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	fmt.Println("Начинаю обратный отсчет до запуска ракеты Falcon Heavy ...")
	fmt.Println("Нажмите <Enter> для прерывания:")
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-tick:
			// Ничего не делаем
		case <-abort:
			fmt.Println("Запуск отменен!")
			return
		}
	}
	launch()
}

func launch() {
	fmt.Println("Поехали !")
}

// !-
