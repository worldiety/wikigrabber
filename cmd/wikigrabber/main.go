package main

import (
	"flag"
	"fmt"
	"github.com/worldiety/wikigrabber/internal/app"
	"github.com/worldiety/wikigrabber/internal/config"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if err := start(); err != nil {
		panic(err)
	}
}

func start() error {
	log.SetFlags(0) // we don't want extra timestamp, systemd et al cares already for us

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	tmpDir := filepath.Join(os.TempDir(), "wikigrabber")
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create tmp dir")
	}

	cfgFile := flag.String("cfg", filepath.Join(homeDir, ".wikigrabber.yaml"), "the yaml file to configure the grabber")
	cmdTmpDir := flag.String("tmpDir", tmpDir, "the temporary folder to clone and pull remotes into")
	help := flag.Bool("help", false, "shows this help")

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return nil
	}

	cfg := config.Config{}
	cfg.TmpDir = *cmdTmpDir
	if err := cfg.Load(*cfgFile); err != nil {
		return err
	}

	return app.Process(cfg)
}
