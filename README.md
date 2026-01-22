# Base Upgrader CLI

This is a CLI tool for upgrading applications to newer base app versions. It will compare the diff between the 
old/new base app versions to the files in your target repo, and generate a basic report.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/c3-trentbuckholz/baseupgrader.git && cd baseupgrader
   ```
2. Install go [via the go website](https://go.dev/doc/install)
3. Bulid the tool:
   ```bash
   make build
   ```
4. Run the tool:
   ```bash
   ./bin/baseupgrader --help
   ```

## Usage

Run with `--help` to see all available options.

### Example Command

```bash
bin/baseupgrader --baseRepoUrl="https://github.com/c3-e/c3pso" --oldCommit="2758cdb5f3fc5e854616ea9c83a8771137fffc5c" --newCommit="28bcf7ac5b81feac20e4e9d70a6e0658b75454b6" --targetR
epoUrl="https://github.com/c3-e/c3fed-hii" --ghAuthToken="<GITHUB_AUTH_TOKEN>" --targetRepoPath=apps/nns/pso --outputType=html

