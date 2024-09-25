package main

import (
    "bytes"
    "encoding/gob"
    "fmt"
    "net"
    "sync"
)

// Request structure according to the specified packet layout
type Request struct {
    RequestReply int // 0 for request
    FunctionNum  int // Assuming function number is 1
    Partition     int
    Range1       int
    Range2       int
    N            int // total number of terms
    T            int // exponent
}

// Reply structure according to the specified packet layout
type Reply struct {
    RequestReply   int // 1 for reply
    FunctionNum    int // Assuming function number is 1
    Partition       int
    NumberOfTerms   int
    N                int // total number of terms
    T                int // exponent
    PartitionSum     float64
    IncludeExclude    string // I/X (include/exclude terms)
}

func sendTaskToServer(serverAddr string, task Request, wg *sync.WaitGroup, results chan<- Reply) {
    defer wg.Done()

    conn, err := net.Dial("udp", serverAddr)
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        return
    }
    defer conn.Close()

    // Encode the request
    buf := new(bytes.Buffer)
    encoder := gob.NewEncoder(buf)
    err = encoder.Encode(task)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return
    }

    // Send the request
    _, err = conn.Write(buf.Bytes())
    if err != nil {
        fmt.Println("Error writing to UDP:", err)
        return
    }

    // Prepare to receive reply
    var reply Reply
    dec := gob.NewDecoder(conn)
    err = dec.Decode(&reply)
    if err != nil {
        fmt.Println("Error receiving reply:", err)
        return
    }

    results <- reply
}

func main() {
    var t, n int
    fmt.Print("Enter the exponent (t): ")
    fmt.Scan(&t)

    fmt.Print("Enter the value of n (total terms): ")
    fmt.Scan(&n)

    if n%5 != 0 { // Ensure n is divisible by 5 for 5 servers
        fmt.Println("n must be divisible by 5.")
        return
    }

    // Server addresses to include ports
    servers := []string{
        "localhost:8080",
        "localhost:8081",
        "localhost:8082",
        "localhost:8083",
        "localhost:8084",
    }

    partitionSize := n / 5
    var wg sync.WaitGroup
    results := make(chan Reply, len(servers))

    for p, server := range servers {
        start := p * partitionSize
        end := (p + 1) * partitionSize - 1

        task := Request{
            RequestReply: 0, // Indicate this is a request
            FunctionNum:  1, // Example function number
            Partition:     p,
            Range1:       start,
            Range2:       end,
            N:            n,
            T:            t,
        }

        wg.Add(1)
        go sendTaskToServer(server, task, &wg, results)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    finalSum := 0.0
    totalTerms := 0

    for reply := range results {
        // Print received sum for each partition
        fmt.Printf("Received sum for partition %d: %.2f\n", reply.Partition, reply.PartitionSum)
        
        // Accumulate the partition sum
        finalSum += reply.PartitionSum
        
        // Update total terms
        totalTerms += reply.NumberOfTerms
    }

    // Print the final results
    fmt.Printf("\nTotal sum: %.2f\n", finalSum)
    fmt.Printf("Total number of terms: %d\n", totalTerms)
    fmt.Printf("Exponentiation (t): %d\n", t)
}
