package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter domain to check. Press Ctrl+C to exit.")
	fmt.Println("Returns: domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord")

	for scanner.Scan() {
		s := spinner.New(spinner.CharSets[11], 200*time.Millisecond)
		s.Color("green", "bold")
		s.Suffix = "  Checking domain..."
		s.Start()
		time.Sleep(4 * time.Second)

		checkDomain(scanner.Text())
		s.Stop()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln("Error: ", err)
	}
}

func checkDomain(domain string) {
	fmt.Println()
	var hasMx, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecords string

	// Check MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		fmt.Println()
		log.Fatalln("Result: DNS name does not exist! ")
	}
	if len(mxRecords) > 0 {
		hasMx = true
	}

	// Check SPF records
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	// Check DMARC records
	dmarcValues, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	for _, record := range dmarcValues {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecords = record
			break
		}
	}

	// Print the results TODO- COMPREHENSIVE OUTPUT
	fmt.Printf("%v,%v,%v,%v,%v,%v\n", domain, hasMx, hasSPF, spfRecord, hasDMARC, dmarcRecords)
}
