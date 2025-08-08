GRPC — Streaming Examples (Go)
A single repository containing five gRPC examples demonstrating the four streaming types and a small helper/shared example.
This README explains everything you need to know: what each example does, how to generate protobuf code, how to run each server & client, expected demo output, troubleshooting tips (Windows-friendly), git push steps, .gitignore, and recommended workflow.

Assumption: your repo root is a folder named grpc and it contains five subdirectories (one per example). If your actual directory names differ, substitute them in the commands below.

Table of contents
Project overview

Directory structure

Prerequisites

Generate protobuf code

How to run each example (detailed)

1. Unary (request/response)

2. Server streaming

3. Client streaming

4. Bidirectional streaming (real-world demo)

5. Shared / helper (shared protos, utils)

Demo outputs (examples)

Windows-specific notes

Troubleshooting & common errors

Recommended git flow & push to GitHub (one-shot commands)

Suggested .gitignore

Should I commit generated .pb.go files?

Contributing, license, contact

Project overview
This repo demonstrates:

Unary RPC — classic request → response

Server streaming — client sends one request, server streams multiple responses

Client streaming — client streams multiple requests, server replies once

Bidirectional streaming — both sides send streams independently (live stock feed example)

Shared — shared .proto files or helper utilities used by one or more examples

Each example is a small, self-contained Go program (server + client). The goal: learn gRPC streaming patterns and run real samples locally.

Directory structure (example)
bash
Copy
Edit
grpc/
├── unary/                   # unary request/response example
│   ├── proto/
│   │   └── greet.proto
│   ├── server/
│   │   └── main.go
│   └── client/
│       └── main.go
├── server_streaming/
│   ├── proto/
│   │   └── news.proto
│   ├── server/main.go
│   └── client/main.go
├── client_streaming/
│   ├── proto/
│   │   └── upload.proto
│   ├── server/main.go
│   └── client/main.go
├── bidi_streaming/
│   ├── proto/
│   │   └── stock.proto
│   ├── server/main.go
│   └── client/main.go
└── shared/
    ├── protos/              # reusable .proto files or common utils
    └── README.md
If your structure is different, follow the same logical flow: locate the .proto, generate code, go mod init per example (or use a workspace).

Prerequisites
Go (1.20+ recommended) installed. Confirm with:

bash
Copy
Edit
go version
protoc (Protocol Buffer Compiler). Confirm with:

bash
Copy
Edit
protoc --version
Go protoc plugins:

bash
Copy
Edit
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
Make sure $GOPATH/bin (or $HOME/go/bin) is in your PATH.

Optional: grpcurl for quick testing (nice-to-have).

Generate protobuf code
For each example, run protoc to generate .pb.go files. Example (run from the example directory or adjust paths):

bash
Copy
Edit
# from grpc/unary (example)
cd grpc/unary
protoc --go_out=. --go-grpc_out=. proto/greet.proto
General pattern:

bash
Copy
Edit
protoc --go_out=. --go-grpc_out=. <path-to-proto-file>
If you have many .proto files or a shared/protos folder, you can run protoc for each file or write a small script.

Generated files will appear next to the package locations (commonly under the same package folder). If you prefer not to commit generated files, document the generation step in README and ensure CI regenerates them.

How to run each example (detailed)
Best practice: run server in one terminal, client in another.

1. Unary (request/response)
Purpose: Demonstrates a single request and single response (e.g., Greet RPC).

Typical files:

proto/greet.proto

server/main.go

client/main.go

Commands:

bash
Copy
Edit
# generate pb
cd grpc/unary
protoc --go_out=. --go-grpc_out=. proto/greet.proto

# in server folder
cd server
go mod init github.com/<your-username>/grpc/unary/server
go mod tidy
go run main.go

# in client folder (new terminal)
cd client
go mod init github.com/<your-username>/grpc/unary/client
go mod tidy
go run main.go
What happens: client sends GreetRequest → server responds with GreetResponse. Useful to check baseline connectivity and service definition.

2. Server streaming
Purpose: Client sends one request (e.g., NewsRequest) and server streams multiple responses (e.g., NewsResponse headlines).

Typical files:

proto/news.proto

server/main.go (server streams multiple messages using stream.Send)

client/main.go (client loops over stream.Recv() until io.EOF)

Commands: same pattern as unary but using the server_streaming directory.

Key code points:

Server: for { stream.Send(&pb.NewsResponse{...}) }

Client: for { res, err := stream.Recv(); if err == io.EOF { break } ... }

3. Client streaming
Purpose: Client sends many messages (e.g., UploadChunk) and server replies once with a summary (e.g., UploadSummary).

Typical files:

proto/upload.proto

client/main.go (calls stream.Send() multiple times then CloseAndRecv())

server/main.go (reads with stream.Recv() until io.EOF, then replies)

Commands: same generation + run pattern.

Key code points:

Client sends multiple messages, then uses stream.CloseAndRecv() to get final server response.

Server accumulates data -> on EOF sends a single result.

4. Bidirectional streaming (real-world demo)
Purpose: Both sides send streams independently (our live stock price subscription example).

Files (example):

proto/stock.proto

server/main.go

client/main.go

Important behaviors:

Client can send new subscription requests at any time (stream.Send()).

Server keeps sending updates for all subscribed symbols (stream.Send(...) periodically).

Both use stream.Recv() concurrently — common pattern is to start a goroutine for receiving while main goroutine handles sending.

