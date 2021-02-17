<div align=center><h1>Parcelivery</h1></div>
<div align=center><strong>Parcel delivery simulation software</strong></div>

## About

**Parcelivery** is a parcel delivery simulation program, that can be used to simulate different real-world delivery strategies and iterate in a timely fashion to meet the needs of the modern international delivery and shipping industry.

## How to Build

Once the project's repository have been cloned, it should be sufficient to run:

```bash
go build
```

to compile the code into a binary called `parcelivery`.

## How to Run

To invoke the program to behave like the subject mandated, simply run:

```bash
./parcelivery INPUT_FILE
```

To invoke the program to see the turns one-by-one and in a visual fashion as an ASCII map, you can run:

```bash
./parcelivery --interactive INPUT_FILE
# or
./parcelivery -i INPUT_FILE
```

Here is the output of the program's help message to learn how to invoke it:

```plain
NAME:
   Parcelivery - parcel delivery simulation software

USAGE:
   parcelivery [global options] INPUT_FILE

GLOBAL OPTIONS:
   --interactive, -i  Enable interactive mode (step-by-step) (default: false)
   --help, -h         Display this help message (default: false)
```

## Project Layout

The project's files are split based on what component of the system it implements, to make it easy to locate any given part of the code.

For instance:

- to find the code that implements the transports' decision making, simply look into `transport.go`
- if one's wondering about how parcels are represented, he'd look no further than `parcel.go`
- when wanting to fix a bug related to how trucks departs, `truck.go` is where to do it
- if the input file was incorrectly parsed, `parse.go` contains all the answers about why
- etc...

We do not have any orthogonal/general-purpose features to expose or factor out, so we don't have any subpackages

## About the Implementation

We use the A* pathfinding algorithm to find our way to parcels and then to the truck.  
We use a nearest-first approach, meaning that transports will prioritize the closest parcels.  
The code is organized in a kind-of object-oriented fashion, in the simple sense that each piece of the simulation has its own struct to represent it and its methods to interact with it.  
