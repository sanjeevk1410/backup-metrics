package main

import (
        "flag"
        "fmt"
        "log"
        "math/rand"
        "net/http"
        "os"
        "os/signal"
        "syscall"
        "time"

        "github.com/braintree/manners"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
        // How often our /hello request durations fall into one of the defined buu
ckets.
        // We can use default buckets or set ones we are interested in.
        )
)

// init registers Prometheus metrics.
func init() {
        prometheus.MustRegister(duration)
        prometheus.MustRegister(counter)
}

func main() {
        addr := flag.String("http", "127.0.0.1:8000", "HTTP server address")
        flag.Parse()

        go func() {
                sigchan := make(chan os.Signal, 1)
                signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
                <-sigchan
                manners.Close()
        }()

        mux := http.NewServeMux()
        mux.HandleFunc("/hello", helloHandler)
        mux.Handle("/metrics", promhttp.Handler())
        if err := manners.ListenAndServe(*addr, mux); err != nil {
                log.Fatal(err)
        }
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
        var status int

        defer func(begun time.Time) {
                duration.Observe(time.Since(begun).Seconds())

                // hello_requests_total{status="200"} 2385
                counter.With(prometheus.Labels{
                        "status": fmt.Sprint(status),
                }).Inc()
        }(time.Now())

        status = doSomeWork()
        w.WriteHeader(status)
        w.Write([]byte("Hello, World!\n"))
}

func doSomeWork() int {
        time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

        statusCodes := [...]int{
                http.StatusOK,
                http.StatusBadRequest,
                http.StatusUnauthorized,
                http.StatusInternalServerError,
        }
        return statusCodes[rand.Intn(len(statusCodes))]
}
