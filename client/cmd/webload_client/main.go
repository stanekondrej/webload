package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/uptime"
	"github.com/stanekondrej/webload/client/pkg/pb/github.com/stanekondrej/webload/protobuf"
	"google.golang.org/protobuf/proto"
)

type env struct {
	serverUrl      *url.URL
	updateInterval time.Duration
}

func getEnv() *env {
	serverUrlEnv, ok := os.LookupEnv("WEBLOAD_SERVER")
	if !ok {
		log.Fatal("No server url provided")
	}

	serverUrl, err := url.Parse(serverUrlEnv)
	if err != nil {
		log.Fatal("Failed to parse url")
	}

	updateIntervalEnv, ok := os.LookupEnv("WEBLOAD_UPDATE_INTERVAL")
	updateInterval := time.Second
	if ok {
		i, err := time.ParseDuration(updateIntervalEnv)
		if err != nil {
			log.Fatal("Invalid duration in WEBLOAD_UPDATE_INTERVAL")
		}

		updateInterval = i
	}

	return &env{
		serverUrl,
		updateInterval,
	}
}

func getSystemStats() *protobuf.Stats {
	stats := &protobuf.Stats{}

	mem, err := memory.Get()
	if err == nil {
		stats.MemMax = mem.Total
		stats.MemUsed = mem.Used
		stats.MemUsage = float32(mem.Used / mem.Total)
	}

	cpu, err := cpu.Get()
	if err == nil {
		stats.CpuUsage = float32(cpu.Idle / (cpu.User + cpu.Nice + cpu.System +
			cpu.Idle + cpu.Iowait + cpu.Irq + cpu.Softirq + cpu.Steal + cpu.Guest +
			cpu.GuestNice))
	}

	uptime, err := uptime.Get()
	if err == nil {
		stats.Uptime = uint64(uptime.Abs().Milliseconds())
	}

	return stats
}

func main() {
	env := getEnv()

	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/provide", env.serverUrl.String()), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	t, d, err := conn.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	if t == websocket.TextMessage {
		log.Printf("Session ID: %s", string(d))
		log.Println("(Use this to view usage stats on the frontend)")
	}

	go func() {
		for {
			b, err := proto.Marshal(getSystemStats())
			if err != nil {
				log.Fatal(err)
			}

			err = conn.WriteMessage(websocket.BinaryMessage, b)
			if err != nil {
				log.Println(err)
			}

			time.Sleep(env.updateInterval)
		}
	}()

	for {
		t, d, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}

		if t != websocket.TextMessage {
			log.Printf("Unknown message type %d", t)
			continue
		}

		log.Printf("Incoming message: %s", string(d))
	}
}
