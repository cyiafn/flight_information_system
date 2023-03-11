# How to run the server.

1. Install go1.19
2. Navigate to the root directory in your terminal/commandprompt/powershell.
3. Run `go run main.go -amo true` to run the server in at most once invocation mode and `go run main.go -amo false` in at least once invocation mode.

# Building it for distribution
1. Install go1.19
2. Navigate to root directory in your terminal/commandprompt/powershell.
3. Run `go build -o output`. This will build a binary (on macOS and linux) or executable (windows) into the output folder for your operating system and system architecture.