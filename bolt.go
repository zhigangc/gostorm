package gostorm

type Bolt interface {
	Init(string)
	ReadTuple()*Tuple
	Process(*Tuple)
	Debug(string)
	Ack(string)
    Log(string)
}

func RunBolt(b Bolt) {
    b.Init("bolt")
    for {
        b.Log("bolt.ReadTuple")
    	tup := b.ReadTuple()
        b.Log("bolt.Process")
    	b.Process(tup)
        b.Log("bolt.Ack")
    	b.Ack(tup.Id)
    }
}