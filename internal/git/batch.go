package git

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// CloneOrPull detects what to do and either clones or pulls a repo
func CloneOrPull(remote, targetDir string) error {
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to ensure target director: %w", err)
	}

	clone := true
	if stat, err := os.Stat(filepath.Join(targetDir, ".git")); err == nil {
		if !stat.IsDir() {
			return fmt.Errorf("found non-.git file")
		}
		clone = false
	}

	if clone {
		return cloneRepo(remote, targetDir)
	}

	return pullRepo(targetDir)
}

func cloneRepo(remote, targetDir string) error {
	log.Printf("cd %s\n", targetDir)
	log.Printf("git clone %s .\n", remote)

	cmd := exec.Command("git", "clone", "--progress", remote, ".")
	cmd.Env = os.Environ()
	cmd.Dir = targetDir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func pullRepo(targetDir string) error {
	log.Printf("cd %s\n", targetDir)
	log.Printf("git pull")

	cmd := exec.Command("git", "pull", "--progress")
	cmd.Env = os.Environ()
	cmd.Dir = targetDir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
