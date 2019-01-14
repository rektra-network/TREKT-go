// Copyright 2018 REKTRA Network, All Rights Reserved.

package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/rektra-network/trekt-go/pkg/trekt"
)

type app struct {
	trekt           *trekt.Trekt
	logSubscription *trekt.LogSubscription
	client          *client
}

func (app *app) close() {
	if app.client != nil {
		app.client.close()
		app.client = nil
	}
	if app.logSubscription != nil {
		app.logSubscription.Close()
		app.logSubscription = nil
	}
	if app.trekt != nil {
		app.trekt.Close()
		app.trekt = nil
	}
}

func (app *app) startLogListening(request string) {
	var err error
	app.logSubscription, err = app.trekt.Log.Subscribe(request)
	if err != nil {
		log.Fatalf(`Failed to subscribe: "%s".`, err)
	}
	go app.logSubscription.Handle(
		func(message trekt.LogMessage) {
			log.Printf("Log event: %s\t%s: %s",
				strings.ToUpper(message.GetLevel()),
				message.GetNodeID(),
				message.GetRecord())
		})
}

func (app *app) startAccessPointReading(host, path string, isUnsecured bool) {
	url := url.URL{Scheme: "wss", Host: host, Path: path}
	if isUnsecured {
		url.Scheme = "ws"
	}
	app.client = createClient(url)
	go func() {
		app.client.run()
		app.client.close()
		app.client = nil
	}()
}
