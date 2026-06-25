package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/fileutil"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
)

const (
	sourceDirEnvKey = "BITRISE_SOURCE_DIR"
	pwdEnvKey       = "PWD"
)

// Input holds raw step inputs parsed from environment variables.
type Input struct {
	Path         string `env:"path,required"`
	IsCreatePath bool   `env:"is_create_path,opt[true,false]"`
	UseSudo      bool   `env:"use_sudo,opt[true,false]"`
}

// Config holds validated step configuration.
type Config struct {
	path         string
	isCreatePath bool
	useSudo      bool
}

// Step changes the working directory for subsequent steps.
type Step struct {
	inputParser  stepconf.InputParser
	logger       log.Logger
	pathChecker  pathutil.PathChecker
	pathModifier pathutil.PathModifier
	exporter     export.Exporter
	cmdFactory   command.Factory
}

// NewStep creates a Step with injected dependencies.
func NewStep(inputParser stepconf.InputParser, logger log.Logger, pathChecker pathutil.PathChecker, pathModifier pathutil.PathModifier, exporter export.Exporter, cmdFactory command.Factory) Step {
	return Step{
		inputParser:  inputParser,
		logger:       logger,
		pathChecker:  pathChecker,
		pathModifier: pathModifier,
		exporter:     exporter,
		cmdFactory:   cmdFactory,
	}
}

// ProcessConfig parses and validates step inputs.
func (s Step) ProcessConfig() (Config, error) {
	var input Input
	if err := s.inputParser.Parse(&input); err != nil {
		return Config{}, fmt.Errorf("failed to parse inputs: %w", err)
	}
	stepconf.Print(input)
	return Config{
		path:         input.Path,
		isCreatePath: input.IsCreatePath,
		useSudo:      input.UseSudo,
	}, nil
}

// Run ensures the target directory exists and returns its expanded absolute path.
func (s Step) Run(cfg Config) (string, error) {
	exists, err := s.pathChecker.IsDirExists(cfg.path)
	if err != nil {
		return "", fmt.Errorf("failed to check if directory exists (%s): %w", cfg.path, err)
	}

	if !exists {
		if !cfg.isCreatePath {
			return "", fmt.Errorf("directory (%s) does not exist and creating it is disabled (is_create_path=false)", cfg.path)
		}

		s.logger.Printf("Directory (%s) does not exist, creating it...", cfg.path)
		if err := s.createDir(cfg.path, cfg.useSudo); err != nil {
			return "", fmt.Errorf("failed to create directory (%s): %w", cfg.path, err)
		}
	}

	absPath, err := s.pathModifier.AbsPath(cfg.path)
	if err != nil {
		return "", fmt.Errorf("failed to expand path (%s): %w", cfg.path, err)
	}

	return absPath, nil
}

// createDir creates the directory (including parents), falling back to an
// elevated mkdir (when useSudo is set) if the step lacks permission. On stacks
// where the step runs as a non-root user, creating a directory inside a
// root-owned parent requires sudo. `-n` prevents sudo from blocking on a
// password prompt.
func (s Step) createDir(path string, useSudo bool) error {
	err := os.MkdirAll(path, 0755)
	if err == nil {
		return nil
	}
	if !useSudo || !errors.Is(err, fs.ErrPermission) {
		return err
	}

	s.logger.Warnf("Permission denied creating %s, retrying with sudo", path)
	var stderr bytes.Buffer
	cmd := s.cmdFactory.Create("sudo", []string{"-n", "mkdir", "-p", path}, &command.Opts{
		Stderr: &stderr,
	})
	if runErr := cmd.Run(); runErr != nil {
		return fmt.Errorf("elevated directory creation failed (%s): %w", strings.TrimSpace(stderr.String()), runErr)
	}
	return nil
}

// ExportOutputs sets BITRISE_SOURCE_DIR (and PWD) so subsequent steps run in the new directory.
func (s Step) ExportOutputs(workdir string) error {
	s.logger.Donef("Changing working directory to %s", workdir)
	for _, key := range []string{sourceDirEnvKey, pwdEnvKey} {
		if err := s.exporter.ExportOutput(key, workdir); err != nil {
			return fmt.Errorf("failed to export %s: %w", key, err)
		}
	}
	return nil
}

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger()
	envRepository := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepository)
	cmdFactory := command.NewFactory(envRepository)
	step := NewStep(inputParser, logger, pathutil.NewPathChecker(), pathutil.NewPathModifier(), export.NewExporter(cmdFactory, fileutil.NewFileManager()), cmdFactory)

	cfg, err := step.ProcessConfig()
	if err != nil {
		logger.Errorf(err.Error())
		return 1
	}

	workdir, err := step.Run(cfg)
	if err != nil {
		logger.Errorf(err.Error())
		return 1
	}

	if err := step.ExportOutputs(workdir); err != nil {
		logger.Errorf(err.Error())
		return 1
	}

	return 0
}
