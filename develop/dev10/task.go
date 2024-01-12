package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

/*
	Реализовать простейший telnet-клиент.

	Примеры вызовов:
	go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

	Требования:
	1) Программа должна подключаться к указанному хосту (ip или доменное имя + порт) по протоколу TCP.
	После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
	2) Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s)
	3) При нажатии Ctrl+D программа должна закрывать сокет и завершаться.
	Если сокет закрывается со стороны сервера, программа должна также завершаться.
	При подключении к несуществующему сервер, программа должна завершаться через timeout
*/

var (
	ErrInvalidTimeout   = errors.New("invalid timeout")
	ErrConnectionClosed = errors.New("connection is closed")
)

type Flags struct {
	Timeout time.Duration
	Address string
}

func Telnet(flg *Flags) {
	// Установка соединения
	conn, err := net.DialTimeout("tcp", flg.Address, flg.Timeout)
	if err != nil {
		log.Fatalf("failed to open connection: %v", err)
	}

	// Установка первоначального времени ожидания
	if err := conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second)); err != nil {
		log.Fatalf("failed to set read deadline: %v", err)
	}

	go func(conn net.Conn) {
		for {
			fmt.Print(">") // Запрос ввода команды пользователя

			// Отправка данных
			if _, err := io.Copy(conn, os.Stdin); err != nil {
				log.Fatalf("failed to send data to server: %v", err)
			}

			// Установка времени ожидания
			if err := conn.SetReadDeadline(time.Now().Add(time.Duration(500) * time.Millisecond)); err != nil {
				log.Fatalf("failed to set read deadline: %v", err)
			}

			// Получение данных
			if _, err := io.Copy(os.Stdout, conn); err != nil {
				switch {
				case errors.Is(err, net.ErrClosed):
					log.Fatal(ErrConnectionClosed)
				default:
					log.Fatalf("failed to receive data from server: %v", err)
				}
			}
		}
	}(conn)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit

	log.Print("shutting down")

	// Закрытие соединения
	if err := conn.Close(); err != nil {
		log.Fatalf("failed to close connection: %v", err)
	}
}

// Парсит флаги
func parseFlags() *Flags {
	var timeoutStr string
	flag.StringVar(&timeoutStr, "timeout", "10s", "connection timeout")
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatal("usage: go-telnet --timeout=<timeout> host port")
	}

	host := flag.Arg(0)
	portStr := flag.Arg(1)
	port, err := strconv.Atoi(portStr)
	if err != nil || port == 0 {
		log.Fatalf("invalid port: %s", portStr)
	}

	timeout, err := parseTimeout(timeoutStr)
	if err != nil {
		log.Fatal(err)
	}

	return &Flags{
		Address: fmt.Sprintf("%s:%d", host, port),
		Timeout: timeout,
	}
}

// Парсит и обрабатывает timeout
func parseTimeout(timeoutStr string) (time.Duration, error) {
	if len(timeoutStr) < 2 {
		return 0, ErrInvalidTimeout
	}

	valueStr := timeoutStr[:len(timeoutStr)-1]
	value, err := strconv.Atoi(valueStr)
	if err != nil || value < 0 {
		return 0, ErrInvalidTimeout
	}

	switch timeoutStr[len(timeoutStr)-1] {
	case 's':
		return time.Duration(value) * time.Second, nil
	case 'm':
		return time.Duration(value) * time.Minute, nil
	case 'h':
		return time.Duration(value) * time.Hour, nil
	default:
		return 0, ErrInvalidTimeout
	}
}

func main() {
	flg := parseFlags()
	Telnet(flg)
}
