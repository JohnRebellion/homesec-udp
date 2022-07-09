package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	httpUtils "github.com/JohnRebellion/go-utils/http"
)

func main() {
	httpUtils.Client.New(&http.Client{})
	port, err := strconv.ParseInt(getEnv("PORT"), 10, 32)

	if err == nil {
		conn, err := net.ListenUDP("udp", &net.UDPAddr{
			Port: int(port),
			IP:   net.ParseIP("0.0.0.0"),
		})
		if err != nil {
			panic(err)
		}

		defer conn.Close()
		fmt.Printf("server listening %s\n", conn.LocalAddr().String())

		lastData := ""

		for {
			message := make([]byte, 1024)
			rlen, _, err := conn.ReadFromUDP(message[:])

			if err == nil {
				data := strings.TrimSpace(string(message[:rlen]))

				if data != lastData && len(data) > 0 {
					logStash := new(LogStash)
					err = json.Unmarshal([]byte(data), logStash)

					if err == nil {
						httpUtils.RequestJSON(http.MethodPost, fmt.Sprintf("http://%s:%s/api/v1/logStash/stash", getEnv("HOMESEC_HOST"), getEnv("HOMESEC_PORT")), &logStash, nil, http.Header{})
					} else {
						fmt.Println(err)
					}
				}

				lastData = data
			}
		}
	}
}

type LogStash struct {
	Timestamp       int    `json:"timestamp"`
	Sensor          string `json:"sensor"`
	Severity        string `json:"severity"`
	SourceIP        string `json:"src_ip"`
	SourcePort      string `json:"src_port"`
	DestinationIP   string `json:"dst_ip"`
	DestinationPort string `json:"dst_port"`
	Protocol        string `json:"proto"`
	Type            string `json:"type"`
	Trail           string `json:"trail"`
	Info            string `json:"info"`
	Reference       string `json:"reference"`
}

func getEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
