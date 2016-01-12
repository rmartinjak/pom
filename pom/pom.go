package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"

	"github.com/rmartinjak/pom/pomrpc"
)

func main() {
	flag.Parse()
	args := flag.Args()
	cmd := ""
	if len(args) >= 1 {
		cmd = args[0]
	}

	client, err := rpc.DialHTTP("unix", pomrpc.SocketAddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()

	var call = ""
	request := pomrpc.Args{}
	var reply pomrpc.State
	switch cmd {
	case "next":
		call = "Pom.Next"
	case "set":
		if len(args) < 2 {
			log.Fatal("command \"set\" requires an argument")
		}
		request.NewState = pomrpc.State(args[1])
		call = "Pom.Set"
	default:
		call = "Pom.Get"
	}
	if err = client.Call(call, &request, &reply); err != nil {
		log.Fatal("call:", err)
	}
	fmt.Println(reply)
}
