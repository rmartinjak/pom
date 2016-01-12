package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"time"

	"github.com/rmartinjak/pom/pomrpc"
)

func wait(ch chan bool, d *time.Duration) bool {
	if d == nil {
		<-ch
		return true
	}
	log.Println("waiting", *d)
	select {
	case <-ch:
		return true
	case <-time.After(*d):
		return false
	}
}

var durations = map[pomrpc.State]*time.Duration{
	pomrpc.Work:  flag.Duration("workTime", 25*time.Minute, "amount of time to spend in \"work\" state"),
	pomrpc.Pause: flag.Duration("pauseTime", 5*time.Minute, "amount of time to spend in \"pause\" state"),
}

var optScript = ""

func runScript(s pomrpc.State) {
	if optScript == "" {
		return
	}
	cmd := exec.Command(optScript, string(s))
	if err := cmd.Start(); err != nil {
		log.Printf("command didn't start: %v", err)
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("command finished with error: %v", err)
		}
	}()
}

func init() {
	flag.StringVar(&optScript, "script", optScript, "script to run when state changes, e.g. to send a desktop notification")
}

func main() {
	flag.Parse()
	p := &pomrpc.Pom{State: pomrpc.WorkPending, Chan: make(chan bool)}
	rpc.Register(p)
	rpc.HandleHTTP()

	_ = os.Remove(pomrpc.SocketAddr)
	l, err := net.Listen("unix", pomrpc.SocketAddr)
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	go http.Serve(l, nil)

	for {
		r := wait(p.Chan, durations[p.State])
		if !r {
			pomrpc.Transition(p)
		}
		runScript(p.State)
	}
}
