package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"system/daemon/progress"
)

var Commands = map[string]string{
	"start":   "start",
	"restart": "restart",
	"stop":    "stop",
	"status":  "status",
}

type Client struct {
	path string
}

func NewClient(path string) *Client {
	return &Client{path}
}

type AskData struct {
	Key string `json:"key,omitempty"`
}

// Manage by daemon commands or run the daemon
func (c *Client) Listen() error {
	if len(os.Args) < 3 {
		return errors.New("The input task is empty. Usage: mydaemon [command] [task name]")
	}
	cmd, ok := Commands[os.Args[1]]
	if !ok {
		return errors.New("Usage: daemon start | restart | stop | status")
	}

	var obj AskData
	for _, info := range os.Args[2:] {
		obj.Key = obj.Key + " " + info
	}
	obj.Key = obj.Key[1:]
	// obj.Key = os.Args[2]
	// if len(os.Args) > 3 {
	// 	args := ""
	// 	for _, info := range os.Args[3:] {
	// 		args = args + " " + info
	// 	}
	// 	obj.Key = obj.Key + "{|}" + args[1:len(args)]
	// }
	log.Println("Key", obj.Key)
	ask, err := json.Marshal(obj)
	if err != nil {
		log.Printf("result to json error:%v", err)
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://localhost%s/server/%s", c.path, cmd), "", bytes.NewBuffer(ask))
	if err != nil {
		// log.Printf("Get symbol init error:%s", err.Error())
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read body error:", err)
		return err
	}
	log.Println("result:", string(data))
	// var ret Result
	// if err := json.Unmarshal(data, &ret); err != nil {
	// 	log.Println("json to struct error:", err, string(data))
	// 	// log.Println("json to struct error:", err)
	// 	return err
	// }

	// log.Println("json to struct error:", ret)

	return nil
}

func TaskTest() {
	go LogUpdate()
	time.Sleep(time.Second * 2)
	for _, key := range []string{"./test/daemon.tasktest", "./test/daemon.1", "./test/daemon.2"} {
		task := progress.NewTask(key, logPath)
		if err := task.Start(); err != nil {
			log.Println("Start:", err)
			return
		}
		Tasks.Set(key, task)
	}
	select {}
}
