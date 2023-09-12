package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func udpEchoServer(ip string, port int) {
	uniqueIPs := make(map[string]struct{})

	// Load initial unique IPs from file
	file, err := os.Open("unique_ips.log")
	if err == nil {
		defer file.Close()
		var ip string
		for {
			_, err := fmt.Fscanf(file, "%s\n", &ip)
			if err != nil {
				break
			}
			uniqueIPs[ip] = struct{}{}
		}
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("Failed to resolve address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Listening on %s for UDP packets...\n", addr)

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		clientData := strings.ReplaceAll(string(buf[:n]), "\n", "")
		logToFile(
			"udp_all_traffic.log",
			fmt.Sprintf("Received from %s: %s\n", clientAddr, clientData),
		)
		// Check for unique IP and log if needed
		if _, exists := uniqueIPs[clientAddr.IP.String()]; !exists {
			uniqueIPs[clientAddr.IP.String()] = struct{}{}
			logToFile("unique_ips.log", clientAddr.IP.String()+"\n")
		}

		// Echo the data back to the client
		conn.WriteToUDP(buf[:n], clientAddr)
	}
}

func logToFile(filename, data string) {
	fullPath := "/var/log/" + filename

	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func main() {
	// Set your IP and port here
	IP := "0.0.0.0"
	PORT := 5060
	udpEchoServer(IP, PORT)
}
