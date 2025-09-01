[![GitHub release](https://img.shields.io/github/release/seqrateam/seqra.svg)](https://github.com/seqrateam/seqra/releases)

# Seqra — security-focused static analyzer for Java

[Issues](https://github.com/seqrateam/seqra/issues) | [FAQ](docs/faq.md) | [Discord](https://discord.gg/FtKRPv8n) | [seqradev@gmail.com](mailto:seqradev@gmail.com)

#### Why Seqra?

* **CodeQL power + Semgrep simplicity**:
  - Write security rules using familiar patterns while getting cross-module dataflow analysis
* **Free and source-available**:
  - Use for any purpose except competing commercial offerings for free
* **Workflow ready**:
  - CLI tool with SARIF output for seamless CI/CD integration



#### Table of Contents
* [License](#license)
* [Quick Start](#quick-start)
* [CI/CD Integration](#cicd-integration)
* [Troubleshooting](#troubleshooting)
* [Changelog](#changelog)

## License

This project is released under the MIT License.

The [core analysis engine](https://github.com/seqrateam/seqra-jvm-sast) is source-available under the [Functional Source License (FSL-1.1-ALv2)](https://fsl.software/), which converts to Apache 2.0 two years after each release. You can use Seqra for free, including for commercial use, except for competing products or services.

## Quick Start

### Prerequisites

- **Docker** (used to run the analysis engine in a container)
  - [Install Docker](https://docs.docker.com/get-started/get-docker/)
  - Ensure Docker is running and accessible from your terminal
  - *For Apple Silicon Mac users*: The Docker image is currently built for x86_64/amd64 architecture. [Enable x86_64/amd64 emulation in Docker Desktop](https://docs.docker.com/desktop/settings/mac/#general) to run x86 containers on your ARM-based Mac.

### 1. Install Seqra CLI

- ### Option A: Download Pre-built Binary (Linux)

  *One-liner install:*

  ```bash
  curl -L https://github.com/seqrateam/seqra/releases/latest/download/seqra_linux_amd64.tar.gz -o seqra.tar.gz && tar -xzf seqra.tar.gz seqra && sudo mv seqra /usr/local/bin/ && rm seqra.tar.gz && seqra --version
  ```

  *Step-by-step:*

  ```bash
  # 1. Download
  curl -L https://github.com/seqrateam/seqra/releases/latest/download/seqra_linux_amd64.tar.gz -o seqra.tar.gz

  # 2. Extract
  tar -xzf seqra.tar.gz seqra

  # 3. Install globally (optional)
  sudo mv seqra /usr/local/bin/

  # 4. Remove archive
  rm seqra.tar.gz

  # 5. Verify
  seqra --version
  ```

- ### Option B: Install via Go (Linux/macOS)


  > **Note:** **Support Apple Silicon Mac is experemental** you need [Enable x86_64/amd64 emulation in Docker Desktop](https://docs.docker.com/desktop/settings/mac/#general)

  Install
  ```bash
  go install github.com/seqrateam/seqra@latest
  ```
  
  Verify

  ```bash
  $(go env GOPATH)/bin/seqra --version
  ```

  > **Optional:** Add `GOPATH` to path

  * bash
    ```bash
    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc && source ~/.bashrc
    ```
  * zsh (macOS)
    ```bash
    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc && source ~/.zshrc
    ```


### 2. Run Your First Scan

Scan a Java project and generate SARIF report

```bash
seqra scan --output results.sarif /path/to/your/java/project
```

### 3. View and Analyze Results

Seqra generates results in the standard SARIF format, which can be viewed and analyzed in multiple ways:

- #### **VS Code Integration**

  Open `results.sarif` with the [SARIF Viewer](https://marketplace.visualstudio.com/items?itemName=MS-SarifVSCode.sarif-viewer) extension for a rich, interactive experience.

- #### **GitHub Integration**

  Upload results to [GitHub code scanning](https://docs.github.com/en/code-security/code-scanning/integrating-with-code-scanning/uploading-a-sarif-file-to-github) for security alerts and code quality insights.

- #### **Command Line Preview**

  Quick preview of findings

  ```bash
  seqra summary --show-findings results.sarif
  ```

- #### **CodeChecker Integration**

  Use [CodeChecker](https://github.com/Ericsson/codechecker) for advanced result management, tracking, and team collaboration.


## CI/CD Integration

For seamless integration with your CI/CD pipelines, check out our dedicated integration repositories:

- **[seqra-action](https://github.com/seqrateam/seqra-action)** - GitHub Action for easy integration with GitHub workflows
- **[seqra-gitlab](https://github.com/seqrateam/seqra-gitlab)** - GitLab CI template for automated security scanning


## Troubleshooting

### Docker not running
  - Ensure Docker is installed and running on your system
  - Run `docker info` to verify Docker is accessible

### Build Issues
  > **Note:** **only Maven and Gradle projects are supported**
  - Ensure your Java project builds successfully with its native build tools
  - If the Docker image lacks required dependencies, use `seqra scan --compile-type native --output /path/project/model /path/to/your/project` to build the project directly on your machine instead

### Logs and Debugging
  - Run with `--verbosity debug` for detailed logs
  - Check the log file at `~/.seqra/logs/`

## Changelog
See [CHANGELOG](CHANGELOG.md).
