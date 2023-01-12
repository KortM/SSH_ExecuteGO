package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/crypto/ssh"
)

// Structure connection and command execution
type Executor struct {
	config  *ssh.ClientConfig //Config SSH
	session *ssh.Session      //SSH Session
	stdin   io.WriteCloser    //Standard input thread
	stdout  io.Reader         //Standart output thread
	wr      chan []byte       //Channel for write to session
}

// Execute SSH command
func (e *Executor) ExecuteCommand(command string) error {
	_, err := e.stdin.Write([]byte(command + "\n"))
	if err != nil {
		return err
	}
	return nil
}

// Printing result
func (e *Executor) PrintResult() {
	scanner := bufio.NewScanner(e.stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

}

// Connecting to Host
func (e *Executor) ConnectToHost(host string) error {
	conn, err := ssh.Dial("tcp", host, e.config)
	if err != nil {
		return err
	}
	e.session, err = conn.NewSession()
	if err != nil {
		return err
	}
	e.stdin, err = e.session.StdinPipe()
	if err != nil {
		return err
	}
	e.stdout, err = e.session.StdoutPipe()
	if err != nil {
		return err
	}
	if err := e.session.Shell(); err != nil {
		return err
	}
	return nil
}

// Pars host html file
func ParseHost(pathToFile string) []string {
	f, err := os.Open(pathToFile)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	value := make([]string, 1)
	doc.Find("tbody").Each(func(i int, s *goquery.Selection) {
		td := s.Find("td")
		//fmt.Println(td.Nodes[3].FirstChild.Data)
		value = append(value, td.Nodes[3].FirstChild.Data)
	})
	return value
}

// For goroutines
func run(host string, c chan bool) {
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("Mar02031812"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	enablePasswd := "Mar02031812"
	executor := new(Executor)
	executor.config = config
	defer func() {
		if r := recover(); r != nil {
			c <- false
		}
	}()
	err := executor.ConnectToHost(host + ":22")
	if err != nil {
		fmt.Printf("Host: %s, Error: %s \n", host, err)
		c <- false
	}
	err = executor.ExecuteCommand("cs_console")
	if err != nil {
		c <- false
	}
	err = executor.ExecuteCommand("en")
	if err != nil {
		c <- false
	}
	err = executor.ExecuteCommand(enablePasswd)
	if err != nil {
		c <- false

	}
	err = executor.ExecuteCommand("sh run")
	if err != nil {
		c <- false
	}
	err = executor.ExecuteCommand("conf t")
	if err != nil {
		c <- false
	}
	err = executor.ExecuteCommand("no logging 192.168.2.2")
	if err != nil {
		c <- false
	}
	err = executor.ExecuteCommand("end")
	if err != nil {
		c <- false
	}
	executor.PrintResult()
	c <- true
}
func main() {
	//hosts := ParseHost("hosts.html")
	hosts := []string{"176.32.0.1", "176.32.0.2"}
	done := make(chan bool, 10)

	for _, host := range hosts {
		if len(host) != 0 {
			go run(host, done)
		}
	}
	for i := 0; i < len(hosts); i++ {
		fmt.Println(<-done)
	}
	/*conn, err := ssh.Dial("tcp", "176.32.0.1:22", config)
	if err != nil {
		log.Fatalf("unable to connect: %s", err)
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Не получилось запустить сессию %s", err)
	}
	defer session.Close()
	//var buf bytes.Buffer
	//session.Stdout = &buf

	stdin, err := session.StdinPipe()
	if err != nil {
		fmt.Println("Unable to setup STDIN")
		fmt.Println("Error : ", err.Error())
	}

	if err := session.Shell(); err != nil {
		fmt.Print(err)
	}
	wr := make(chan []byte, 20)

	go func() {
		for {
			select {
			case d := <-wr:
				_, err := stdin.Write(d)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}()
	//session.Shell()
	for {
		fmt.Println("$")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		fmt.Println(text)
		wr <- []byte(text + "\n")
	} */
}
