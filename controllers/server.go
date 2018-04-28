package controllers

import (
	"fmt"
	"log"
	"time"

	"system/daemon/progress"
	"system/daemon/public"
)

// Create your own business code
var (
	Tasks   = progress.NewTaskIndex()
	logPath = newLogPath()
)

type Server struct {
	*Controller
}

func newLogPath() string {
	return fmt.Sprintf("log/%s", time.Now().Format("20060102"))
}

func LogUpdate() {
	log.Println("LogUpdate")
	if err := public.NewDir(logPath); err != nil {
		log.Println("create log path error:", err)
		return
	}
	logdate := time.NewTicker(time.Hour)
	// logdate := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-logdate.C:
			logPath = newLogPath()
			if err := public.NewDir(logPath); err != nil {
				continue
			}
			for _, task := range Tasks.List() {
				task.NewFile(logPath)
			}
		}
	}
}

/*************************************************************/
func (c *Server) Start() {
	var ask AskData
	c.RequestStruct(&ask)

	if task := Tasks.Get(ask.Key); task != nil {
		status, err := task.Status()
		if err != nil {
			c.ResultJson(105, err.Error())
			return
		}
		c.ResultJson(106, fmt.Sprintf("Task had start, %v", status))
		return
	}

	task := progress.NewTask(ask.Key, logPath)
	if err := task.Start(); err != nil {
		c.ResultJson(102, err.Error())
		return
	}

	Tasks.Set(ask.Key, task)
	c.ResultJson(0, "success")
}

func (c *Server) Restart() {
	var ask AskData
	c.RequestStruct(&ask)

	if task := Tasks.Get(ask.Key); task != nil {
		if err := task.Stop(); err != nil {
			c.ResultJson(104, err.Error())
			// return
		}
	}

	task := progress.NewTask(ask.Key, logPath)
	if err := task.Start(); err != nil {
		c.ResultJson(102, err.Error())
		return
	}
	Tasks.Set(ask.Key, task)
	c.ResultJson(0, "success")
}

func (c *Server) Stop() {
	var ask AskData
	c.RequestStruct(&ask)
	task := Tasks.Get(ask.Key)
	if task == nil {
		c.ResultJson(107, fmt.Sprintf("Can not find %s task", ask.Key))
		return
	}

	if err := task.Stop(); err != nil {
		Tasks.Del(ask.Key)
		c.ResultJson(104, err.Error())
		return
	}
	Tasks.Del(ask.Key)
	c.ResultJson(0, "success")
}

func (c *Server) Status() {
	var ask AskData
	c.RequestStruct(&ask)
	task := Tasks.Get(ask.Key)
	if task == nil {
		c.ResultJson(107, fmt.Sprintf("Can not find %s task", ask.Key))
		return
	}

	status, err := task.Status()
	if err != nil {
		c.ResultJson(105, err.Error())
		return
	}
	c.ResultJson(0, status)
}
