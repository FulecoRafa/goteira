package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type Ports []int64

// Set implements flag.Value.
func (p *Ports) Set(str string) error {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*p = append(*p, num)
	return nil
}

// String implements flag.Value.
func (p *Ports) String() string {
	var b strings.Builder
	for _, port := range *p {
		fmt.Fprint(&b, port)
		b.WriteString(", ")
	}
	return b.String()

}

var _ flag.Value = &Ports{}

var ports Ports

func init() {
	flag.Var(&ports, "p", "Ports to listen")
}
func Ping(port int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf(":%d was pinged\n", port)
		w.Write([]byte("Pong"))
	}
}

func httpServerFactory(wg *sync.WaitGroup, p int64) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Ping(p))

	srvr := &http.Server{
		Addr:    fmt.Sprintf(":%d", p),
		Handler: mux,
	}

	go func() {
		defer wg.Done()
		fmt.Printf("http://localhost:%d\n", p)
		srvr.ListenAndServe()
	}()

	return srvr
}

func closeSrvrs(srvrs []*http.Server) {
	for _, srvr := range srvrs {
		if err := srvr.Shutdown(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func main() {
	flag.Parse()

	fmt.Println("Serving the following ping servers:")

	var wg sync.WaitGroup
	wg.Add(len(ports))

	srvrs := make([]*http.Server, 0, len(ports))
	for _, p := range ports {
		srvr := httpServerFactory(&wg, p)
		srvrs = append(srvrs, srvr)
	}

	// Listen for interrupt
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer close(signalChan)
	go func() {
		<-signalChan
		closeSrvrs(srvrs)
	}()
	wg.Wait()
}
