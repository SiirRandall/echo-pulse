package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var logDir = "logs"

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
			logDir,
			"udp_all_traffic.log",
			fmt.Sprintf("Received from %s: %s\n", clientAddr, clientData),
		)
		// Check for unique IP and log if needed
		if _, exists := uniqueIPs[clientAddr.IP.String()]; !exists {
			uniqueIPs[clientAddr.IP.String()] = struct{}{}
			logToFile(logDir, "unique_ips.log", clientAddr.IP.String()+"\n")
		}

		// Echo the data back to the client
		conn.WriteToUDP(buf[:n], clientAddr)
	}
}

func logToFile(directory, filename, data string) {
	// Construct the full path for the file
	fullPath := filepath.Join(directory, filename)

	// Ensure the directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	// Check if the file exists and create if not
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		_, err := os.Create(fullPath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		fmt.Printf("File %s created.\n", fullPath)
	}

	// Write the data to the file
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, 0644)
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
