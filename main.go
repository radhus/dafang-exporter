package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	promSetting = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "dafang_setting"}, []string{"setting"})
	promUptime  = promauto.NewGauge(prometheus.GaugeOpts{Name: "dafang_uptime"})
	promSignal  = promauto.NewGauge(prometheus.GaugeOpts{Name: "dafang_signal"})
	promLink    = promauto.NewGauge(prometheus.GaugeOpts{Name: "dafang_link"})
	promNoise   = promauto.NewGauge(prometheus.GaugeOpts{Name: "dafang_noise"})
)

// {"uptime":" 22:04:55 up 34 min,  0 users,  load average: 3.26, 3.25, 2.97",
// "ssid":"LW", "bitrate":"72.2 Mb/s", "signal_level":"100%", "link_quality":"85%", "noise_level":"0%" }
type payload struct {
	Uptime      string
	SSID        string
	Bitrate     string
	SignalLevel string `json:"signal_level"`
	LinkQuality string `json:"link_quality"`
	NoiseLevel  string `json:"noise_level"`
}

type session struct {
	logger *log.Logger
	prefix string
}

func parsePercent(str string) float64 {
	i, err := strconv.ParseFloat(str[:len(str)-1], 64)
	if err != nil {
		return 0
	}
	return i
}

func (s session) status(client mqtt.Client, message mqtt.Message) {
	p := payload{}
	if err := json.Unmarshal(message.Payload(), &p); err != nil {
		s.logger.Println("Error parsing:", err)
	}
	// s.logger.Printf("[%s] %v\n", message.Topic(), p)

	uptime := parseUptime(p.Uptime)
	promUptime.Set(uptime.Seconds())
	promSignal.Set(parsePercent(p.SignalLevel))
	promLink.Set(parsePercent(p.LinkQuality))
	promNoise.Set(parsePercent(p.NoiseLevel))
}

func (s session) settings(client mqtt.Client, message mqtt.Message) {
	trim := len(s.prefix) + 1
	if len(message.Topic()) < trim {
		return
	}

	setting := message.Topic()[trim:]
	s.logger.Printf("[%s] %s=%s\n", message.Topic(), setting, message.Payload())

	var value float64
	switch p := string(message.Payload()); p {
	case "ON":
		value = 1
	case "OFF":
		value = 0
	default:
		i, err := strconv.ParseFloat(p, 64)
		if err != nil {
			s.logger.Printf("[%s] unsupported value: %s", setting, p)
		}
		value = i
	}

	promSetting.WithLabelValues(setting).Set(value)
}

func main() {
	logger := log.New(os.Stdout, "main: ", 0)
	logger.Println("starting...")

	s := session{
		logger: log.New(os.Stdout, "receiver: ", 0),
		prefix: os.Args[2],
	}

	// mqtt.DEBUG = log.New(os.Stdout, "mqtt/debug: ", 0)
	mqtt.ERROR = log.New(os.Stdout, "mqtt/error: ", 0)
	opts := mqtt.NewClientOptions().AddBroker(os.Args[1]).SetClientID("dafang-exporter")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	opts.OnConnect = func(c mqtt.Client) {
		if token := c.Subscribe(s.prefix, 0, s.status); token.Wait() && token.Error() != nil {
			logger.Fatalln(token.Error())
		}
		if token := c.Subscribe(s.prefix+"/#", 0, s.settings); token.Wait() && token.Error() != nil {
			logger.Fatalln(token.Error())
		}

	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logger.Fatalln(token.Error())
	}

	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	// <-sig

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
