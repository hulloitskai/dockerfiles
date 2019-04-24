package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	ess "github.com/unixpickle/essentials"
)

const dockerConfigPath = "/kaniko/.docker/config.json"

func main() {
	// Read and validate config.
	cfg, err := newConfigFromEnv()
	if err != nil {
		ess.Die("Error: reading config from env:", err)
	}
	if err = cfg.Validate(); err != nil {
		ess.Die("Error: invalid configuration:", err)
	}

	// Open and rewrite Docker config file.
	file, err := os.OpenFile(dockerConfigPath, os.O_RDWR, 0600)
	defer file.Close()
	if err = cfg.EditDockerConfig(file); err != nil {
		ess.Die("Error: editing Docker config:", err)
	}

	if cfg.PluginDryRun {
		if _, err = file.Seek(0, io.SeekStart); err != nil {
			ess.Die("Error: seeking file start:", err)
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			ess.Die("Error: reading file (dry-run):", err)
		}
		fmt.Fprintf(os.Stderr, "Docker config file contents: '%q'\n", data)
	}

	if err = file.Close(); err != nil {
		ess.Die("Error: closing Docker config file:", err)
	}

	// Build and execute command.
	cmd, err := cfg.Command()
	if err != nil {
		ess.Die("Error: building Kaniko command:", err)
	}

	// During dry-run, print command to stderr and exit.
	if cfg.PluginDryRun {
		fmt.Fprintf(
			os.Stderr,
			"Should run command '%s' with arguments: %v\n",
			cmd.Path, cmd.Args,
		)
		os.Exit(0)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		if eerr, ok := err.(*exec.ExitError); ok {
			os.Exit(eerr.ExitCode())
		}
		ess.Die("Error: running Kaniko command:", err)
	}
}
