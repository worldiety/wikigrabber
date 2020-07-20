package index

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/worldiety/wikigrabber/internal/markdown"
	"os"
)

type Document struct {
	ProjectPath string
	SrcFile     string
	PlainText   string
	Config      *markdown.PageConfig
}

type SearchEngine struct {
	idx bleve.Index
}

func NewTruncatedSearchEngine(fname string) *SearchEngine {
	_ = os.RemoveAll(fname)
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(fname, mapping)
	if err != nil {
		panic(err)
	}

	return &SearchEngine{idx: index}

}

func (s *SearchEngine) Add(project markdown.TransformedProject) error {
	//batch := s.idx.NewBatch()

	fmt.Printf(" 00%%\n")

	lastProgress := 0
	for i, page := range project.Pages {
		doc := Document{
			ProjectPath: project.Project.Path,
			SrcFile:     page.SrcFile,
			PlainText:   page.Plain,
			Config:      page.Config,
		}

		tmp := sha256.Sum256([]byte(page.Plain))
		if err := s.idx.Index(hex.EncodeToString(tmp[:]), doc); err != nil {
			return err
		}

		progress := int(float64(i) / float64(len(project.Pages)) * 100)
		if progress != lastProgress {
			fmt.Print("\b \r\r\r\r")
			fmt.Printf(" %0.2d%%\n", progress)
			lastProgress = progress
		}
	}

	//return s.idx.Batch(batch)
	return nil
}
