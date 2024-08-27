package main

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tatsushid/go-fastping"
)

var devices = make(map[string]int)
var pinger = fastping.NewPinger()

func main() {
	pinger.MaxRTT = time.Millisecond * time.Duration(config.PingIntervalMs)
	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		remoteIP := addr.String()
		//log.Printf("IP Addr: %s receive, RTT: %v\n", remoteIP, rtt)
		_, exists := devices[remoteIP]
		if exists {
			devices[remoteIP] = 0
		}
	}
	pinger.OnIdle = func() {
		for remoteIP, fails := range devices {
			if fails > 0 {
				log.Println(remoteIP, fails)
			}
			if fails > config.MaxFails {
				log.Println("Device went offline", remoteIP)
				removeDevice(remoteIP)
				reportSignal(remoteIP, "off")
			} else {
				devices[remoteIP]++
			}
		}
	}
	pinger.AddIP("127.0.0.1")
	pinger.RunLoop()
	go func() {
		time.Sleep(time.Millisecond * time.Duration(config.PingIntervalMs-10))
		pinger.RemoveIP("127.0.0.1")
		<-pinger.Done()
		if err := pinger.Err(); err != nil {
			log.Fatalln("Ping failed:", err)
		}
	}()
	startHttp()
}

func addDevice(remoteIP string) {
	pinger.AddIP(remoteIP)
	devices[remoteIP] = 0
}

func removeDevice(remoteIP string) {
	pinger.RemoveIP(remoteIP)
	delete(devices, remoteIP)
}

func buildURL(path ...string) string {
	return config.ReportUrl + "/" + strings.Join(path, "/")
}

func reportSignal(signal string, power string) {
	resp, err := http.Post(buildURL(signal, power), "", nil)
	if err != nil {
		log.Println("Failed to report signal", err)
	} else {
		defer resp.Body.Close()
		log.Println("Response code:", resp.StatusCode)
	}
}

func startHttp() {
	app := gin.Default()
	app.POST("/", func(ctx *gin.Context) {
		remoteIP := ctx.RemoteIP()
		log.Println("Device went online", remoteIP)
		addDevice(remoteIP)
		reportSignal(remoteIP, "on")
		ctx.Status(201)
	})
	app.Run(":" + strconv.Itoa(config.Port))
}
