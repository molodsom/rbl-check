package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Could not open the %s file", path))
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln(fmt.Sprintf("Could not read the %s file", path))
		}
	}(file)
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		log.Fatalln(scanner.Err())
	}
	return lines
}

func main() {
	hosts := os.Args[1:]
	if len(hosts) == 0 {
		hosts = readLines("hosts.txt")
	}
	srvs := readLines("dnsbl.txt")
	for _, h := range hosts {
		ips, err := net.LookupIP(h)
		if err == nil {
			ip := ips[0].String()
			log.Println(fmt.Sprintf("Checking %s [%s]:", h, ip))
			s := strings.Split(ip, ".")
			ipr := fmt.Sprintf("%s.%s.%s.%s", s[3], s[2], s[1], s[0])
			for _, srv := range srvs {
				rec, err := net.LookupTXT(fmt.Sprintf("%s.%s", ipr, srv))
				if err == nil {
					log.Println(fmt.Sprintf("Response %s [%s] (%s): %s", h, ip, srv, rec))
				}
			}
		}
	}
	log.Println("Completed")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}
