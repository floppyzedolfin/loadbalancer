package timeservice

import (
	"math/rand"
	"testing"

	"github.com/floppyzedolfin/loadbalancer/api"
	"github.com/stretchr/testify/assert"
)

func TestTimeService_Run(t *testing.T) {
	// let's make sure that, when we kill the service, this properly closes the ReqChan
	ts := TimeService{
		Dead:            make(chan struct{}, 0),
		ReqChan:         make(chan api.Request, 10),
		AvgResponseTime: rand.Float64() * 3,
	}

	// launch it
	go ts.Run()

	// kill it
	ts.Dead <- struct{}{}

	// check the ReqChan - it's closed, which means reading from it returns false (writing to it would panic)
	_, ok := <- ts.ReqChan
	assert.False(t, ok)
}
