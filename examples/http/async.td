// tender@v0.0.8
import "http"
import "json"

async_fetch := go fn() {
    req := http.get("https://api.github.com")
    resp := req.execute()
    
    return resp.body
}()

final_status := async_fetch.result()

debug(final_status |> json.decode)