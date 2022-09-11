package main

import (
	"log"
	"net"
)

func run() error {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		log.Println("Accepted connection from", conn.RemoteAddr())

		go func() {
			defer conn.Close()

			buf := make([]byte, 4096)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					return
				}
				if _, err := conn.Write(buf[:n]); err != nil {
					return
				}
			}
		}()
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
