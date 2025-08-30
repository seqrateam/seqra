package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/seqra/seqra/internal/container_run"
	"github.com/seqra/seqra/internal/globals"
	"github.com/seqra/seqra/internal/utils"
	"github.com/seqra/seqra/internal/utils/log"
)

var OutputProjectModelPath string
var ProjectPath string

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile <project>",
	Short: "Compile your Java project",
	Args:  cobra.MinimumNArgs(1), // require at least one argument
	Long: `This command takes a required path to the project, automatically detects Java build system, modules and dependencies and compile project model.

Arguments:
  <project>  - the path to the project to compile (required)
`,
	Run: func(cmd *cobra.Command, args []string) {
		ProjectPath = args[0]
		compile(OutputProjectModelPath, ProjectPath, globals.CompileType)
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)

	compileCmd.Flags().StringVarP(&OutputProjectModelPath, "output", "o", "", `A path to the result project model`)
	compileCmd.Flags().StringVar(&globals.CompileType, "compile-type", "docker", "Environment for run compile command (docker, native)")
	_ = compileCmd.MarkFlagRequired("output")
}

func compile(outputProjectModelPath, projectRoot, compileType string) {
	logrus.Infof("Logging to file: %s", globals.LogPath)
	outputProjectModelPath = filepath.Clean(outputProjectModelPath)
	absOutputProjectModelPath := log.AbsPathOrExit(outputProjectModelPath, "output")

	projectRoot = filepath.Clean(projectRoot)
	absProjectRoot := log.AbsPathOrExit(projectRoot, "project path")

	logrus.Infof("Build project: %s", absProjectRoot)
	logrus.Infof("Result write to: %s", absOutputProjectModelPath)

	if _, err := os.Stat(absOutputProjectModelPath); err == nil {
		logrus.Fatalf("Output directory already exist: %s", absOutputProjectModelPath)
	}

	appendFlags := []string{}

	switch globals.VerboseLevel {
	case "info":
		appendFlags = append(appendFlags, "--verbosity=info")
	case "debug":
		appendFlags = append(appendFlags, "--verbosity=debug")
	}

	switch compileType {
	case "docker":
		compileWithDocker(absOutputProjectModelPath, absProjectRoot, appendFlags)
	case "native":
		compileWithNative(absOutputProjectModelPath, absProjectRoot, appendFlags)
	default:
		logrus.Fatalf("compile-type must be one of \"docker\", \"native\"")
	}
}

func compileWithDocker(absOutputProjectModelPath, absProjectRoot string, appendFlags []string) {
	autobuilderFlags := []string{
		"--project-root-dir", "/data/project",
		"--build", "portable",
		"--result-dir", "/data/build",
	}

	autobuilderFlags = append(autobuilderFlags, appendFlags...)

	hostConfig := &container.HostConfig{}

	// Get the current user's UID and GID
	containerUID := fmt.Sprintf("%d", os.Getuid())
	containerGID := fmt.Sprintf("%d", os.Getgid())

	envCont := []string{"CONTAINER_UID=" + containerUID, "CONTAINER_GID=" + containerGID}

	var copyToContainer = make(map[string]string)
	copyToContainer[absProjectRoot] = "/data/project"

	var copyFromContainer = make(map[string]string)
	copyFromContainer["/data/build"] = absOutputProjectModelPath

	autobuilderImageLink := utils.GetImageLink(globals.AutobuilderVersion, globals.AutobuilderDocker)
	container_run.RunGhcrContainer(autobuilderImageLink, autobuilderFlags, envCont, hostConfig, copyToContainer, copyFromContainer)
}

func compileWithNative(absOutputProjectModelPath, absProjectRoot string, appendFlags []string) {
	autobuilderJarPath, err := utils.GetAutobuilderJarPath(globals.AutobuilderVersion)
	if err != nil {
		logrus.Fatalf("Unexpected error occurred while trying to construct path to the autobuilder: %s", err)
	}

	if _, err := os.Stat(autobuilderJarPath); errors.Is(err, os.ErrNotExist) {
		err := utils.DownloadGithubReleaseAsset(globals.RepoOwner, globals.AutobuilderRepoName, globals.AutobuilderVersion, globals.AutobuilderAssetName, autobuilderJarPath, globals.GithubToken)
		if err != nil {
			logrus.Fatalf("Unexpected error occurred while trying to download autobuilder: %s", err)
		}
	}

	autobuilderCommand := []string{
		"-Xmx1G",
		"-jar",
		autobuilderJarPath,
		"--project-root-dir", absProjectRoot,
		"--build", "portable",
		"--result-dir", absOutputProjectModelPath,
	}
	autobuilderCommand = append(autobuilderCommand, appendFlags...)

	cmd := exec.Command("java", autobuilderCommand...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		logrus.Fatalf("Unexpected error occurred while trying to start autobuilder: %s", err)
	}
	err = cmd.Wait()
	if err != nil {
		logrus.Fatalf("Unexpected error occurred while trying to wait autobuilder: %s", err)
	}
}
