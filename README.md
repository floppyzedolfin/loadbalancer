# Load Balancer

## Discussion

### Implementation choices



### Code structure

I've reorganized the code in different packages, each one with its own purpose.
The packages I've worked on are
- loadbalancer: implementation of the task
- timeservice: I had to add a line inside `TimeService.Run`
- twig: see below

## Added utilities

### Twig

During the implementation / debugging of my code, I needed to be able to print
information. As I didn't want to pollute the standard output, I wrote a
tiny `twig` utility (it's a small `log` package). Since I really didn't want to
make it too big, I used a global (local) variable to set its state. A "better "
implementation would use a `type Twig struct { status bool }`
object and a `func (t *Twig) SetTo(newStatus bool)` method.

The debug mode has been plugged into the CLI through the `debug` option .

