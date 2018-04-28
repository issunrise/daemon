# This is a daemon in linux.

### For exmple:


### Run daemon in linux
```
go build 
nohup ./daemon &
```

### Run task in daemon. 

usage: mydaemon [command] [task name]

command: start | restart | stop | status
```
./daemon start taskname
```

### For exmple to use:
go build 
nohup ./daemon &
./daemon start ./test/tasktest

### route:

```
func init() {
	RouteMap = map[string]func(http.ResponseWriter, *http.Request){
		"/server/": ServerRoute,
	}
}

func ServerRoute(w http.ResponseWriter, r *http.Request) {
	log.Println("r.URL.Path", r.URL.Path)
	if parts := strings.Split(r.URL.Path, "/"); len(parts) >= 3 {
		if method := reflect.ValueOf(&conn.Server{conn.NewController(w, r, cfg)}).MethodByName(strings.Title(parts[2])); method.IsValid() {
			method.Call(nil)
		}
		return
	}
	conn.NewController(w, r, cfg).Error()
}

```

### controller:

```
var Tasks = progress.NewTasks()

type Server struct {
	*Controller
}

func (c *Server) Start() {
	var ask datas.AskData
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

	task := new(progress.Task)
	if err := task.Init(ask.Key); err != nil {
		c.ResultJson(101, err.Error())
		return
	}
	if err := task.Start(); err != nil {
		c.ResultJson(102, err.Error())
		return
	}

	Tasks.Set(ask.Key, task)
	c.ResultJson(0, "success")
}

func (c *Server) Restart() {
	var ask datas.AskData
	c.RequestStruct(&ask)

	if task := Tasks.Get(ask.Key); task != nil {
		if err := task.Stop(); err != nil {
			c.ResultJson(104, err.Error())
			return
		}
	}

	task := new(progress.Task)
	if err := task.Init(ask.Key); err != nil {
		c.ResultJson(101, err.Error())
		return
	}
	if err := task.Start(); err != nil {
		c.ResultJson(102, err.Error())
		return
	}
	Tasks.Set(ask.Key, task)
	c.ResultJson(0, "success")
}

func (c *Server) Stop() {
	var ask datas.AskData
	c.RequestStruct(&ask)
	task := Tasks.Get(ask.Key)
	if task == nil {
		c.ResultJson(107, fmt.Sprintf("Can not find %s task", ask.Key))
		return
	}

	if err := task.Stop(); err != nil {
		c.ResultJson(104, err.Error())
		return
	}
	Tasks.Set(ask.Key, nil)
	c.ResultJson(0, "success")
}

func (c *Server) Status() {
	var ask datas.AskData
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

```

### log:

创建日志路径周期 {每天一个日期路径}
```
func newLogPath() string {
	return fmt.Sprintf("log/%s", time.Now().Format("20060102"))
}

```

日志更新
```
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
			//判断路径，若已存:跳过创建新路径
			if err := public.NewDir(logPath); err != nil {
				continue
			}
			for _, task := range Tasks.List() {
				task.NewFile(logPath)
			}
		}
	}
}

```