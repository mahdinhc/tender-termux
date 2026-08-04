// use makechan() for all channels
chan1 := makechan()
chan2 := makechan()
select_chan := makechan() // our central multiplexer

// worker 1
govm(fn() {
    chan1.send(2)
})

// worker 2
govm(fn() {
    chan2.send("pong")
})

// --- multiplexing logic ---

// forwarder for chan1
govm(fn() {
    val := chan1.recv()
    // send an array (or map) containing the source and the value
    select_chan.send(["chan1", val]) 
})

// forwarder for chan2
govm(fn() {
    val := chan2.recv()
    select_chan.send(["chan2", val])
})

// --- the "select" loop ---

for i:=0; i<2; i++ {
    // we only wait on the single multiplexer channel
    msg := select_chan.recv()
    
    source := msg[0]
    val := msg[1]
    
    // handle the cases dynamically based on the source identifier
    if source == "chan1" {
        println("received from chan1:", val)
    } 
	else if source == "chan2" {
        println("received from chan2:", val)
    }
}