package pomrpc

import (
	"fmt"
	"log"
	"flag"
	"sync"
	"os/user"
)

// State is the current pomodoro state
type State string

// Pom is the RPC type
type Pom struct {
	State
	Chan   chan bool
	lock   sync.RWMutex
}

// Args are the RPC arguments for Pom
type Args struct {
	NewState State
}

// States
const (
	Work         State = "work"
	PausePending       = "pause pending"
	Pause              = "pause"
	WorkPending        = "work pending"
)

// SocketAddr is the path of the unix socket used for RPC
var SocketAddr string

func init() {
	u, err := user.Current()
	if err != nil {
		log.Fatal("couldn't look up current user: %v", err)
	}
	SocketAddr = u.HomeDir + "/.pom.sock"
	flag.StringVar(&SocketAddr, "socket", SocketAddr, "path to RPC socket")
}

// Get returns the current state
func (p *Pom) Get(args *Args, reply *State) error {
	*reply = p.State
	return nil
}

// Next performs a manual state transition if there is one pending (or does nothing).
func (p *Pom) Next(args *Args, reply *State) error {
	_ = args
	if p.State == PausePending || p.State == WorkPending {
		Transition(p)
		p.Chan <- true
	}
	*reply = p.State
	return nil
}

// Set sets the current state to the one given in Args
func (p *Pom) Set(args *Args, reply *State) error {
	switch args.NewState {
	case Work:
	case WorkPending:
	case Pause:
	case PausePending:
	default:
		return fmt.Errorf("invalid state: %q", args.NewState)
	}
	p.lock.Lock()
	p.set(args.NewState)
	*reply = p.State
	p.lock.Unlock()
	p.Chan <- true
	return nil
}

func (p *Pom) set(st State) {
	p.State = st
	log.Printf("state is now %q", p.State)
}

var transitions = map[State]State{
	Work:         PausePending,
	PausePending: Pause,
	Pause:        WorkPending,
	WorkPending:  Work,
}

// Transition advances the given Pom to the next state.
// This is not a method so rpc.Register in pomd won't complain about the wrong signature.
func Transition(p *Pom) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.set(transitions[p.State])
}
