package gostorm

type Bolt interface {
	Init(string)
	ReadTuple()*Tuple
	Process(*Tuple)
	Debug(string)
	Ack(string)
    Log(string)
    SetAnchor(*Tuple)
}

func RunBolt(b Bolt) {
    b.Init("bolt")
    for {
        b.Debug("bolt.ReadTuple")
    	tup := b.ReadTuple()
        b.SetAnchor(tup)
        b.Debug("bolt.Process")
    	b.Process(tup)
        b.Debug("bolt.Ack")
    	b.Ack(tup.Id)
    }
}