# GIMPS Go Implementation

GIMPS (Great Internet Mersenne Prime Search) implementation in Go.

## Description

This project is an implementation of Mersenne prime search algorithms in Go. The project includes:

- Lucas-Lehmer test for checking Mersenne number primality
- PRP (Probable Prime) tests
- Factoring algorithms: Trial Factoring, P-1, P+1, ECM
- PrimeNet server integration
- Multi-threaded work unit processing
- State saving and restoration

## Project Structure

```
gimps-go/
├── cmd/gimps/          # Main executable
├── internal/
│   ├── algorithms/     # Testing and factoring algorithms
│   ├── primenet/       # PrimeNet integration
│   ├── worker/         # Worker management
│   ├── storage/         # State persistence
│   └── math/           # Mathematical utilities
└── pkg/config/         # Configuration
```

## Installation

```bash
go build ./cmd/gimps
```

## Usage

```bash
# Start workers
./gimps

# Configuration menu
./gimps -m

# Show status
./gimps -s

# Contact PrimeNet server
./gimps -c

# Torture test
./gimps -t

# Version
./gimps -v
```

## Configuration

Configuration is stored in `prime.ini` file. Main parameters:

- `UsePrimenet` - use PrimeNet (0/1)
- `UserID` - user ID
- `ComputerName` - computer name
- `Workers` - number of workers
- `WorkPreference` - work type preferences
- `DayMemory` / `NightMemory` - memory in MB

## Algorithms

### Lucas-Lehmer Test
Primality test for Mersenne numbers (2^p - 1). Algorithm: s_0 = 4, s_i = (s_{i-1}^2 - 2) mod (2^p - 1).

### PRP Test
Probable Prime test: computes base^(N-1) mod N.

### Trial Factoring
Search for small divisors of Mersenne numbers. Uses the property: divisors have the form 2kp+1.

### P-1 Factoring
Pollard P-1 factoring method with two stages.

### P+1 Factoring
Similar to P-1, but uses P+1 group.

### ECM
Elliptic Curve Method for factoring.

## PrimeNet

The project supports integration with PrimeNet server for receiving assignments and submitting results.

## License

This project is based on the original Prime95 source code from Mersenne Research, Inc.
