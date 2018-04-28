package progress

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"system/daemon/public"
	"time"
)

type Task struct {
	dir         string
	args        []string
	name        string
	pid         int
	done        bool
	file        *FileClient
	cmd         *exec.Cmd
	alertSignal public.AlertChan
	errCnt      int
}

func NewTask(key, path string) *Task {

	args := strings.Split(key, " ")
	dirs := strings.Split(args[0], "/")
	name := dirs[len(dirs)-1]
	dir := strings.Replace(args[0], name, "", -1)
	args[0] = fmt.Sprintf("./%s", name)
	log.Printf("Start----------------- Init:%v, path:%v, c.args:%v", name, dir, args)
	file := NewFileClient(fmt.Sprintf("%s/%s.log", path, name))

	return &Task{
		dir:         dir,
		args:        args,
		name:        name,
		done:        false,
		alertSignal: public.AlertSignal,
		file:        file,
	}
}

func (c *Task) NewFile(path string) error {
	return c.file.Update(fmt.Sprintf("%s/%s.log", path, c.name))
}

// Start the service
func (c *Task) Start() error {
	go func() {
		defer c.file.Close()
		for {
			if c.done {
				break
			}
			// c.cmd = exec.Command(fmt.Sprintf("./%s.exe", c.name))
			// c.cmd = exec.Command("nohup", fmt.Sprintf("./%s", c.name))
			c.cmd = exec.Command("nohup")
			c.cmd.Dir = c.dir
			for _, v := range c.args {
				c.cmd.Args = append(c.cmd.Args, v)
			}
			c.cmd.Stdout = c.file
			c.cmd.Stderr = c.file
			if err := c.cmd.Start(); err != nil {
				c.alert(fmt.Sprintf("%s start error: %v", c.name, err))
				time.Sleep(time.Second * 2)
				continue
			}
			c.pid = c.cmd.Process.Pid

			c.alert(fmt.Sprintf("%s start success. pid：%v", c.name, c.pid))
			if err := c.cmd.Wait(); err != nil {
				c.alert(fmt.Sprintf("%s exit error：%v", c.name, err))
				time.Sleep(time.Second * 2)
			}
			log.Println("next:", c.name, c.pid)
		}
	}()
	return nil
}

func (c *Task) alert(msg string) {
	log.Println(msg)
	if c.errCnt > 5 {
		return
	}
	c.alertSignal <- msg
	c.errCnt++
}

// Stop the service
func (c *Task) Stop() error {

	c.done = true
	if err := c.cmd.Process.Kill(); err != nil {
		return err
	}

	// process, err := os.FindProcess(c.Pid)
	// if err != nil {
	// 	log.Println("FindProcess:", err)
	// 	return err
	// }
	// log.Println("do---1")
	// if err := process.Kill(); err != nil {
	// 	return err
	// }
	return nil
}

// Status - Get service status
func (c *Task) Status() (interface{}, error) {
	// process, err := os.FindProcess(c.Pid)
	// if err != nil {
	// 	log.Println("FindProcess:", err)
	// 	return nil, err
	// }
	// log.Println("do---1", process)

	// log.Println("doSysUsage1", c.Cmd.ProcessState.SysUsage())

	// log.Printf("%s status:%s", c.nowtask.Name, processState.String())
	status := fmt.Sprintf("%s status: %v", c.name, *c.cmd.Process)
	return status, nil
}
