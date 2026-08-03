// Create a 1-slot buffered channel using new syntax: chan(size)
job_queue := chan(1)

// Spawn background worker using new syntax: go fn()
worker := go fn() {
    // Simulate some heavy math or database lookup
    secret_answer := 42
    
    // Send to channel using new syntax: ch <- val
    job_queue <- secret_answer
}()

// Main script thread receives payload using new syntax: <-ch
payload := <- job_queue
println("Received payload: " + string(payload))

// Verify worker handle wait
worker.wait()

// Legacy method compatibility also works:
ch2 := makechan(1)
w2 := govm(fn() {
    ch2.send(100)
})
println("Legacy payload: " + string(ch2.recv()))
w2.wait()