package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	portsOpen   int
	portso      []int
	username    string
	concurrency = 10
)

func main() {
	initUsername()

	greeting()

	for {
		fmt.Printf("Hello %s, please enter target host name (or q to quit):", username)
		host := input()

		if host == "q" || host == "Q" {
			break
		}

		fmt.Printf("Enter number of ports to scan:")
		p := input()

		port, _ := strconv.Atoi(p)
		scanPorts(host, port)

		printResults()
		resetResults()
	}
}

func initUsername() {
	currentUser, err := user.Current()
	if err != nil {
		username = "user"
	} else {
		username = currentUser.Username
	}
}

func scanPorts(host string, numPorts int) {
	var wg sync.WaitGroup
	portCh := make(chan int, concurrency)

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go worker(host, portCh, &wg)
	}

	for port := 0; port <= numPorts; port++ {
		portCh <- port
	}
	close(portCh)

	wg.Wait()
}

func worker(host string, portCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range portCh {
		if isPortOpen(host, port) {
			portsOpen++
			portso = append(portso, port)
		}
	}
}

func isPortOpen(host string, port int) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", target, time.Millisecond*100)
	if err != nil {
		fmt.Println(target)
		return false
	}
	defer conn.Close()
	return true
}

func greeting() {
	fmt.Println("\033[34mWelcome to Lennart's Portscanner V1.0")
	fmt.Println("\033[34m-------------------------------------")
	fmt.Println("\033[0m")
}

func input() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}

func printResults() {
	fmt.Printf("\033[31m%d port(s) are open:\n", portsOpen)
	if portsOpen > 0 {
		fmt.Println(portso)
	}
	fmt.Println("\033[0m")
}

func resetResults() {
	portso = nil
	portsOpen = 0
}
