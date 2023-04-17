package main

//⛺︎⛺︎⛺︎⛺︎
// this is a script to view the sprbus

/*
TODO

* pager subview when selecting items
    * json / key/value colorized output

m.Index() (list.Model.Index()) in delegate.go
how to share models? delegates vs. list

*/

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spr-networks/sprbus"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
)

var (
	socket   = "/state/api/eventbus.sock"
	fgYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render
	fgCyan   = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render
)

func checkAccess(fileName string, checkWrite bool) error {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.Fatal(err)
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

	if uid != euid && gid != egid && string(mode.String()[7]) != "r" ||
		(checkWrite && string(mode.String()[8]) != "w") {
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

func ParseDesc(topic string, json string) string {
	gjsons := map[string]string{
		"www:log:access": "@values",
		"devices:save":   "@keys",
	}

	gjsonFilter := ""
	if gjsons[topic] != "" {
		gjsonFilter = gjsons[topic]
	}

	desc := ""
	if topic == "www:log:access" {
		desc = gjson.Get(json, gjsonFilter).String()
	} else if topic == "devices:save" {
		desc = gjson.Get(json, gjsonFilter).String()
	} else if strings.HasPrefix(topic, "device:") {
		desc = gjson.Get(json, "MAC").String()
	} else if strings.HasPrefix(topic, "wifi") {
		desc = gjson.Get(json, "@values").String()
	} else if topic == "dns:block:event" {
		desc = fmt.Sprintf("%v %v", gjson.Get(json, "ClientIP").String(),
			gjson.Get(json, "Name"))
	} else if topic == "dns:serve:event" {
		desc = fmt.Sprintf("%v %v", gjson.Get(json, "Type").String(),
			gjson.Get(json, "Q.0.Name"))
	} else if strings.HasPrefix(topic, "nft") {
		if gjson.Get(json, "DNS.Questions").Exists() {
			q := gjson.Get(json, "DNS.Questions.0.Name")
			name, _ := b64.StdEncoding.DecodeString(q.String())
			desc = fmt.Sprintf("%v lookup %v",
				gjson.Get(json, "IP.SrcIP"), string(name))
		} else if gjson.Get(json, "UDP.SrcPort").Exists() {
			desc = fmt.Sprintf("%v:%v -> %v:%v\n",
				gjson.Get(json, "IP.SrcIP"), gjson.Get(json, "UDP.SrcPort"),
				gjson.Get(json, "IP.DstIP"), gjson.Get(json, "UDP.DstPort"))
		} else {
			desc = fmt.Sprintf("%v:%v -> %v:%v",
				gjson.Get(json, "IP.SrcIP"), gjson.Get(json, "TCP.SrcPort"),
				gjson.Get(json, "IP.DstIP"), gjson.Get(json, "TCP.DstPort"))
		}
	} else {
		desc = gjson.Get(json, "@values").String()
	}

	return desc
}

// TODO make same as filterSub
func handleData(topic string, json string, filter string, exclude string, callback func(string, string)) {
	if exclude != "" && strings.HasPrefix(topic, exclude) {
		return
	}

	if filter != "" && !StartsWithAny(topic, strings.Split(filter, ",")) {
		return
	}

	callback(topic, json)
}

func filterSub(filter string, exclude string, callback func(string, string)) {
	go sprbus.HandleEvent("", func(topic string, json string) {
		handleData(topic, json, filter, exclude, callback)
	})
}

func filterWS(addr string, filter string, exclude string, callback func(string, string)) {
	authString := os.Getenv("TOKEN")
	if authString == "" {
		log.Fatal("missing TOKEN in environment")
	}

	go ConnectWebsocket(addr, authString, func(topic string, json string) {
		handleData(topic, json, filter, exclude, callback)
	})
}

// json representation of sprbus logs
type BusLogItem struct {
	Topic string
	Value map[string]interface{}
}

func preloadItems(fileName string) []list.Item {
	maxNum := 1

	var data []byte
	var err error

	if data, err = os.ReadFile(fileName); err != nil {
		log.Fatal(err)
	}

	items := []list.Item{}

	result := gjson.Parse(string(data))
	for i, event := range result.Array() {
		title := gjson.Get(event.String(), "Topic").String()
		json := gjson.Get(event.String(), "Value").String()
		description := ParseDesc(title, json)

		v := item{title: title, description: description, json: json}
		items = append(items, v)

		if i >= maxNum {
			break
		}
	}

	return items
}

func initGUI() *tea.Program {
	items := []list.Item{} //TODO preloadItems() if -filename

	m := NewModel(items)
	p := tea.NewProgram(m)

	return p
}

func main() {
	help := flag.Bool("help", false, "show help")
	addr := flag.String("addr", "", "http service address, example: 192.168.2.1:80")
	dumpJSON := flag.Bool("j", false, "show json output")
	topic := flag.String("t", "", "topic(s) to filter. Example www:log,wifi:auth")
	exclude := flag.String("e", "nft:lan:in", "exclude topic(s). Example nft:lan:in,nft:wan:out")
	publish := flag.String("p", "", "publish, raw json data")
	timeout := flag.Int("timeout", 20, "exit timeout")
	verbose := flag.Bool("v", false, "verbose output")
	noGUI := flag.Bool("nogui", false, "dont show gui, print to stdout")

	//TODO add these
	//fileNameIn := flag.String("f", "", "read sprbus json data from file")
	//fileNameOut := flag.String("o", "", "write sprbus json data to file")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	isRemote := *addr != ""

	if *dumpJSON {
		*noGUI = true
	}

	if !isRemote {
		if err := checkAccess(socket, *publish != ""); err != nil {
			fmt.Println(fgYellow("NOTE"), err)
		}
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

		if isRemote {
			log.Println("TODO not implemented")
			return
		}

		sprbus.PublishString(*topic, *publish)
		return
	}

	var p *tea.Program

	if !*noGUI {
		p = initGUI()
	}

	// we get same data from either sprbus or websocket here
	retSub := func(title string, json string) {
		description := ParseDesc(title, json)

		//send to view
		if *noGUI {
			if *dumpJSON {
				fmt.Printf("{\"Topic\": \"%s\", \"Value\": %s}\n", title, json)

				return
			}

			//else just print to stdout
			fmt.Printf("%v %v\n", fgCyan(title), description)

			return
		}

		p.Send(EventMsg{title: title, description: description, json: json})

		return
	}

	if isRemote {
		filterWS(*addr, *topic, *exclude, retSub)
	} else {
		filterSub(*topic, *exclude, retSub)
	}

	if !*noGUI {
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		return
	}

	if *timeout != 0 {
		time.Sleep(time.Second * time.Duration(*timeout))
	} else {
		for {
			time.Sleep(time.Second)
		}
	}
}
