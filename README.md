# Load Balancer

## Get it

Use the following command to get the contents of the repo:
```
git clone https://github.com/floppyzedolfin/loadbalancer.git
```

## Try it!

Here are a couple of tests one can run to ensure this works as expected:

- One instance of the service returns a message

```
$ make run
> time
Timeout
> spawn
> time
2021-05-11 01:26:42.084701526 +0200 CEST m=+18.723647170
> exit
```

- One instance, two consecutive calls

```
$ make run
> spawn
> time
2021-05-11 01:32:44.566558756 +0200 CEST m=+4.089792811
> time
2021-05-11 01:32:45.826326638 +0200 CEST m=+5.349560917
> exit
```

- Two instances of the service return one message

```
$ make run
> spawn
> spawn
> time
2021-05-11 01:28:27.110303832 +0200 CEST m=+5.322090573
> exit
```

- Two instances of the service down to one instance returns a message

```
$ make run
> spawn
> spawn
> kill
> time
2021-05-11 01:28:41.33827261 +0200 CEST m=+4.875787782
> exit
```

- Two instances of the service, down to none, returns a timeout

```
$ make run
> spawn
> spawn
> kill
> kill
> time
Timeout
> exit
```

- Debug traces

```
$ make run
> spawn
> debug
debug ON
> spawn
registering instance 3a8948e0-70b8-4bfd-80e1-1c01dbe4341a
> time
message sent to instance 7a09258a-e7de-4b4a-99cb-acf7b1b48e90
2021-05-11 01:30:02.175815406 +0200 CEST m=+8.896377778
> kill
> spawn
registering instance debe63e1-9a4e-4592-97f0-5f741a446030
> kill
> time
instance debe63e1-9a4e-4592-97f0-5f741a446030 appears to be dead
message sent to instance 7a09258a-e7de-4b4a-99cb-acf7b1b48e90
removed instance(s) [debe63e1-9a4e-4592-97f0-5f741a446030], 2 instance(s) left
2021-05-11 01:30:24.272103558 +0200 CEST m=+30.992665852
> time
message sent to instance 7a09258a-e7de-4b4a-99cb-acf7b1b48e90
2021-05-11 01:31:45.837742061 +0200 CEST m=+112.558304345
> time
instance 3a8948e0-70b8-4bfd-80e1-1c01dbe4341a appears to be dead
message sent to instance 7a09258a-e7de-4b4a-99cb-acf7b1b48e90
removed instance(s) [3a8948e0-70b8-4bfd-80e1-1c01dbe4341a], 1 instance(s) left
2021-05-11 01:31:46.824926293 +0200 CEST m=+113.545488580
> debug
> time
2021-05-11 01:31:47.824926293 +0200 CEST m=+114.545488580
> exit
```

## Discussion

### Change scope

My main target here was to not alter existing code (apart from moving it to
separate packages). I've achieved this by implementing everything in
`loadbalancer`, but I needed to write something in `TimeService`.

### Transfering the death of the service information

The first problem I faced was indeed that, when we kill a service, there is no
way to tell the load balancer that the instance was killed. For instance, other
implementations could have instance IDs, and when we randomly kill one, the ID
of the killed instance is returned. Implementing this could have been done, but
it would have violated my rule #1 of not altering existing code
(in this case, the `Kill` function). I therefore needed to transfer the shutdown
of the TimeService to its channel. This was achieved by adding a
`close(ts.ReqChan)` call when we detect that the TimeService has been killed.

### MyLoadBalancer structure

I've decided to store my instances inside the `MyLoadBalancer` as a map indexed
by a key rather than a `[]chan api.Request`. The main reason for doing so was
that, in Golang, maps are randomly iterated through. I didn't want to always be
shooting messages to the same channel "because it's the first". Also, using a
map made the cleaning of the dead instances easier - I needed to call `delete`
for each entry rather than doing some uncomfortable slice operations (like
iterating through a slice whilst removing some of its elements).

### Detecting dead services

The contract I've decided to have both the services and the load balancer agree
upon is that a service needs to close its `ReqChan` to let the rest of the world
know it's dead. This allowed me to detect that a service was dead by trying to
write to its channel. Indeed, in Golang, writing to a closed channel raises
a `panic` (while writing to a non-closed channel never does). That's how I've
decided to acknowledge the death of a service on the LoadBalancer side.
Using `recover` around that channel enqueueing prevents the program from
exploding the hands of the user.

### Dealing with dead services

Once a channel has been detected as faulty, we add it to a list of "things 
to clean up". After trying the sending of a message, we clean the dead 
instances from the LoadBalancer "database" (the map in the memory). This 
ensures that further calls to `Request` won't try and send messages to these 
dead instances.

However, it's to be noticed that, should we succeed on sending a message on 
the first attempt, we won't check for other dead services.

### Code structure

I've reorganized the code in different packages, each one with its own purpose.
The packages I've worked on are

- `loadbalancer`: implementation of the task
- `timeservice`: I had to add a line inside `TimeService.Run`
- `twig`: see below

## Added utilities

### Twig

During the implementation / debugging of my code, I needed to be able to print
information. As I didn't want to pollute the standard output, I wrote a
tiny `twig` utility (it's a small `log` package). Since I really didn't want to
make it too big, I used a global (local) variable to set its state. A "better"
implementation would use a `type Twig struct { status bool }`
object and a `func (t *Twig) Switch(newStatus bool)` method.

The debug mode has been plugged into the CLI through the `debug` option .

### Unit Tests

I've written some unit tests to check the code I've written. The UT for the load
balancer uses a fake implementation of the service (always replying the ID of
the channel).

### CI

I've copied some gitlab CI from my other projects to ensure the build and the
tests run smoothly.

## TODOs

- The current implementation doesn't clean dead services until it realizes 
  they are dead. This means we could store in the memory of the LoadBalancer 
  a huge list of dead services that we "never" (based on the randomness of 
  Go's map iterations) clean up. A direct communication between the 
  TimeServiceManager and the LoadBalancer could help here.
- The test only allows for asynchronous calls to the LoadBalancer. I haven't 
  added a mutex on the Request call - this could be something we want to do, for several reasons:
    - first, this would help prevent overloading an instance with tons of 
      requests
    - second, we might want to avoid shenanigans of deleting a channel in 
      one of the parallel calls to the LoadBalancer while another one is 
      still sending a message to the service
    - third, not call `delete` on the same key more than once -- Go 
      tolerates this, but we can avoid it.
- The current implementation doesn't handle load at all. It merely routes to 
  an instance. We would need to add a gauge somewhere (either on the load 
  balancer or on the endpoint to achieve this.
  
