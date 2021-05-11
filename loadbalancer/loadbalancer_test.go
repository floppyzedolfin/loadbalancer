package loadbalancer

import (
	"testing"

	"github.com/floppyzedolfin/loadbalancer/api"
	"github.com/stretchr/testify/assert"
)

func TestMyLoadBalancer_RegisterInstance(t *testing.T) {
	tt := map[string]struct {
		instances []chan api.Request
		size      int
	}{
		"nothing": {
			instances: nil,
			size:      0,
		},
		"one instance": {
			instances: []chan api.Request{make(chan api.Request, 5)},
			size:      1,
		},
		"three instances": {
			instances: []chan api.Request{make(chan api.Request, 7), make(chan api.Request, 7), make(chan api.Request, 10)},
			size:      3,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			lb := &MyLoadBalancer{}
			for _, i := range tc.instances {
				lb.RegisterInstance(i)
			}
			assert.Equal(t, tc.size, len(lb.instances))
		})
	}
}

func TestMyLoadBalancer_Request(t *testing.T) {
	lb := &MyLoadBalancer{instances: map[string]chan api.Request{
		"channel1": make(chan api.Request, 10),
		"channel2": make(chan api.Request, 10),
		"channel3": make(chan api.Request, 10),
	}}

	// create fake responders
	fakePongs(lb)

	{
		// Send a message and check the result
		resChan := lb.Request(nil)
		// read the message
		msg := <-resChan
		// we can't know whether it was 1 or 3, but one of them answered
		assert.Regexp(t, "hello, channel[123]", msg)
	}

	// Now kill some services
	close(lb.instances["channel1"])
	close(lb.instances["channel2"])

	{
		// Send the message, check the contents, and try and see if we have removed instances
		resChan := lb.Request(nil)
		msg := <-resChan
		assert.Equal(t, "hello, channel3", msg)
		// we can't guarantee we removed instances
		assert.GreaterOrEqual(t, 3, len(lb.instances))
	}
}

func fakePongs(lb *MyLoadBalancer) {
	for id, c := range lb.instances {
		go func(chanID string, inChan chan api.Request) {
			for {
				select {
				// fake response
				case r := <-inChan:
					r.RspChan <- "hello, " + chanID
				}
			}
		}(id, c)
	}
}
