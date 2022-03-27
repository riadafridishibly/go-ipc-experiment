package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"ipc-unix/common"
	"log"
)

func main() {
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial(common.Protocol, common.SocketAddr, &config)
	// conn, err := net.Dial(common.Protocol, common.SocketAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	decoder := json.NewDecoder(conn)

	for {
		var data common.Data
		data.Reset()
		err := decoder.Decode(&data)
		if err == io.EOF {
			fmt.Println("data from server ended")
			break
		}

		if err != nil {
			panic(err)
		}

		fmt.Println(data.String())
	}
}
