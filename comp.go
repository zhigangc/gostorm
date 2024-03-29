package gostorm

import "os"
import "encoding/json"
import "syscall"
import "fmt"
import "io/ioutil"
import "bufio"

type Component struct {
	Conf Values
	Context Values
	PendingCommands []Values
	PendingTaskids []TaskIds
	debug *bufio.Writer
    Type string
    Anchor *Tuple
    reader *bufio.Reader
}

func (comp *Component) Debug(prefix string, msg interface{}) {
    if msg != nil {
        if str, ok := msg.(string); ok {
            comp.debug.WriteString(prefix+ ": " + str + "\n")
        } else {
        	data, err := json.Marshal(msg)
            if err != nil {
                panic(err.Error())
            }
            comp.debug.WriteString(prefix+ ": " + string(data) + "\n")
        }
    } else {
        comp.debug.WriteString(prefix+ ": nil\n")
    }
    
	comp.debug.Flush()
}

func (comp *Component) Init(typ string) {
	tmpfile, err := os.OpenFile("/tmp/go.out."+typ, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	comp.debug = bufio.NewWriter(tmpfile)
    comp.reader = bufio.NewReader(os.Stdin)

	setupInfo := comp.ReadValues()
    pidDir, _ := setupInfo.GetString("pidDir")
    comp.SendPid(pidDir)
    conf, _ := setupInfo.GetValues("conf")
    context, _ := setupInfo.GetValues("context")
    comp.Conf = conf
    comp.Context = context
    comp.Type = typ
    comp.reader = bufio.NewReader(os.Stdin)
}

func (comp *Component) SetAnchor(anchor *Tuple) {
    comp.Anchor = anchor
}

func (comp *Component) Run() {}
func (comp *Component) Process(tup *Tuple) {}

func (comp *Component) ReadTuple() *Tuple {
    vals := comp.ReadCommand()
    tup := NewTuple(vals, comp)
    return tup
}

func (comp *Component) ReadMsg() []byte {
    var msgBuf []byte
    for {
        var line []byte
        data, hasMore, err := comp.reader.ReadLine()
        for hasMore {
        	line = append(line, data...)
            data, hasMore, err = comp.reader.ReadLine()
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

func (comp *Component) ReadValues() Values {
	msg := comp.ReadMsg()
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
        msg := comp.ReadMsg()
        for {
        	vals, err := ParseValues(msg)
            if err != nil {
                break
            }
            comp.PendingCommands = append(comp.PendingCommands, vals)
            msg = comp.ReadMsg()
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
        msg := comp.ReadMsg()
        for {
            taskIds, err := ParseTaskIds(msg)
            if err != nil {
                break
            }
            comp.PendingTaskids = append(comp.PendingTaskids, taskIds)
            msg = comp.ReadMsg()
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
    fmt.Fprintf(os.Stdout, string(data)+"\n")
    fmt.Fprintf(os.Stdout, "end\n")
    os.Stdout.Sync()
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

func (comp *Component) Emit(data []interface{}, stream string, id string, directTask int) {
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
    anchorIds := make([]string, 0, 1)
    if comp.Anchor != nil {
        anchorIds = append(anchorIds, comp.Anchor.Id)
    }
    m.Set("anchors", anchorIds)
    m.Set("tuple", data)
    comp.SendMsgToParent(m)
    comp.ReadTaskIds()
}