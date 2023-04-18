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
	"bufio"
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
	"strconv"
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
	var data []byte
	var err error

	if data, err = os.ReadFile(fileName); err != nil {
		log.Fatal(err)
	}

	parseLine := func(event string) item {
		title := gjson.Get(event, "Topic").String()
		json := gjson.Get(event, "Value").String()
		description := ParseDesc(title, json)

		return item{title: title, description: description, json: json}
	}

	items := []list.Item{}
	// either array or json object on each line
	// -o output is object/line & no array
	if string(data)[0] == '{' {
		//result := gjson.Get(string(data), "..")
		gjson.ForEachLine(string(data), func(line gjson.Result) bool {
			v := parseLine(line.String())
			items = append(items, v)

			return true
		})
	} else {
		result := gjson.Parse(string(data))
		for _, event := range result.Array() {
			v := parseLine(event.String())
			items = append(items, v)
		}
	}

	return items
}

func initGUI(items []list.Item) *tea.Program {
	m := NewModel(items)
	p := tea.NewProgram(m)

	return p
}

// enable/disable forwarding of all events from eventbus to websocket
func SetEventNotifications(api Api, enable bool) (string, error) {
	body, err := api.Get("/notifications")
	if err != nil {
		return "", err
	}

	if enable {
		gotEntry := gjson.Get(body, "#(Conditions.Prefix==\"\")#")
		if len(gotEntry.Array()) > 0 {
			fmt.Println("already have entry for sprbus, will be removed")
			return body, nil
		}

		// Conditions.Prefix="", Notification=false
		settingEntry := ConditionEntry{}
		data, _ := json.Marshal(settingEntry)
		body = string(data)

		res, err := api.Put("/notifications", body)

		return res, err
	} else {
		// find index for sprbus entry - TODO gjson method for this?
		result := gjson.Parse(body)
		keyIndex := int64(-1)
		result.ForEach(func(key, value gjson.Result) bool {
			if gjson.Get(value.String(), "Conditions.Prefix").String() == "" {
				keyIndex, _ = strconv.ParseInt(key.String(), 10, 64)
				return false
			}

			return true // keep iterating
		})

		if keyIndex < 0 {
			return "", errors.New("sprbus event entry not found")
		}

		res, err := api.Delete(fmt.Sprintf("/notifications/%d", keyIndex), "")
		return res, err
	}
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
	fileNameIn := flag.String("f", "", "read sprbus json data from file")
	fileNameOut := flag.String("o", "", "write sprbus json data to file")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	isRemote := *addr != ""

	if *dumpJSON {
		*noGUI = true
	}

	if isRemote {
		token := os.Getenv("TOKEN")
		if token == "" {
			log.Fatal("missing token")
		}

		api := NewApi("192.168.2.1", token)

		// 0. get notification settings
		// 1. set: filter out entries with empty prefix, append {Notification: false}
		// 2. on exit, set: notification settings from 0
		SetEventNotifications(api, true)
		defer SetEventNotifications(api, false)
	} else {
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
		items := []list.Item{}
		if *fileNameIn != "" {
			items = preloadItems(*fileNameIn)
		}

		p = initGUI(items)
	}

	var fileWriter *bufio.Writer

	if *fileNameOut != "" {
		f, err := os.Create(*fileNameOut)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		fileWriter = bufio.NewWriter(f)
	}

	// we get same data from either sprbus or websocket here
	retSub := func(title string, json string) {
		description := ParseDesc(title, json)

		lineJSON := fmt.Sprintf("{\"Topic\": \"%s\", \"Value\": %s}\n", title, json)

		if fileWriter != nil {
			fileWriter.WriteString(lineJSON)
			fileWriter.Flush()
		}

		line := ""

		if *noGUI {
			if *dumpJSON {
				line = lineJSON
			} else {
				//else just print to stdout
				line = fmt.Sprintf("%v %v\n", fgCyan(title), description)
			}

			fmt.Println(line)

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
