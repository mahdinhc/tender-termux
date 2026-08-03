// tender@v0.0.8

// New syntax
job_queue := chan(1)

worker := go fn() {
    secret_answer := 42
    job_queue <- secret_answer
}()
payload := <- job_queue
println("Received payload: " + string(payload))
worker.wait()

// Legacy method compatibility also works:
ch2 := makechan(1)
w2 := govm(fn() {
    ch2.send(100)
})
println("Legacy payload: " + string(ch2.recv()))
w2.wait()