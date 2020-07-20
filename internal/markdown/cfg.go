package markdown

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Project struct {
	Name        string
	Path        string
	Markdown    []string
	Attachments []string
}

type Config struct {
	SearchPaths []string
	OutDir      string
}

func (c *Config) Collect() ([]Project, error) {
	var res []Project
	for _, root := range c.SearchPaths {
		p := Project{
			Path: root,
			Name: filepath.Base(root),
		}

		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.HasPrefix(info.Name(), ".") {
				if info.IsDir() {
					return filepath.SkipDir
				}

				return nil
			}

			if info.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(root, path)
			if err != nil {
				panic("cannot happen")
			}

			if strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
				p.Markdown = append(p.Markdown, relPath)
			} else {
				p.Attachments = append(p.Attachments, relPath)
			}

			return nil
		})

		if err != nil {
			return res, err
		}

		res = append(res, p)
	}

	return res, nil
}

type PageConfig struct {
	// Alternate title for the page
	Title string `yaml:"title"`

	// Menu is slash separated, like my/cool/menu
	Menu string `yaml:"menu"`

	// ManagedBy names someone who is in charge of the page
	ManagedBy string `yaml:"managedBy"`

	// whatever that means
	Date time.Time `yaml:"date"`

	// Tags for an alternate hierarchyless navigation
	Tags []string `yaml:"tags"`
}
