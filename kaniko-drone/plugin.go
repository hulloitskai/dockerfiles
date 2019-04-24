package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/kelseyhightower/envconfig"
	errors "golang.org/x/xerrors"
)

const (
	dockerRegistry  = "https://index.docker.io/v1/"
	executorCommand = "/kaniko/executor"
	tagRegex        = `[\w][\w.-]{0,127}`
)

type config struct {
	Registry   string
	Username   string
	Password   string
	Context    string
	Dockerfile string
	BuildArgs  []string `split_words:"true"`
	Cache      bool
	CacheDir   string `split_words:"true"`
	CacheRepo  string `split_words:"true"`
	NoPush     bool   `split_words:"true"`
	Repo       string
	Tags       []string
	TarPath    string `split_words:"true"`
	Verbosity  string

	PluginDryRun     bool `split_words:"true"` // for internal use
	DisableGCRHelper bool `envconfig:"disable_gcr_helper"`
}

func defaultConfig() *config {
	return &config{
		Registry: dockerRegistry,
		Context:  ".", // current directory
	}
}

func newConfigFromEnv() (*config, error) {
	cfg := defaultConfig()
	if err := envconfig.Process("plugin", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *config) Validate() error {
	if (len(cfg.Tags) > 0) && (cfg.Repo == "") {
		return errors.New("tag(s) were specified, but no repo was set")
	}
	// Context must be an absolute path.
	if !filepath.IsAbs(cfg.Context) {
		var err error
		cfg.Context, err = filepath.Abs(cfg.Context)
		if err != nil {
			return errors.Errorf("determining absolute context path: %w", err)
		}
	}
	return nil
}

func (cfg *config) Command() (*exec.Cmd, error) {
	args := []string{"--context", cfg.Context}
	if cfg.Cache {
		args = append(args, "--cache")
	}
	if cfg.NoPush {
		args = append(args, "--no-push")
	}
	if cfg.Dockerfile != "" {
		args = append(args, "--dockerfile", cfg.Dockerfile)
	}
	if cfg.CacheDir != "" {
		args = append(args, "--cache-dir", cfg.CacheDir)
	}
	if cfg.TarPath != "" {
		args = append(args, "--tarPath", cfg.TarPath)
	}
	if cfg.Verbosity != "" {
		args = append(args, "--verbosity", cfg.Verbosity)
	}

	// Parse build args using shell.
	for _, arg := range cfg.BuildArgs {
		val, err := cfg.resolveWithShell(arg)
		if err != nil {
			return nil, errors.Errorf("resolving build arg '%s': %w", arg, err)
		}
		args = append(args, "--build-arg", val)
	}

	// Parse and validate tags and destinations.
	regex := regexp.MustCompile(tagRegex)
	for _, tag := range cfg.Tags {
		val, err := cfg.resolveWithShell(tag)
		if err != nil {
			return nil, errors.Errorf("resolving tag '%s': %w", tag, err)
		}
		if !regex.MatchString(val) {
			return nil, errors.Errorf("'%q' is not a valid tag", val)
		}
		args = append(args, "--destination", fmt.Sprintf("%s:%s", cfg.Repo, val))
	}
	return exec.Command(executorCommand, args...), nil
}

func (cfg *config) resolveWithShell(value string) (string, error) {
	var (
		wrapper = fmt.Sprintf("VAL=\"%s\"; printf %%s \"$VAL\"", value)
		cmd     = exec.Command("sh", "-ec", wrapper)
	)
	cmd.Dir = cfg.Context
	raw, err := cmd.Output()
	if err != nil {
		if eerr, ok := err.(*exec.ExitError); ok && (len(eerr.Stderr) > 0) {
			stderr := bytes.TrimSpace(eerr.Stderr)
			return "", errors.Errorf("%v: %s", eerr, stderr)
		}
		return "", err
	}
	return string(bytes.TrimSpace(raw)), nil
}

func (cfg *config) EditDockerConfig(file *os.File) error {
	type authConfig struct {
		Auth string `json:"auth"`
	}
	var dockerConfig struct {
		Auths       map[string]authConfig `json:"auths"`
		CredHelpers map[string]string     `json:"credHelpers,omitempty"`
	}

	// Read config from rw.
	if err := json.NewDecoder(file).Decode(&dockerConfig); err != nil {
		return errors.Errorf("reading file as JSON: %w", err)
	}

	// Truncate file, and reset file cursor to beginning.
	if err := file.Truncate(0); err != nil {
		return errors.Errorf("truncating file: %w", err)
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return errors.Errorf("seeking file start: %w", err)
	}

	// Update registry auths, cred helpers.
	var (
		creds    = fmt.Sprintf("%s:%s", cfg.Username, cfg.Password)
		b64creds = base64.StdEncoding.EncodeToString([]byte(creds))
	)
	dockerConfig.Auths[cfg.Registry] = authConfig{b64creds}
	if cfg.DisableGCRHelper {
		dockerConfig.CredHelpers = nil
	}

	// Save config to rw.
	if err := json.NewEncoder(file).Encode(&dockerConfig); err != nil {
		return errors.Errorf("writing to file: %w", err)
	}
	return nil
}
