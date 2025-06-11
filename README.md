# ReleaseController MCP Server

Inspired by https://github.com/manusa/podman-mcp-server

This MCP server provides tools to interact with the OKD and Openshift release controllers to:
- provide the list of release streams and releases within the stream
- Latest accepted release in a release stream
- Latest rejected release in a release stream
- List failed jobs in a release run
- Identifying failures for a given job
- List CVE fixes in a given release
- List BUGS fixed in a given release
- List feature changes in a given release

Getting Started (with goose AI agent):

1. Clone the repo:
```
git clone git@github.com:Prashanth684/releasecontroller-mcp-server.git
```
2. Build:
```
make build
```
3. Run in SSE mode:
```
./releasecontroller-mcp-server --sse-port 8080
```
4. Add your MCP server to the goose config file (~/.config/goose/config.yaml)
```
GOOSE_MODEL: gemini-2.0-flash
extensions:
  releasecontroller:
    description: null
    enabled: true
    envs: {}
    name: releasecontroller
    timeout: 300
    type: sse
    uri: http://0.0.0.0:8080/sse
```
5. Start goose
```
goose session
```

Sample query flow:
- Find the latest accepted release in the 4.20.0-0.okd-scos stream
- List the failed jobs in this release
- For the gcp job, look at logs and list the failures

Samples query flow if the stream, release and failing job is known:
- From the OKD release controller, fetch all failed jobs for release 4.20.0-0.okd-scos-2025-06-09-184607 in the 4.20.0-0.okd-scos stream, use the prow job url for the gcp job, list the failures and generate a succinct report not exceeding 1000 words
