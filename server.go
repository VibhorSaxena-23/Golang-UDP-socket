package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

// Request structure as defined in client
type Request struct {
	RequestReply int // 0 for request
	FunctionNum  int // Function number
	Partition    int
	Range1       int
	Range2       int
	N            int // total number of terms
	T            int // exponent
}

// Reply structure as defined in client
type Reply struct {
	RequestReply   int // 1 for reply
	FunctionNum    int // Function number
	Partition      int
	NumberOfTerms  int
	N              int // total number of terms
	T              int // exponent
	PartitionSum   float64
	IncludeExclude string // I/X (include/exclude terms)
}

// Function to compute the sum of the partition
func computePartitionSum(range1, range2 int) float64 {
	sum := 0.0
	for i := range1; i <= range2; i++ {
		sum += float64(i) * float64(i)
	}
	return sum
}

func main() {
	// Set up the UDP address
	addr := net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("localhost"),
	}

	// Start listening for UDP connections
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	// Improved logging for server address

	fmt.Printf("Server listening on %s:%d\n", addr.IP.String(), addr.Port)

	for {
		buf := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		// Decode the request
		var request Request
		decoder := gob.NewDecoder(bytes.NewReader(buf[:n]))
		err = decoder.Decode(&request)
		if err != nil {
			fmt.Println("Error decoding request:", err)
			continue
		}

		// Compute the sum for the requested partition
		partitionSum := computePartitionSum(request.Range1, request.Range2)

		// Prepare the reply
		reply := Reply{
			RequestReply:   1, // Indicate this is a reply
			FunctionNum:    request.FunctionNum,
			Partition:      request.Partition,
			NumberOfTerms:  request.Range2 - request.Range1 + 1, // Total terms in this partition
			N:              request.N,
			T:              request.T,
			PartitionSum:   partitionSum,
			IncludeExclude: "Include", // Can change based on your logic
		}

		// Encode the reply
		bufReply := new(bytes.Buffer)
		encoder := gob.NewEncoder(bufReply)
		err = encoder.Encode(reply)
		if err != nil {
			fmt.Println("Error encoding reply:", err)
			continue
		}

		// Send the reply back to the client
		_, err = conn.WriteToUDP(bufReply.Bytes(), remoteAddr)
		if err != nil {
			fmt.Println("Error sending reply:", err)
			continue
		}

		// Print the computed partition sum for logging
		fmt.Printf("Processed partition %d: sum = %.2f\n", request.Partition, partitionSum)

	}
}
