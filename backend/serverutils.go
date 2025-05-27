package main

import (
	"backend/internal/repository"
	"backend/internal/utils"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

func getMigrationImagePaths() []string {
	paths := []string{}
	sqlFiles := []string{"migrations/000002_insert_data.up.sql"}

	// find image path strings inside sql files
	re := regexp.MustCompile(`'/?data/uploads/[^']+'`) // from '/data/uploads/ or 'data/uploads/ to next '
	for _, file := range sqlFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading", file, ":", err)
			continue
		}
		matches := re.FindAllString(string(content), -1)
		for _, match := range matches {
			// Remove the surrounding single quotes
			cleanPath := match[1 : len(match)-1]
			paths = append(paths, cleanPath)
		}
	}

	return paths
}

func deleteUnusedImages() {
	imgPaths, err := repository.GetUploadedImagePaths() // image paths in database
	if err != nil {
		fmt.Println(err)
		return
	}
	imgPaths = append(imgPaths, getMigrationImagePaths()...) // image paths in migration files

	// map paths for fast lookups
	usedPaths := map[string]bool{}
	for _, path := range imgPaths {
		usedPaths[filepath.Base(path)] = true // remove folder from path; slashes flip in different OSs, filenames are unique enough
		//fmt.Println(filepath.Base(path))
	}

	directories := []string{
		"data/uploads/avatars",
		"data/uploads/comments",
		"data/uploads/posts",
	}

	for _, dir := range directories {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			// delete if file in directory is not in used paths and is an image
			if !usedPaths[filepath.Base(path)] && utils.IsAllowedImageExtension(filepath.Ext(path)) {
				fmt.Println("Deleting unused file:", path)

				if err := os.Remove(path); err != nil {
					return fmt.Errorf("failed to delete %s: %w", path, err)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking the path %s: %v\n", dir, err)
		}
	}
}
