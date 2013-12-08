package gostorm

type Spout interface {
	Init(string)
	NextTuple()
	Ack(string)
	Fail(string)
	Sync()
	ReadCommand() Values
	Debug(string, interface{})
}

func RunSpout(s Spout) {
	s.Init("spout")
    for {
    	cmd := s.ReadCommand()
    	if val, ok := cmd.GetString("command"); ok {
    		if val == "next" {
    			s.NextTuple()
			} /*else if val == "ack" {
				if id, ok := cmd.GetString("id"); ok {
					//s.Ack(id)
				}
			} else if val == "fail" {
				if id, ok := cmd.GetString("id"); ok {
					//s.Fail(id)
				}
			}*/
    	}
    	s.Sync()
    }
}