How to run (example):

bash
Copy
Edit
cd grpc/bidi_streaming
protoc --go_out=. --go-grpc_out=. proto/stock.proto

# server
cd server
go mod init github.com/<your-username>/grpc/bidi_streaming/server
go mod tidy
go run main.go

# client (new terminal)
cd client
go mod init github.com/<your-username>/grpc/bidi_streaming/client
go mod tidy
go run main.go
Notes: We recommended ctx usage for cancellations and proper CloseSend() from client when exiting.

5. Shared / helper (shared protos, utilities)
Purpose: If you use the same .proto in multiple examples, keep a shared/protos folder. Generate code once and import it in each example (or commit generated pb Go code to reduce friction).

Tip: If multiple examples need the same generated package path, run protoc with the --go_opt=paths=source_relative option to keep paths stable.

Demo outputs (examples)
Below are example outputs you should expect. Your actual numbers and timestamps will differ.

Bidirectional (stock) — Server
yaml
Copy
Edit
🚀 Stock Price Server listening on port 50051
📡 Client connected for stock updates
✅ Subscribed to: AAPL
📤 Sent AAPL update: $743.91
✅ Subscribed to: TSLA
📤 Sent AAPL update: $756.34
📤 Sent TSLA update: $1021.55
❌ Client stopped sending symbols
Bidirectional (stock) — Client
bash
Copy
Edit
📥 Enter stock symbols to subscribe (e.g., AAPL). Type 'exit' to quit:
AAPL
💹 AAPL -> $743.91 at 2025-08-08T14:25:12+05:30
TSLA
💹 AAPL -> $756.34 at 2025-08-08T14:25:14+05:30
💹 TSLA -> $1021.55 at 2025-08-08T14:25:14+05:30
exit
📴 Server closed the stream
Server streaming (news) — Client
pgsql
Copy
Edit
📥 Requesting news for "technology"
📨 Headline: "Startup X raises $100M"
📨 Headline: "New Go release announced"
📨 Headline: "gRPC patterns for real-time apps"
📴 Stream ended by server
Client streaming (uploads) — Server
arduino
Copy
Edit
📤 Receiving chunk 1 (size 512)
📤 Receiving chunk 2 (size 1024)
📤 Receiving chunk 3 (size 256)
✅ Upload complete: 1792 bytes received
Unary (greet) — Client
yaml
Copy
Edit
▶ Request: Hello server
◀ Response: Hello, client — welcome!
Windows-specific notes
Use protoc.exe (download appropriate release) and add its folder to PATH. Example:

powershell
Copy
Edit
setx PATH "%PATH%;C:\tools\protoc\bin"
Restart terminal after setx.

When running protoc on Windows PowerShell, use forward slashes or quote paths. Example:

powershell
Copy
Edit
protoc --go_out=. --go-grpc_out=. .\proto\stock.proto
go mod init and go run work the same on Windows. If git commands fail because of CRLF, configure with:

bash
Copy
Edit
git config --global core.autocrlf true
Troubleshooting & common errors
protoc: command not found

Install protoc and ensure it's in PATH. Verify protoc --version.

protoc-gen-go: program not found or is not executable

Make sure protoc-gen-go is installed and in PATH ($GOPATH/bin or $HOME/go/bin).

bash
Copy
Edit
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
rpc error: code = Unauthenticated or auth errors

You might need to adjust TLS credentials. For local testing you can use insecure for development, but prefer TLS for production:

go
Copy
Edit
import "google.golang.org/grpc/credentials/insecure"
grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
remote: Repository not found or error: remote origin already exists

Make sure the remote URL is correct and you have permissions.

To update an existing origin:

bash
Copy
Edit
git remote set-url origin https://github.com/YOUR_USERNAME/YOUR_REPO.git
To remove and re-add:

bash
Copy
Edit
git remote remove origin
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO.git
go: module declares its path or import path issues

Ensure go mod init module path matches how you import packages, or use replace or go.work for monorepos.

Recommended git flow & push to GitHub (one-shot commands)
Run these (update <your-username> and <repo-name>):

bash
Copy
Edit
cd path/to/grpc             # go to your grpc repo root
git init                    # if not already a git repo
git add .
git commit -m "Initial commit — gRPC streaming examples"

# if remote 'origin' already exists and you want to overwrite:
git remote remove origin
git remote add origin https://github.com/<your-username>/<repo-name>.git

git branch -M main
git push -u origin main
If you prefer to set the remote URL (without removing it):

bash
Copy
Edit
git remote set-url origin https://github.com/<your-username>/<repo-name>.git
git push -u origin main
Suggested .gitignore
gitignore
Copy
Edit
# Go
/bin/
/pkg/
*.exe
*.dll
*.so
*.dylib

# Editor/OS
.vscode/
.idea/
.DS_Store
Thumbs.db

# Generated protobufs (optionally ignore if you will regenerate)
# *.pb.go

# Go workspace / cache
*.cache
You may choose to commit generated .pb.go files so others can run examples without installing protoc. If so, do not list *.pb.go in .gitignore. I explain more below.

Should I commit generated .pb.go files?
Pros:

Easiest for others to run the examples immediately (no protoc required).

Simpler for CI that doesn’t want to install protoc.

Cons:

Generated files can become stale if .proto changes and you forget to regenerate.

Slight repo noise.

