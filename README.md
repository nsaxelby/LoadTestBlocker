# LoadTestBlocker
Load testing a URL with a L7 Flood.


## Running
Optional: run `docker compose up -d` to start the local httpbin (provides a local test destination)

Run the web app
`go run ./cmd/web/main.go`

Then visit `http://localhost:8080` in a browser.

### Running a load test
Enter in the URL you want to load test, and press the 'Start Test' button.

Note: the 'rate tracking' feature doesn't work too well at the moment, Issue is open: https://github.com/nsaxelby/LoadTestBlocker/issues/1

### Screenshot
![Screenshot1](screenshot1.png)




