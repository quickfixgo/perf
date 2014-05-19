package main

import (
	"flag"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/fix/enum"
	"github.com/quickfixgo/quickfix/fix/field"
	"github.com/quickfixgo/quickfix/fix42/newordersingle"
	"log"
	"os"
	"time"
)

var fixconfig = flag.String("fixconfig", "outbound.cfg", "FIX config file")
var sampleSize = flag.Int("samplesize", 1000, "Expected sample size")

var SessionID quickfix.SessionID
var start = make(chan interface{})
var app = &OutboundRig{}

func main() {
	flag.Parse()

	cfg, err := os.Open(*fixconfig)
	if err != nil {
		log.Fatal(err)
	}

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		log.Fatal(err)
	}

	logFactory, err := quickfix.NewFileLogFactory(appSettings)
	if err != nil {
		log.Fatal(err)
	}

	initiator, err := quickfix.NewInitiator(app, appSettings, logFactory)
	if err != nil {
		log.Fatal(err)
	}
	if err = initiator.Start(); err != nil {
		log.Fatal(err)
	}

	<-start

	for i := 0; i < *sampleSize; i++ {
		order := newordersingle.Builder(
			field.NewClOrdID("100"),
			field.NewHandlInst("1"),
			field.NewSymbol("TSLA"),
			field.NewSide(enum.Side_BUY),
			&field.TransactTimeField{},
			field.NewOrdType(enum.OrdType_MARKET))

		quickfix.SendToTarget(order, SessionID)
		time.Sleep(1 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
}

type OutboundRig struct {
	quickfix.SessionID
}

func (e OutboundRig) OnCreate(sessionID quickfix.SessionID) {}
func (e *OutboundRig) OnLogon(sessionID quickfix.SessionID) {
	SessionID = sessionID
	start <- "START"
}
func (e OutboundRig) OnLogout(sessionID quickfix.SessionID)                                    {}
func (e OutboundRig) ToAdmin(msgBuilder quickfix.MessageBuilder, sessionID quickfix.SessionID) {}
func (e OutboundRig) ToApp(msgBuilder quickfix.MessageBuilder, sessionID quickfix.SessionID) (err error) {
	return
}

func (e OutboundRig) FromAdmin(msg quickfix.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	return
}

func (e OutboundRig) FromApp(msg quickfix.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	return
}
