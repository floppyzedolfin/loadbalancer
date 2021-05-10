package loadbalancer

import (
	"github.com/google/uuid"

	"github.com/floppyzedolfin/loadbalancer/api"
	"github.com/floppyzedolfin/loadbalancer/twig"
)

// LoadBalancer is used for balancing load between multiple instances of a service.
type LoadBalancer interface {
	Request(payload interface{}) chan api.Response
	RegisterInstance(chan api.Request)
}

// MyLoadBalancer implements the LoadBalancer interface
type MyLoadBalancer struct {
	instances map[string]chan api.Request
}

// Request sends a request to an instance and returns the channel of the response
func (lb *MyLoadBalancer) Request(payload interface{}) chan api.Response {
	req := api.Request{RspChan: make(chan api.Response, 1), Payload: payload}
	deadInstances := make([]string, 0)
	// loop over the available instances
	// we rely on the fact that Go's maps are browsed with randomness to ensure we don't spam the same instance
	for k := range lb.instances {
		// we are blind here - the best way to know whether an instance is still alive is to send it the request
		// we'll let the sendTo function take care of closed channels
		success := lb.sendTo(k, req)
		if !success {
			deadInstances = append(deadInstances, k)
			twig.Printf("instance %s appears to be dead", k)
		} else {
			twig.Printf("message sent to instance %s", k)
			// let's not try and find other dead instances -- yet
			break
		}
	}
	// cleanup - remove dead instances
	for _, deadInstance := range deadInstances {
		delete(lb.instances, deadInstance)
	}
	if len(deadInstances) > 0 {
		twig.Printf("removed instances %v, %d instances left", deadInstances, len(lb.instances))
	}

	// if no instance was found, this RspChan will never be populated
	return req.RspChan
}

// sendTo tries to send a message to an instance channel, and returns true if everything went fine
func (lb *MyLoadBalancer) sendTo(instanceKey string, req api.Request) (success bool) {
	defer func() {
		// writing to a closed channel will cause a panic. Let's catch it here
		r := recover()
		// recover returns nil if no panic was called; otherwise it will contain the panic's parameter
		success = r == nil
	}()

	// this might panic(). check the defer for the handling of that panic.
	lb.instances[instanceKey] <- req
	return
}

// RegisterInstance registers an instance to the load balancer
func (lb *MyLoadBalancer) RegisterInstance(ch chan api.Request) {
	if lb.instances == nil {
		lb.instances = make(map[string]chan api.Request)
	}
	// generate a random key
	key := uuid.New().String()
	lb.instances[key] = ch
	twig.Printf("registering instance %s", key)
	return
}
