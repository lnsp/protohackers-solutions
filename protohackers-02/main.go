package main

import (
	"encoding/binary"
	"log"
	"net"
)

type asset struct {
	timestamp int32
	value     int32
}

func handle(conn net.Conn) {
	defer conn.Close()

	inserts := []asset{}
	for {
		var (
			msglower, msgupper int32
			msgtype            byte
		)
		if err := binary.Read(conn, binary.BigEndian, &msgtype); err != nil {
			return
		}
		if err := binary.Read(conn, binary.BigEndian, &msglower); err != nil {
			return
		}
		if err := binary.Read(conn, binary.BigEndian, &msgupper); err != nil {
			return
		}

		switch msgtype {
		case 'I':
			inserts = append(inserts, asset{msglower, msgupper})
		case 'Q':
			var n, sum int64
			for i := range inserts {
				if msgupper < msglower || msglower > inserts[i].timestamp || inserts[i].timestamp > msgupper {
					continue
				}
				sum += int64(inserts[i].value)
				n++
			}
			var mean int32
			if n > 0 {
				mean = int32(sum / n)
			}
			binary.Write(conn, binary.BigEndian, mean)
		}
	}
}

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

		go handle(conn)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
