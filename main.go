package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/panic", handler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello this is fun")
	})

	server := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		stop := make(chan os.Signal, 1)
		defer close(stop)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		fmt.Printf("Received signal, %s\n", (<-stop).String())
		fmt.Println("Stopping service")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			// handle err
		}
	}()

	fmt.Println("Listening at :8080")

	if err := server.ListenAndServe(); err != nil {

	}

}

type data struct {
	ID string
}

var pool = sync.Pool{New: func() interface{} { return new(data) }}

func handler(w http.ResponseWriter, r *http.Request) {

	d := pool.Get().(*data)
	defer func() {
		d.ID = ""
		pool.Put(d)
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()

	jsonStr := `{ "id":"TestID" }`
	if err := json.Unmarshal([]byte(jsonStr), &d); err != nil {
		fmt.Println("Error ", err)
		return
	}
	fmt.Fprintln(w, d)
	// panic("Force panic")
}
