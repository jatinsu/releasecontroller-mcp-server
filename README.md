# ReleaseController MCP Server

The releasecontroller-mcp-server is a robust Model Context Protocol (MCP) server designed to provide comprehensive tooling for interacting with OpenShift release controllers and extracting valuable cluster-related information from Prow job artifacts.

This server acts as a powerful interface, enabling users to programmatically query detailed release information (such as available release streams, latest accepted/rejected releases, failed job summaries, and component versions) and fetch in-depth cluster state details (including pod statuses, cluster operator health, and node specifics) by leveraging the rich data available in Prow job artifacts. It aims to simplify the analysis and debugging of OpenShift and OKD releases and cluster states.

The releasecontroller-mcp-server exposes a wide array of tools categorized into Release Controller interactions and Cluster Information retrieval:

### Release Controller Tools

- List Release Controllers: Get a list of all supported release controllers (e.g., OKD, OpenShift, Multi-arch, ARM64, PPC64LE, S390X).
- Get Specific Release Controller URL: Retrieve the base URL for a particular release controller (e.g., OKD, OpenShift).
- List Release Streams: Enumerate all available release streams within a specified release controller.
- Latest Accepted/Rejected Release: Identify the most recent accepted or rejected release for a given stream.
- List Failed Jobs in Release: Obtain a list of all failed jobs associated with a specific release, including their corresponding Prow job URLs.
- List Components in Release: Display the versions of key components (like kubectl, kubernetes, coreos, and tests) included in a release.
- List Test Failures for Release: Extract and present a summary of failing tests from a given Prow job URL. If no failures are found, a clear message is returned.
- Get Flaky Tests for Release: Identify and list tests that have been marked as flaky within a specific Prow job.
- Get Risk Analysis Data: Fetch the detailed risk analysis data available for a particular Prow job.
- Analyze Job Failures for Release: Download and analyze the build log file for a given Prow job, providing a succinct summary of critical errors and failures. This tool supports log compaction with configurable thresholds (aggressive, moderate, conservative) to manage large logs.
- List Feature Changes: Identify and list feature-related issues from updated image commits within a release, explicitly excluding bugs (OCPBUGS) and CVEs.
- List Bug Fixes: List all bug fixes (specifically OCPBUGS) introduced by updated image commits in a release.
- List CVE Fixes: Enumerate Common Vulnerabilities and Exposures (CVEs) addressed by updated image commits in a release.

### Cluster Information Tools (from Prow Job Artifacts)

These tools leverage gather-extra artifacts from Prow jobs to provide insights into cluster state at the time of the job run.

- Get Pods by State: Retrieve a list of pods in specific states (e.g., CrashLoopBackOff, Pending, Init, Error, Running, or All pods).
- Get Pods by Namespace: Filter and list pods belonging to a particular Kubernetes namespace.
- Get Pods by Node: Identify and list pods scheduled on a specific cluster node.
- Get Containers in Pod: List all container names within a designated pod.
- Get Container Logs: Fetch the logs of a specific container within a pod and analyze them for important events, failures, and errors.
- Get Cluster Operator Status Summary: Provide an overview of the status (Available, Progressing, Degraded) for all cluster operators.
- Get Cluster Version Summary: Detail the cluster's current version, desired version, available updates, and historical update records.
- Get Nodes Info: Retrieve comprehensive details for all cluster nodes, including architecture, OS image, kernel version, and other relevant hardware/software specifics.
- Get Node Info by Name: Obtain detailed information for a specific node by its name.
- Get Node Labels/Annotations by Name: Fetch and display Kubernetes labels or annotations applied to a specified node.
- Get All Nodes Labels/Annotations/Conditions: Retrieve aggregated labels, annotations, or conditions across all nodes in the cluster.



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

Samples query if the stream, release and failing job is known:
- From the OCP release controller, fetch only blocking jobs which have failed for the latest rejected in the 4.19.0-0.nightly stream, use the prow job url for the failing job, clearly list the names of tests that have failed and analyze the logs to see why these particular tests have failed
- From the OKD release controller, fetch all failed jobs for the latest accepted in the 4.20.0-0.okd-scos stream, use the prow job url for the gcp failing job, clearly list the names of tests that have failed and analyze the logs to see why these particular tests have failed
