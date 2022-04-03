package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

const reqMethodSep = "/"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type fooHandler struct{}

func (h *fooHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	if err := conn.Reply(ctx, req.ID,
		fmt.Sprintf("[%s] foo handler, #%s: %s",
			req.Method,
			req.ID.String(),
			string(*req.Params))); err != nil {
		log.Println("Error: ", err.Error())
	}

	time.Sleep(1 * time.Second)
	if err := conn.Notify(ctx, "progress", fmt.Sprintf("%d%% done", 25)); err != nil {
		log.Println("Error: ", err.Error())
	}

	time.Sleep(1 * time.Second)
	if err := conn.Notify(ctx, "progress", fmt.Sprintf("%d%% done", 75)); err != nil {
		log.Println("Error: ", err.Error())
	}

	time.Sleep(1 * time.Second)
	if err := conn.Notify(ctx, "progress", fmt.Sprintf("%d%% done", 100)); err != nil {
		log.Println("Error: ", err.Error())
	}

	time.Sleep(1 * time.Second)
	if err := conn.Notify(ctx, "done", fmt.Sprintf("task done")); err != nil {
		log.Println("Error: ", err.Error())
	}
}

type barHandler struct{}

func (h *barHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	if err := conn.Reply(ctx, req.ID,
		fmt.Sprintf("[%s] bar handler, #%s: %s",
			req.Method,
			req.ID.String(),
			string(*req.Params))); err != nil {
		log.Println("Error: ", err.Error())
	}

	time.Sleep(1 * time.Second)
	if err := conn.Notify(ctx, "progress", fmt.Sprintf("%d%% done", 25)); err != nil {
		log.Println("Error: ", err.Error())
	}

	time.Sleep(1 * time.Second)
	if err := conn.Notify(ctx, "progress", fmt.Sprintf("%d%% done", 75)); err != nil {
		log.Println("Error: ", err.Error())
	}

	time.Sleep(1 * time.Second)
	if err := conn.Notify(ctx, "progress", fmt.Sprintf("%d%% done", 100)); err != nil {
		log.Println("Error: ", err.Error())
	}

	time.Sleep(1 * time.Second)
	if err := conn.Notify(ctx, "done", fmt.Sprintf("task done")); err != nil {
		log.Println("Error: ", err.Error())
	}
}

type mainHandler struct {
	mu         sync.RWMutex
	handlerMap map[string]jsonrpc2.Handler
}

func (h *mainHandler) Register(workspace string, handler jsonrpc2.Handler) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.handlerMap == nil {
		h.handlerMap = make(map[string]jsonrpc2.Handler)
	}

	if workspace == "" {
		panic("workspace name can't be empty")
	}

	h.handlerMap[workspace] = handler
}

func (h *mainHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	method := req.Method

	if i := strings.Index(method, reqMethodSep); i > 0 {
		handler, ok := h.handlerMap[method[:i]]
		if ok {
			req.Method = method[i+1:]
			handler.Handle(ctx, conn, req)

			return
		}
	}

	// No specific handler, fallback to default handler
	conn.Reply(ctx, req.ID, "replying with default handler")
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {

	}

	handler := &mainHandler{}
	handler.Register("foo", &fooHandler{})
	handler.Register("bar", &barHandler{})

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Accepted conn: ", conn.RemoteAddr().String())

		jsonrpc2.NewConn(
			context.Background(),
			jsonrpc2.NewBufferedStream(conn, jsonrpc2.VarintObjectCodec{}),
			handler,
		)
	}
}
