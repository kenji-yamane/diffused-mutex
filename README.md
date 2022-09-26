# diffused-mutes

Implementation of the Ricart-Agrawala algorithm on a simulation of a distributed system with ipc

## Instructions

To build all targets:
```bash
make build-all
```

Then, to execute a simulation of three different processes
attempting to access a shared resource, first instantiate
the processes simulators:
```bash
make p1
make p2
make p3
```
On three different terminals.

Then, instantiate the shared resource simulator:
```bash
make cs
```

In case you want to build your own simulation with a different
number of processes, you can do so executing the binary directly.
For instance, to simulate four different processes:
```bash
./bin/process 1 10005 10004 10003 10002
./bin/process 2 10005 10004 10003 10002
./bin/process 3 10005 10004 10003 10002
./bin/process 4 10005 10004 10003 10002
```
The first argument indicates to the simulator which is his own port,
from the arguments that follows.
