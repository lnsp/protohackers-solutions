package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
)

type request struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func isPrime(n float64) bool {
	if math.Floor(n) != n {
		return false
	}

	j := int(math.Floor(n))
	if j < 2 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(j))); i++ {
		if j%i == 0 {
			return false
		}
	}
	return true
}

func validate(req *request) bool {
	if req.Method == nil || *req.Method != "isPrime" {
		return false
	}

	if req.Number == nil {
		return false
	}

	return true
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

		go func() {
			defer conn.Close()

			scanner := bufio.NewScanner(conn)

			for scanner.Scan() {
				data := scanner.Bytes()
				req := request{}

				if err := json.Unmarshal(data, &req); err != nil {
					fmt.Fprintln(conn, string(data))
					return
				}

				// validate request
				if ok := validate(&req); !ok {
					fmt.Fprintln(conn, string(data))
					return
				}

				// check if number is prime
				prime := isPrime(*req.Number)
				data, _ = json.Marshal(&response{"isPrime", prime})

				// write response
				fmt.Fprintln(conn, string(data))
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
