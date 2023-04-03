package main

// this is a script to view the sprbus

import (
	b64 "encoding/base64"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/spr-networks/sprbus"
	"github.com/tidwall/gjson"
	"os"
	"strings"
	"time"
)

//var client *sprbus.Client
//var socket = "/tmp/test.sock"
var socket = "/home/spr/super/state/api/eventbus.sock"

func StartsWithAny(t string, list []string) bool {
	for _, m := range list {
		if strings.HasPrefix(t, m) {
			return true
		}
	}

	return false
}

func sub(filter string, exclude string, dumpJSON bool) {
	go sprbus.HandleEvent("", func(topic string, json string) {
		ct := color.New(color.FgCyan).SprintFunc()

		if exclude != "" && strings.HasPrefix(topic, exclude) {
			return
		}

		if filter != "" && !StartsWithAny(topic, strings.Split(filter, ",")) {
			return
		}

		if dumpJSON {
			fmt.Println(json)
			return
		}

		gjsons := map[string]string{
			"www:log:access": "@values",
			"devices:save":   "@keys",
		}

		gjsonFilter := ""
		if gjsons[topic] != "" {
			gjsonFilter = gjsons[topic]
		}

		if topic == "www:log:access" {
			fmt.Printf("%12v %v\n", ct(topic), gjson.Get(json, gjsonFilter))
		} else if topic == "devices:save" {
			fmt.Printf("%12s %v\n", ct(topic), gjson.Get(json, gjsonFilter))
		} else if strings.HasPrefix(topic, "device:") {
			fmt.Printf("%12s %v\n", ct(topic), gjson.Get(json, "MAC"))
		} else if strings.HasPrefix(topic, "wifi") {
			fmt.Printf("%12v %v\n", ct(topic), gjson.Get(json, "@values"))
		} else if strings.HasPrefix(topic, "nft") {
			if gjson.Get(json, "DNS.Questions").Exists() {
				q := gjson.Get(json, "DNS.Questions.0.Name")
				name, _ := b64.StdEncoding.DecodeString(q.String())
				fmt.Printf("%v %v lookup %v\n", ct(topic),
					gjson.Get(json, "IP.SrcIP"), string(name))
			} else if gjson.Get(json, "UDP.SrcPort").Exists() {
				fmt.Printf("%v %v:%v -> %v:%v\n", ct(topic),
					gjson.Get(json, "IP.SrcIP"), gjson.Get(json, "UDP.SrcPort"),
					gjson.Get(json, "IP.DstIP"), gjson.Get(json, "UDP.DstPort"))
			} else {
				fmt.Printf("%v %v:%v -> %v:%v\n", ct(topic),
					gjson.Get(json, "IP.SrcIP"), gjson.Get(json, "TCP.SrcPort"),
					gjson.Get(json, "IP.DstIP"), gjson.Get(json, "TCP.DstPort"))
			}
		} else {
			fmt.Printf("%12v %v\n", ct(topic), json)
		}
	})
}

func main() {
	help := flag.Bool("help", false, "show help")
	dumpJSON := flag.Bool("j", false, "json")
	topic := flag.String("t", "", "topic")
	exclude := flag.String("e", "", "exclude")
	publish := flag.String("p", "", "publish, raw json data")
	timeout := flag.Int("timeout", 20, "exit timeout")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *publish != "" && *topic != "" {
		sprbus.Publish(*topic, *publish)
	} else {
		sub(*topic, *exclude, *dumpJSON)
	}

	if *timeout != 0 {
		time.Sleep(time.Second * time.Duration(*timeout))
	} else {
		for {
			time.Sleep(time.Second)
		}
	}
}
