package repository

import (
	"backend/internal/database"
	"fmt"
)

func GetUploadedImagePaths() ([]string, error) {
	var paths []string

	queries := []string{
		`SELECT avatar_path FROM users WHERE avatar_path != '' AND avatar_path NOT LIKE '%default%';`,
		`SELECT image_path FROM posts WHERE image_path != '' AND image_path NOT LIKE '%default%';`,
		`SELECT image_path FROM comments WHERE image_path != '' AND image_path NOT LIKE '%default%';`,
		`SELECT image_path FROM group_posts WHERE image_path != '' AND image_path NOT LIKE '%default%';`,
		`SELECT image_path FROM group_comments WHERE image_path != '' AND image_path NOT LIKE '%default%';`,
	}

	for _, q := range queries {
		rows, err := database.DB.Query(q)
		if err != nil {
			fmt.Println("query error at GetUploadedImagePaths", err)
			return nil, err
		}

		for rows.Next() {
			path := ""
			if err = rows.Scan(&path); err != nil {
				fmt.Println("scan error at GetUploadedImagePaths", err)
				return nil, err
			}

			if len(path) > 0 && path[0] == '/' { // remove leading forward slash
				path = path[1:]
			}

			paths = append(paths, path)
		}
	}

	return paths, nil
}
