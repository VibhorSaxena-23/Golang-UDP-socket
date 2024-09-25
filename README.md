# Golang UDP Socket - Distributed Math Series Computation

## Problem Overview

The goal of this project is to compute the sum of a math series of the form:

S(n, t) = 0^t + 1^t + 2^t + ... + (n-1)^t


using a client-server architecture in Go. The computation is distributed across 5 servers, and the final sum is computed on the client-side.

- **n**: Total number of terms (100, 200, 300, 400, 500).
- **t**: Exponent (fixed non-negative integer).
  
The client sends requests to 5 servers, each calculating a partition of the sum, and the client aggregates the results to compute the final sum.

## Architecture

### Client
- Sends the values of `t` (exponent) and `p` (partition number) to each server.
- Waits for responses from all 5 servers.
- Aggregates the partial sums received from each server.
- Outputs the final computed result.

### Server
- Receives the partition number and the values of `n` and `t` from the client.
- Computes the sum of the partition range:  

((n/5)p)^t + ((n/5)p + 1)^t + ... + ((n/5)(p+1)-1)^t

- Sends the result back to the client.

## Packet Structure

### Request Packet (Client → Server):
1. Request/Reply flag
2. Function number
3. Partition number (0 to 4)
4. `range1`: Starting range of terms
5. `range2`: Ending range of terms
6. `n`: Total number of terms
7. `t`: Exponentiation value

### Reply Packet (Server → Client):
1. Request/Reply flag
2. Function number
3. Partition number
4. Number of terms in the partition
5. `n`: Total number of terms
6. `t`: Exponentiation value
7. Sum of the partition

## Client Implementation Steps:
1. **Create UDP Socket**: Establish a UDP socket for communication.
2. **Send Requests**:
 - Divide the series into 5 equal partitions.
 - Send partition-specific requests to 5 servers.
3. **Receive Responses**:
 - Wait for each server to return its computed partition sum.
 - Store each result.
4. **Compute Final Sum**: 
 - Aggregate all partition sums to get the final result.
 - Print the final sum, number of terms, and exponent.

## Server Implementation Steps:
1. **Create UDP Socket**: Initialize a UDP socket to listen for incoming requests.
2. **Listen for Requests**: Decode incoming requests to extract `t`, `n`, and partition number.
3. **Compute Partition Sum**:
 - Calculate the sum for the given partition range.
4. **Send Reply**: Send the computed partition sum back to the client.

## Example Flow:

### Client Flow:
- Initialize UDP connection.
- Calculate partition ranges based on `n`.
- Send requests to 5 servers.
- Collect responses and compute the final sum.

### Server Flow:
- Initialize UDP connection.
- Listen for requests.
- Decode the request and compute the partition sum.
- Send the result back to the client.

## Error Handling:
- Handle socket communication errors (timeouts, unavailable servers, etc.).
- Validate input values for `n` and `t`.

## Key Go Libraries:
- `net`: For socket communication using UDP.
- `math`: For performing power calculations.
- `encoding/binary`: For encoding/decoding packets if needed.

## How to Run:
1. Start the server processes on different ports.
2. Run the client to send requests to the servers and receive computed results.

## Future Enhancements:
- Add timeout handling for client-server communication.
- Extend to handle more complex error cases and validation.
