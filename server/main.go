package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"ipc-unix/common"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cleanup := func() {
		if _, err := os.Stat(common.SocketPath); err == nil {
			if err := os.RemoveAll(common.SocketPath); err != nil {
				log.Fatal(err)
			}
		}
	}

	cleanup()
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	listener, err := tls.Listen(common.Protocol, common.SocketAddr, &config)

	// listener, err := net.Listen(common.Protocol, common.SocketAddr)
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("ctrl-c pressed!")
		close(quit)
		cleanup()
		os.Exit(0)
	}()

	fmt.Println("> Server started")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(">>> accepted: ", conn.RemoteAddr().Network())
		go handleClient(conn)
	}
}

func getBytes(r io.Reader) <-chan chan []byte {
	ch := make(chan chan []byte)
	buf := make([]byte, 1024)
	go func() {
		n, err := r.Read(buf)
		if err == io.EOF {
			ch <- nil
			return
		}

		bufChan := make(chan []byte)
		go func() {
			bufChan <- buf[:n]
		}()
		ch <- bufChan
	}()
	return ch
}

func getData(stdout, stderr io.Reader) (*common.Data, error) {
	stdOutChan, stdErrChan := getBytes(stdout), getBytes(stderr)
	stdOutBytesChan := <-stdOutChan
	stdErrBytesChan := <-stdErrChan
	if stdOutBytesChan == nil && stdErrBytesChan == nil {
		return &common.Data{Msg: "END"}, io.EOF
	}

	data := &common.Data{}

	if stdOutBytesChan != nil {
		data.Stdout = <-stdOutBytesChan
	}

	if stdErrBytesChan != nil {
		data.Stderr = <-stdErrBytesChan
	}

	return data, nil
}

func execute(enc *json.Encoder) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "./counter/counter")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)

		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)

		return
	}

	err = cmd.Start()
	if err != nil {
		log.Println(err)

		return
	}

	go func() {
		for {
			data, err := getData(stdout, stderr)
			// fmt.Println(data)
			if err == io.EOF {
				err = enc.Encode(&common.Data{
					Msg: "END",
				})

				if err != nil {
					log.Println(err)

					return
				}

				break
			}

			if err != nil {
				log.Println(err)

				return
			}

			err = enc.Encode(data)
			if err != nil {
				if errors.Is(err, syscall.EPIPE) {
					log.Println("connection droped or probably closed by client")

					return
				}

				log.Println(err)

				return
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		log.Println(err)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	// 	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	execute(encoder)
}
