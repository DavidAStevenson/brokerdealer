package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/go-nats"
)

func emitEvent(nc *nats.Conn) {
	subj := "trades.booking"
	e := []byte("trade.booking.booked")

	nc.Publish(subj, e)
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subj, e)
	}
}

func bookTrade(t time.Time) {
	log.Println(t, ": Trade DONE and booked!")
}

func main() {
	log.Println("Starting trade_booking...")

	var (
		url = flag.String("url", nats.DefaultURL, "The NATS server URL to connect to")
	)
	flag.Parse()

	var nc *nats.Conn
	var err error

	maxAttempts := 10
	for attempts := 0; attempts < maxAttempts; attempts++ {
		log.Println("Attempting to connect to NATS on", *url)
		nc, err = nats.Connect(*url)
		if err == nil {
			break
		}
		log.Println(err)
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	log.Println("trade_booking started successfully.")

	// use WaitGroup to run indefinitely (until signalled to shutdown)
	wg := sync.WaitGroup{}
	wg.Add(1)
	log.Println("WaitGroup set to 1.")

	// simulate booking of trades, for starters every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func(nc *nats.Conn) {
		for {
			select {
			case t := <-ticker.C:
				bookTrade(t)
				emitEvent(nc)
			}
		}
	}(nc)

	// Gracefully shutdown on SIGINT or SIGNTERM
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Printf("Signal received: [%s]\n", sig)
		wg.Done()
	}()

	// run until a termination signal is received
	wg.Wait()

	log.Printf("trade_booking shutting down gracefully...\n")
}
