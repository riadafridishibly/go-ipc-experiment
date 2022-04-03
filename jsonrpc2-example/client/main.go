package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type clientHandler struct {
	c context.CancelFunc
}

func (h *clientHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	defer func() {
		if req.Method == "done" {
			h.c()
		}
	}()

	if req.Notif {
		log.Println("go notification params: ", string(*req.Params))
		return
	}

	log.Printf("... got unexpected request %+v\n", req)
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	ctx2, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	ctx := context.Background()
	ch := clientHandler{cancel}
	cc := jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, jsonrpc2.VarintObjectCodec{}), &ch)
	go func() {
		<-ctx2.Done()
		if err := cc.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Simple
	var got string
	// jsonrpc2.PickID(jsonrpc2.ID{Str: "STRID", IsString: true})
	if err := cc.Call(ctx, "foo/myFooFunc", []string{"hello", "world"}, &got); err != nil {
		log.Fatal(err)
	}

	log.Println("got result:", got)

	// err = cc.Notify(ctx, "notify-001", []string{"notify"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if want := fmt.Sprintf("hello, #%d: [1,2,3]", i); got != want {
	// 	log.Errorf("got result %q, want %q", got, want)
	// }
	<-cc.DisconnectNotify()
}
