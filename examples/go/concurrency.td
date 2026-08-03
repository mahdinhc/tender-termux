// tender@v0.0.8
// ----- Simple goroutine -----
g := go fn() {
    println("Hello from goroutine!")
    return 42   // return a value
}()

// Wait for completion (blocks until done)
g.wait()

// Get the return value
result := g.result()
println("Goroutine returned:", result)  // 42

// Abort a goroutine (if needed, e.g., from another goroutine)
// g.abort()

// ----- Channels -----
ch := chan(10)  // buffered channel with capacity 10

// Send and receive
ch <- "hello"
msg := <-ch
println(msg)  // "hello"
ch.close() // Close the channel (sends no more values)

// ----- Concurrent worker pool -----
fn worker(id, jobs, results) {
    for {
        job := <-jobs        // blocks until a value is available
        if is_null(job) {    // channel closed => receive returns null
            break
        }
		results <- job * 2   // send result
    }
}

jobs := chan(100)
results := chan(100)

// Start 10 workers
for i := 0; i < 10; i++ {
    // govm(worker, i, jobs, results
	go worker(i, jobs, results)
}

// Send 50 jobs
for i := 0; i < 50; i++ {
    jobs <- i
}
jobs.close()  // signal workers to stop

// Collect results (order may vary)
for i := 0; i < 50; i++ {
    res := <-results
    println(res)
}