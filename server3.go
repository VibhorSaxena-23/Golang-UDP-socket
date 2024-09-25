package main

import (
    "bytes"
    "encoding/gob"
    "fmt"
    "math"
    "net"
)

type Request struct {
    Partition int
    Range1    int
    Range2    int
    T         int
}

type Reply struct {
    Partition    int
    PartitionSum float64
}

func computePartitionSum(rangeStart, rangeEnd, t int) float64 {
    sum := 0.0
    for i := rangeStart; i <= rangeEnd; i++ {
        sum += math.Pow(float64(i), float64(t))
    }
    return sum
}

func main() {
    addr, err := net.ResolveUDPAddr("udp", ":8082")
    if err != nil {
        fmt.Println("Error resolving address:", err)
        return
    }
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Println("Error listening on UDP:", err)
        return
    }
    defer conn.Close()

    fmt.Println("UDP server is running on port 8082...")

    for {
        buffer := make([]byte, 2048) // Buffer size
        n, clientAddr, err := conn.ReadFromUDP(buffer) // Capture client address
        if err != nil {
            fmt.Println("Error reading from UDP:", err)
            continue
        }

        var req Request
        buf := bytes.NewReader(buffer[:n])
        decoder := gob.NewDecoder(buf)

        // Decode the request
        err = decoder.Decode(&req)
        if err != nil {
            fmt.Println("Error decoding request:", err)
            fmt.Printf("Received raw data: %v\n", buffer[:n])
            continue
        }

        // Compute the partition sum
        partitionSum := computePartitionSum(req.Range1, req.Range2, req.T)

        // Prepare the reply
        reply := Reply{
            Partition:    req.Partition,
            PartitionSum: partitionSum,
        }

        // Send the reply back to the client
        replyBuf := new(bytes.Buffer)
        encoder := gob.NewEncoder(replyBuf)
        err = encoder.Encode(reply)
        if err != nil {
            fmt.Println("Error encoding reply:", err)
            continue
        }

        // Send the reply to the correct client address
        _, err = conn.WriteToUDP(replyBuf.Bytes(), clientAddr) // Use clientAddr here
        if err != nil {
            fmt.Println("Error sending reply:", err)
        }
    }
}
