package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	repo *Repository
}

func NewServer(repo *Repository) *Server {
	return &Server{
		repo: repo,
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (s *Server) NewCharta(w http.ResponseWriter, r *http.Request) {

	width, err := strconv.Atoi(r.URL.Query().Get("width"))
	height, err1 := strconv.Atoi(r.URL.Query().Get("height"))
	if err != nil || err1 != nil || width <= 0 || width > 20000 || height <= 0 || height > 50000 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chart := s.repo.AddChart(width, height)
	resp := map[string]int{"id": chart.Id}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
func (s *Server) EditCharta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	width, err := strconv.Atoi(r.URL.Query().Get("width"))
	height, err1 := strconv.Atoi(r.URL.Query().Get("height"))
	x, err2 := strconv.Atoi(r.URL.Query().Get("x"))
	y, err3 := strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil || err1 != nil || err2 != nil || err3 != nil || width <= 0 || width > 20000 || height <= 0 || height > 50000 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	idInt, err5 := strconv.Atoi(id)
	if err5 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chart, err := s.repo.GetChart(idInt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	img, err := chart.AddImage(x, y, width, height)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}
	//file, fileHeader, err := r.FormFile("file")
	file, _, err := r.FormFile("file")
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	f, err := os.OpenFile(img.FileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	io.Copy(f, file)
	w.WriteHeader(http.StatusOK)
}
func (s *Server) GetCharta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	width, err := strconv.Atoi(r.URL.Query().Get("width"))
	height, err1 := strconv.Atoi(r.URL.Query().Get("height"))
	x, err2 := strconv.Atoi(r.URL.Query().Get("x"))
	y, err3 := strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil || err1 != nil || err2 != nil || err3 != nil || width <= 0 || width > 5000 || height <= 0 || height > 5000 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idInt, err1 := strconv.Atoi(id)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chart, err := s.repo.GetChart(idInt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if x+width < 0 || x > chart.Width || y+height < 0 || y > chart.Width {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	imgIds := getImages(x, y, width, height, chart)

	file := bytes.NewBuffer(nil)
	background := image.NewRGBA(image.Rect(x, y, x+width, y+height))
	black := image.NewUniform(color.RGBA{})
	draw.Draw(background, background.Bounds(), black, image.Point{}, draw.Src)

	rectOver := background.Bounds().Intersect(image.Rect(0, 0, chart.Width, chart.Heidth))
	for _, id := range imgIds {
		imgFile, _ := os.Open(chart.Images[id].FileName)
		defer imgFile.Close()
		r1 := image.Rect(chart.Images[id].X, chart.Images[id].Y, chart.Images[id].X+chart.Images[id].Width, chart.Images[id].Y+chart.Images[id].Heidth)
		r1 = r1.Bounds().Intersect(rectOver)

		img, _ := bmp.Decode(imgFile)
		draw.Draw(background, r1, img, image.Point{}, draw.Src)
	}

	bmp.Encode(file, background)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	f, _ := ioutil.ReadAll(file)
	w.Write(f)
}

func getImages(x, y, width, height int, chart *Chart) []int {
	var imgIds []int
	for id, img := range chart.Images {
		if (x <= img.X && img.X <= x+width && y <= img.Y && img.Y <= y+height) ||
			(x <= img.X+img.Width && img.X+img.Width <= x+width && y <= img.Y && img.Y <= y+height) ||
			(x <= img.X && img.X <= x+width && y <= img.Y+img.Heidth && img.Y+img.Heidth <= y+height) ||
			(x <= img.X+img.Width && img.X+img.Width <= x+width && y <= img.Y+img.Heidth && img.Y+img.Heidth <= y+height) {
			imgIds = append(imgIds, id)
		}
	}
	return imgIds
}
func (s *Server) DeleteCharta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	idInt, err1 := strconv.Atoi(id)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chart, err := s.repo.GetChart(idInt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, img := range chart.Images {
		os.Remove(img.FileName)
	}
	if err := s.repo.DeleteChart(idInt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}
