package main

import (
	"encoding/xml"
	"io/fs"
	"path/filepath"
	"sort"
)

func search(dir string) FeedbackResults {
	results := make(FeedbackResults, 0)
	filepath.WalkDir(dir,
		func(path string, info fs.DirEntry, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			content, err := checkFile(path)
			if err != nil {
				return nil
			}
			feedback := Feedback{}

			if err := xml.Unmarshal(content, &feedback); err != nil {
				return nil
			}
			results = append(results, parseFeedback(&feedback, path)...)
			return nil
		},
	)

	// sort data in descending order (based on End)
	sort.Sort(sort.Reverse(results))
	return results
}
