package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/floppyzedolfin/loadbalancer/api"
	"github.com/floppyzedolfin/loadbalancer/loadbalancer"
	"github.com/floppyzedolfin/loadbalancer/timeservice"
	"github.com/floppyzedolfin/loadbalancer/twig"
)

// main runs an interactive console for spawning, killing and asking for the
// time.
func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	bio := bufio.NewReader(os.Stdin)
		var lb api.LoadBalancer = &loadbalancer.MyLoadBalancer{}

	manager := &timeservice.TimeServiceManager{}

	for {
		fmt.Printf("> ")
		cmd, err := bio.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command: ", err)
			continue
		}
		switch strings.TrimSpace(cmd) {
		case "kill":
			manager.Kill()
		case "spawn":
			ts := manager.Spawn()
			lb.RegisterInstance(ts.ReqChan)
			go ts.Run()
		case "time":
			select {
			case rsp := <-lb.Request(nil):
				fmt.Println(rsp)
			case <-time.After(5 * time.Second):
				fmt.Println("Timeout")
			}
		case "exit":
			return
		case "debug":
			twig.Switch()
		default:
			fmt.Printf("Unknown command: %s Available commands: time, spawn, kill, debug, exit\n", cmd)
		}
	}
}
