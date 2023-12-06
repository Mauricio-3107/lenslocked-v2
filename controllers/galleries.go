package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/Mauricio-3107/lenslocked-v2/context"
	"github.com/Mauricio-3107/lenslocked-v2/errors"
	"github.com/Mauricio-3107/lenslocked-v2/models"
	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	Templates struct {
		Show  Template
		New   Template
		Edit  Template
		Index Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}
	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")
	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	// 3. Build the data we want to pass to the Show template
	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}

	var data struct {
		ID     int
		Title  string
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title

	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}

	g.Templates.Show.Execute(w, r, data)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGalleries)
	if err != nil {
		return
	}

	// 4. execute the edit template
	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}

	var data struct {
		ID     int
		Title  string
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title

	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}

	g.Templates.Edit.Execute(w, r, data)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGalleries)
	if err != nil {
		return
	}
	// 4. parse the title and query to the DB to update
	title := r.FormValue("title")
	gallery.Title = title
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	// 5. edirect back to the edit page for now
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    int
		Title string
	}
	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())
	galleries, err := g.GalleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}

	g.Templates.Index.Execute(w, r, data)
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGalleries)
	if err != nil {
		return
	}
	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) filename(w http.ResponseWriter, r *http.Request) string {
	filename := chi.URLParam(r, "filename")
	filename = filepath.Base(filename)
	return filename
}

// 1. Define the galleryOpt type, which is really a function that takes in response writer, request, and gallery and will return an error if something isn't okay.
type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

// 2. Update the galleryByID function to accept this type as the last argument via a variadic parameter.
func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	// 1. query the id param to the url
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return nil, err
	}
	// 2. query to the DB with this id
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	// 3. Add a for loop to iterate over all of our functional options, calling each and returning if there is an error.
	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}
	return gallery, nil
}

func userMustOwnGalleries(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized  to edit this gallery", http.StatusForbidden)
		return fmt.Errorf("user does not have access to this gallery")
	}
	return nil
}

// Image controllers
func (g Galleries) Image(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the URL params
	filename := g.filename(w, r)
	galleryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	// 2. Query for the image from your Image method model
	image, err := g.GalleryService.Image(galleryID, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, image.Path)
}

func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the filename param
	filename := g.filename(w, r)

	// 2. Check the gallery ID with the galleryByID helper func and check if this user owns the gallery
	gallery, err := g.galleryByID(w, r, userMustOwnGalleries)
	if err != nil {
		return
	}
	err = g.GalleryService.DeleteImage(gallery.ID, filename)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) UploadImage(w http.ResponseWriter, r *http.Request) {
	// 1. Verify that the current user has access to this gallery
	gallery, err := g.galleryByID(w, r, userMustOwnGalleries)
	if err != nil {
		return
	}
	err = r.ParseMultipartForm(5 << 20) // 5mb
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	fileHeaders := r.MultipartForm.File["images"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		err = g.GalleryService.CreateImage(gallery.ID, fileHeader.Filename, file)
		if err != nil {
			var fileErr models.FileError
			if errors.As(err, &fileErr) {
				msg := fmt.Sprintf("%v has an invalid content type or extension. Only png, gif, and jpg files can be uploaded.", fileHeader.Filename)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}
