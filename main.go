package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type hostentry struct {
	hostname string
	ip       string
}

var hosts []hostentry

const PATH_TO_HOSTS_FILE = "C:/Windows/System32/drivers/etc/hosts"

func main() {
	argsWithoutProg := os.Args[1:]

	// Read the hosts file
	readHostFile()

	// If there are no commandline arguments, show invalid usage message
	if len(argsWithoutProg) == 0 {
		showInvalidUsage()
		return
	}

	if argsWithoutProg[0] == "add" && len(argsWithoutProg) == 3 {
		//TODO Needs admin rights to access the hosts file :/
		hostname := argsWithoutProg[1]
		ip := argsWithoutProg[2]

		addToHostFile(hostentry{hostname, ip})

		fmt.Println("Added hostname to the hosts file.")
	} else if argsWithoutProg[0] == "remove" {
		//TODO Needs to be implemented
	} else if argsWithoutProg[0] == "list" {
		fmt.Printf("Found %s hostnames in hosts file\n", strconv.Itoa(len(hosts)))
	} else if argsWithoutProg[0] == "help" {
		showHelp()
	} else {
		showInvalidUsage()
	}
}

func showInvalidUsage() {
	fmt.Println("Invalid arguments for usage type \"gohost help\"")
}

func showHelp() {
	fmt.Println("gohost help")
	fmt.Println("gohost add <hostname> <ip>")
	fmt.Println("gohost remove <hostname>")
	fmt.Println("gohost list")
}

func addToHostFile(host hostentry) {
	// Create string from hostentry
	var entryString = host.hostname + " " + host.ip + "\n"

	// Open the File with append mode
	f, err := os.OpenFile(PATH_TO_HOSTS_FILE, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	// Close the file at the end of the function
	defer f.Close()

	// Write the string to the file
	if _, err = f.WriteString(entryString); err != nil {
		panic(err)
	}
}

func readHostFile() {
	// open file
	f, err := os.Open(PATH_TO_HOSTS_FILE)
	if err != nil {
		log.Fatal(err)
	}

	// Close the File at the end of the function
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		var text = scanner.Text()
		text = strings.Trim(text, " ")

		// Remove all invisible characters
		text = strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return -1
		}, text)

		// if text starts with # then ignore it
		if len(text) == 0 || text[0] == '#' {
			continue
		}

		// split the line into hostname and ip
		split := strings.Split(text, " ")
		hostname := split[0]
		ip := split[1]

		// Add it to the hosts list
		hosts = append(hosts, hostentry{hostname, ip})

		// Print out the line for debug use
		// fmt.Printf("line: %s\n", text)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
