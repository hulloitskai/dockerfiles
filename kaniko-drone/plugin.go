package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	TagFiles   []string `split_words:"true"`
	TarPath    string   `split_words:"true"`
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
	if (len(cfg.Tags)+len(cfg.TagFiles) > 0) && (cfg.Repo == "") {
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
	for _, arg := range cfg.BuildArgs {
		args = append(args, "--build-arg", arg)
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

	// Parse tag files.
	var (
		tags  = cfg.Tags
		regex = regexp.MustCompile(tagRegex)
	)
	for _, name := range cfg.TagFiles {
		tag, err := cfg.readTagFile(name, regex)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	// Construct destination flags.
	for _, tag := range tags {
		args = append(args, "--destination", fmt.Sprintf("%s:%s", cfg.Repo, tag))
	}
	return exec.Command(executorCommand, args...), nil
}

func (cfg *config) readTagFile(name string, r *regexp.Regexp) (string, error) {
	if !filepath.IsAbs(name) && (cfg.Context != "") {
		name = filepath.Join(cfg.Context, name)
	}
	rawtag, err := ioutil.ReadFile(name)
	if err != nil {
		return "", errors.Errorf("reading tag from '%s': %w", name, err)
	}
	rawtag = bytes.TrimSpace(rawtag)
	if !r.Match(rawtag) {
		return "", errors.Errorf("'%q' is not a valid tag", name)
	}
	return string(rawtag), nil
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
