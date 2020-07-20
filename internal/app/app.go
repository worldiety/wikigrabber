package app

import (
	"fmt"
	"github.com/worldiety/wikigrabber/internal/config"
	"github.com/worldiety/wikigrabber/internal/git"
	"github.com/worldiety/wikigrabber/internal/index"
	"github.com/worldiety/wikigrabber/internal/markdown"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Process pulls all repositories and creates a meta page with hugo
func Process(cfg config.Config) error {
	var searchPath []string

	for _, remote := range cfg.Remote {
		sanitizedName := filepath.Clean(remote)

		// http is to lazy
		if strings.HasPrefix(remote, "https://") {
			sanitizedName = sanitizedName[len("https://")-1:]
		}

		targetDir := filepath.Join(cfg.TmpDir, sanitizedName)
		if err := git.CloneOrPull(remote, targetDir); err != nil {
			return err
		}

		searchPath = append(searchPath, targetDir)
	}

	for _, p := range cfg.Local {
		stat, err := os.Stat(p)
		if err != nil {
			return err
		}

		if !stat.IsDir() {
			return fmt.Errorf("not a directory: %s", p)
		}

		searchPath = append(searchPath, p)
	}



	outDir := filepath.Join(cfg.TmpDir, "generated")
	if err := os.RemoveAll(outDir); err != nil {
		return err
	}

	projects, err := markdown.Transform(markdown.Config{
		SearchPaths: searchPath,
		OutDir:      outDir,
	})

	log.Printf("transformation complete, indexing...\n")

	sPath := filepath.Join(cfg.TmpDir, "generated-index")
	searchEngine := index.NewTruncatedSearchEngine(sPath)
	for i, prj := range projects {
		log.Printf("%0.2f%%\n", float64(i)/float64(len(projects)*100))
		if err := searchEngine.Add(prj); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}
