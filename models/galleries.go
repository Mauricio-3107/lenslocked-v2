package models

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Image struct {
	GalleryID int
	Path      string
	Filename  string
}

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB

	// ImagesDir is used to tell the GalleryService where to store and locate images. If not set, the GalleryService will default to using the "images" directory.
	ImagesDir string
	// Extensions        []string
	// ImageContentTypes []string
}

func (service *GalleryService) Create(title string, userID int) (*Gallery, error) {
	gallery := Gallery{
		UserID: userID,
		Title:  title,
	}
	row := service.DB.QueryRow(`
	INSERT INTO galleries (title, user_id)
	VALUES ($1, $2) RETURNING id;`, gallery.Title, gallery.UserID)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (service *GalleryService) ByID(id int) (*Gallery, error) {
	gallery := Gallery{
		ID: id,
	}
	row := service.DB.QueryRow(`
	SELECT user_id, title
	FROM galleries
	WHERE id=$1;`, gallery.ID)
	err := row.Scan(&gallery.UserID, &gallery.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}
	return &gallery, nil
}

func (service *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := service.DB.Query(`
	  SELECT id, title
	  FROM galleries
	  WHERE user_id=$1`, userID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}
		err := rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	return galleries, nil
}

func (service *GalleryService) Update(gallery *Gallery) error {
	_, err := service.DB.Exec(`
	UPDATE galleries
	SET title = $2
	WHERE id = $1;`, gallery.ID, gallery.Title)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (service *GalleryService) Delete(id int) error {
	_, err := service.DB.Exec(`
	DELETE FROM galleries
	WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}
	err = os.RemoveAll(service.galleryDir(id))
	if err != nil {
		return fmt.Errorf("deleting gallery images: %w", err)
	}

	return nil
}

func (service GalleryService) galleryDir(id int) string {
	imagesDir := service.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

// func (service GalleryService) galleryDir(id int) string {
// 	imagesDir := service.ImagesDir
// 	if imagesDir == "" {
// 		imagesDir = "images"
// 	}
// 	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
// }

// Images

// func (service *GalleryService) Images(galleryID int) ([]Image, error) {
// 	globPattern := filepath.Join(service.galleryDir(galleryID), "*")
// 	allFiles, err := filepath.Glob(globPattern)
// 	if err != nil {
// 		return nil, fmt.Errorf("retrieving gallery images: %w", err)
// 	}

// 	var images []Image
// 	for _, file := range allFiles {
// 		if hasExtension(file, service.extensions()) {
// 			images = append(images, Image{
// 				GalleryID: galleryID,
// 				Path:      file,
// 				Filename:  filepath.Base(file),
// 			})
// 		}
// 	}
// 	return images, nil
// }

// func (service *GalleryService) Image(galleryID int, filename string) (Image, error) {
// 	imagePath := filepath.Join(service.galleryDir(galleryID), filename)
// 	_, err := os.Stat(imagePath)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return Image{}, ErrNotFound
// 		}
// 		return Image{}, fmt.Errorf("querying image: %w", err)
// 	}
// 	return Image{
// 		GalleryID: galleryID,
// 		Filename:  filename,
// 		Path:      imagePath,
// 	}, nil
// }

// func (service *GalleryService) CreateImage(gallleryID int, filename string, contents io.ReadSeeker) error {
// 	err := checkContentType(contents, service.imageContentTypes())
// 	if err != nil {
// 		return fmt.Errorf("creating image %v: %w", filename, err)
// 	}
// 	err = checkExtension(filename, service.extensions())
// 	if err != nil {
// 		return fmt.Errorf("creating image %v: %w", filename, err)
// 	}
// 	galleryDir := service.galleryDir(gallleryID)
// 	err = os.MkdirAll(galleryDir, 0755)
// 	if err != nil {
// 		return fmt.Errorf("creating gallery-%d images directory: %w", gallleryID, err)
// 	}
// 	imagePath := filepath.Join(galleryDir, filename)
// 	dst, err := os.Create(imagePath)
// 	if err != nil {
// 		return fmt.Errorf("creating image file: %w", err)
// 	}
// 	defer dst.Close()

// 	_, err = io.Copy(dst, contents)
// 	if err != nil {
// 		return fmt.Errorf("copying contents to image: %w", err)
// 	}
// 	return nil
// }

// func (service *GalleryService) DeleteImage(galleryID int, filename string) error {
// 	image, err := service.Image(galleryID, filename)
// 	if err != nil {
// 		return fmt.Errorf("deleting image: %w", err)
// 	}
// 	err = os.Remove(image.Path)
// 	if err != nil {
// 		return fmt.Errorf("deleting image: %w", err)
// 	}
// 	return nil
// }

// func (service *GalleryService) imageContentTypes() []string {
// 	imgCntnTypesAllowed := service.ImageContentTypes
// 	if imgCntnTypesAllowed == nil {
// 		imgCntnTypesAllowed = []string{"image/png", "image/jpeg", "image/gif"}
// 	}
// 	return imgCntnTypesAllowed
// }

// func (service *GalleryService) extensions() []string {
// 	extensionsAllowed := service.Extensions
// 	if extensionsAllowed == nil {
// 		extensionsAllowed = []string{".png", ".jpg", ".jpeg", ".gif"}
// 	}
// 	return extensionsAllowed
// }

// func hasExtension(file string, extensions []string) bool {
// 	for _, ext := range extensions {
// 		file = strings.ToLower(file)
// 		ext = strings.ToLower(ext)
// 		if filepath.Ext(file) == ext {
// 			return true
// 		}
// 	}
// 	return false
// }
