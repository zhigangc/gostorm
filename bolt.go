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
    	tup := b.ReadTuple()
        b.SetAnchor(tup)
    	b.Process(tup)
    	b.Ack(tup.Id)
    }
}