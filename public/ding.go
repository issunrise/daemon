package public

import (
	"fmt"
	"log"
	"time"

	"github.com/hugozhu/godingtalk"
)

type AlertChan chan string

var AlertSignal AlertChan
var ding *godingtalk.DingTalkClient

func NewDing(name, corpid, corpsecret, robot string) {
	ding = godingtalk.NewDingTalkClient(corpid, corpsecret)
	AlertSignal = make(AlertChan, 100)
	alert := new(DingSend)
	alert.Robot = robot
	heart := time.NewTicker(time.Hour * 2)
	go func() {
		for {
			select {
			case msg, ok := <-AlertSignal:
				if ok {
					alert.Msg = fmt.Sprintf("%s:\r%s \r%v", name, msg, time.Now())
					alert.Send()
				}
			case <-heart.C:
				alert.Msg = fmt.Sprintf("%s daemon is normal. \r%v", name, time.Now())
				alert.Send()
			}
		}
	}()
}

type DingSend struct {
	Msg   string `json:"msg"`
	Robot string `json:"robot"`
}

func (c *DingSend) Send() {
	cnt := 0
	for cnt < 2 {
		ding.RefreshAccessToken()
		if err := ding.SendRobotTextMessage(c.Robot, c.Msg); err != nil {
			log.Println("send dingtalk error:", err)
			time.Sleep(time.Millisecond * 500)
			cnt++
		} else {
			return
		}
	}
}
