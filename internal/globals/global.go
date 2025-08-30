package globals

const GithubDockerHost = "ghcr.io"

const RepoOwner = "seqrateam"

const AnalyzerDocker = GithubDockerHost + "/" + RepoOwner + "/seqra-jvm-sast/sast-analyzer"
const AnalyzerBindVersion = "2025.08.29.2bedbbb"

const AutobuilderRepoName = "seqra-jvm-autobuilder"
const AutobuilderDocker = GithubDockerHost + "/" + RepoOwner + "/" + AutobuilderRepoName + "/sast-autobuilder"
const AutobuilderBindVersion = "2025.08.28.6def5f5"
const AutobuilderAssetName = "seqra-project-auto-builder.jar"

const RulesRepoName = "seqra-rules"
const RulesBindVersion = "v1.0.0"

var AnalyzerVersion string
var AutobuilderVersion string
var VerboseLevel string
var Quiet bool
var GithubToken string

var CompileType string

var LogPath string
