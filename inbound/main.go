package main

import (
	"flag"
	"github.com/grd/stat"
	"github.com/quickfixgo/quickfix"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var fixconfig = flag.String("fixconfig", "inbound.cfg", "FIX config file")
var sampleSize = flag.Int("samplesize", 1000, "Expected sample size")

var count = 0
var allDone = make(chan interface{})
var app = &InboundRig{}
var metrics stat.IntSlice

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		defer pprof.StopCPUProfile()
	}

	log.Print("NumCPU: ", runtime.NumCPU())
	log.Print("GOMAXPROCS: ", runtime.GOMAXPROCS(-1))

	metrics = make(stat.IntSlice, *sampleSize)

	cfg, err := os.Open(*fixconfig)
	if err != nil {
		log.Fatal(err)
	}

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		log.Fatal(err)
	}

	logFactory := quickfix.NewNullLogFactory()

	acceptor, err := quickfix.NewAcceptor(app, appSettings, logFactory)
	if err != nil {
		log.Fatal(err)
	}
	if err = acceptor.Start(); err != nil {
		log.Fatal(err)
	}

	<-allDone

	metricsUS := make(stat.Float64Slice, *sampleSize)
	for i, durationNS := range metrics {
		metricsUS[i] = float64(durationNS) / 1000.0
	}

	mean := stat.Mean(metricsUS)
	max, maxIndex := stat.Max(metricsUS)
	stdev := stat.Sd(metricsUS)

	log.Printf("Sample mean is %g us", mean)
	log.Printf("Sample max is %g us (%v)", max, maxIndex)
	log.Printf("Standard Dev is %g us", stdev)
}

type InboundRig struct{}

func (e InboundRig) OnCreate(sessionID quickfix.SessionID) {}
func (e InboundRig) OnLogon(sessionID quickfix.SessionID) {
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
	}
}
func (e InboundRig) OnLogout(sessionID quickfix.SessionID)                                    {}
func (e InboundRig) ToAdmin(msgBuilder quickfix.MessageBuilder, sessionID quickfix.SessionID) {}
func (e InboundRig) ToApp(msgBuilder quickfix.MessageBuilder, sessionID quickfix.SessionID) (err error) {
	return
}

func (e InboundRig) FromAdmin(msg quickfix.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	return
}

func (e *InboundRig) FromApp(msg quickfix.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	metrics[count] = int64(time.Since(msg.ReceiveTime))
	count++

	if count == *sampleSize {
		allDone <- "DONE"
	}
	return
}
