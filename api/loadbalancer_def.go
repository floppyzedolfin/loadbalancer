package api

// LoadBalancer is used for balancing load between multiple instances of a service.
type LoadBalancer interface {
	Request(payload interface{}) chan Response
	RegisterInstance(chan Request)
}
