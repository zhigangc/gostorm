package gostorm

import "os"
import "encoding/json"
import "syscall"
import "fmt"
import "io/ioutil"
import "io"
import "bufio"

type Component struct {
	Conf Values
	Context Values
	PendingCommands []Values
	PendingTaskids []TaskIds
	debug *bufio.Writer
    Type string
}

func (comp *Component) Debug(info string) {
	comp.debug.WriteString(info + "\n")
	comp.debug.Flush()
}

func (comp *Component) Init(who string) {
	tmpfile, err := os.OpenFile("/tmp/go.out."+who, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	comp.debug = bufio.NewWriter(tmpfile)

	setupInfo := comp.ReadValues(os.Stdin)
    pidDir, _ := setupInfo.GetString("pidDir")
    comp.SendPid(pidDir)
    conf, _ := setupInfo.GetValues("conf")
    context, _ := setupInfo.GetValues("context")
    comp.Conf = conf
    comp.Context = context
}

func (comp *Component) Run() {}
func (comp *Component) Process(tup *Tuple) {}

func (comp *Component) ReadTuple() *Tuple {
    vals := comp.ReadCommand()
    return NewTuple(vals)
}

func (comp *Component) ReadMsg(r io.Reader) []byte {
    var msgBuf []byte
    bio := bufio.NewReader(r)
    for {
        var line []byte
        data, hasMore, err := bio.ReadLine()
        for hasMore {
        	line = append(line, data...)
            data, hasMore, err = bio.ReadLine()
        }
        line = append(line, data...)
        if err != nil {
            panic(err.Error())
        }
        if string(line) == "end" {
            break
        }
        msgBuf = append(msgBuf, line...)
        msgBuf = append(msgBuf, '\n')
    }
    return msgBuf
}

func (comp *Component) ReadValues(r io.Reader) Values {
	msg := comp.ReadMsg(r)
	vals, err := ParseValues(msg)
    if err != nil {
    	panic(err.Error())
    }
    return vals
}

func (comp *Component) ReadTaskIds() TaskIds {
    if len(comp.PendingTaskids) > 0 {
        var first TaskIds
        first, comp.PendingTaskids = comp.PendingTaskids[0], comp.PendingTaskids[1:]
        return first
    } else {
        msg := comp.ReadMsg(os.Stdin)
        for {
        	vals, err := ParseValues(msg)
            if err != nil {
                break
            }
            comp.PendingCommands = append(comp.PendingCommands, vals)
            msg = comp.ReadMsg(os.Stdin)
        }
        taskIds, err := ParseTaskIds(msg)
        if err != nil {
            panic(err.Error())
        }
        return taskIds
    }
}

func (comp *Component) ReadCommand() Values {
    if len(comp.PendingCommands) > 0 {
        var first Values
        first, comp.PendingCommands = comp.PendingCommands[0], comp.PendingCommands[1:]
        return first
    } else {
        msg := comp.ReadMsg(os.Stdin)
        for {
            taskIds, err := ParseTaskIds(msg)
            if err != nil {
                break
            }
            comp.PendingTaskids = append(comp.PendingTaskids, taskIds)
            msg = comp.ReadMsg(os.Stdin)
        }
        vals, err := ParseValues(msg)
        if err != nil {
            panic(err.Error())
        }
        return vals
    }
}


func (comp *Component) SendMsgToParent(msg interface{}) {
    data, err := json.Marshal(msg)
    if err != nil {
        panic(err.Error())
    }
    comp.Debug("sendtoParent: " + string(data))
    fmt.Fprintf(os.Stdout, string(data)+"\n")
    fmt.Fprintf(os.Stdout, "end\n")
    //os.Stdout.Sync()
}

func (comp *Component) Sync() {
    vals := make(Values)
    vals.Set("command", "sync")
    comp.SendMsgToParent(vals)
}

func (comp *Component) SendPid(heartbeatdir string) {
    pid := syscall.Getpid()
    vals := make(Values)
    vals.Set("pid", pid)
    comp.SendMsgToParent(vals)
    fname := fmt.Sprintf("%s/%d", heartbeatdir, pid)
    ioutil.WriteFile(fname, nil, 06666)
} 

func (comp *Component) Ack(id string) {
    vals := make(Values)
    vals.Set("command", "ack")
    vals.Set("id", id)
    comp.SendMsgToParent(vals)
}

func (comp *Component) Fail(id string) {
    vals := make(Values)
    vals.Set("command", "fail")
    vals.Set("id", id)
    comp.SendMsgToParent(vals)
}

func (comp *Component) Log(msg string) {
    vals := make(Values)
    vals.Set("command", "log")
    vals.Set("msg", msg)
    comp.SendMsgToParent(vals)
}

func (comp *Component) Emit(data []string, stream string, id string, directTask int) {
    m := make(Values)
    m.Set("command", "emit")
    if len(id) > 0 {
    	m.Set("id", id)
    }
    if len(stream) > 0 {
    	m.Set("stream", stream)
    }
    if directTask != 0 {
    	m.Set("task", directTask)	
    }
    m.Set("tuple", data)
    comp.SendMsgToParent(m)
    if comp.Type == "bolt" {
        comp.Log("SendMsgToParent")
    }
    comp.ReadTaskIds()
    if comp.Type == "bolt" {
        comp.Log("ReadTaskIds")
    }
}