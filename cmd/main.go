package main

// this is a script to view the sprbus

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/spr-networks/sprbus"
	"github.com/tidwall/gjson"
	//"log"
	"os"
	"strings"
	"syscall"
	"time"
)

//var client *sprbus.Client
//var socket = "/tmp/test.sock"
var socket = "/home/spr/super/state/api/eventbus.sock"

func checkAccess(fileName string) error {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Printf("Error accessing file %s: %v\n", fileName, err)
		return err
	}
	mode := fileInfo.Mode()

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return errors.New("failed to get file permissions")
	}

	uid := int(stat.Uid)
	gid := int(stat.Gid)
	euid := os.Geteuid()
	egid := os.Getegid()

	if uid != euid && gid != egid && string(mode) != "Srwxrwxrwx" {
		return errors.New(fmt.Sprintf("bad file permissions on %v (%v)\n"+
			"\tchmod a+rw or run with sudo", fileName, mode))
	}

	return nil
}

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
	dumpJSON := flag.Bool("j", false, "show json output")
	topic := flag.String("t", "", "topic")
	exclude := flag.String("e", "", "exclude")
	publish := flag.String("p", "", "publish, raw json data")
	timeout := flag.Int("timeout", 20, "exit timeout")
	verbose := flag.Bool("v", false, "verbose output")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if err := checkAccess(socket); err != nil {
		ct := color.New(color.FgYellow).SprintFunc()
		fmt.Println(ct("NOTE"), err)
	}

	if *publish != "" && *topic != "" {
		if *verbose {
			fmt.Printf("publish topic: %v, value: %v\n", *topic, *publish)
			//NOTE on valid/invalid json

			if !json.Valid([]byte(*publish)) {
				fmt.Println("error: invalid json:", *publish)
				return
			}
		}

		sprbus.PublishString(*topic, *publish)
		return
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
