package main

import (
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
	//fmt.Println("new charta")

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
	fmt.Println(`id := `, id)
	width, err := strconv.Atoi(r.URL.Query().Get("width"))
	height, err1 := strconv.Atoi(r.URL.Query().Get("height"))
	x, err2 := strconv.Atoi(r.URL.Query().Get("x"))
	y, err3 := strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil || err1 != nil || err2 != nil || err3 != nil || width <= 0 || width > 20000 || height <= 0 || height > 50000 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	idInt, _ := strconv.Atoi(id)
	chart, err := s.repo.GetChart(idInt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}
	img, err := chart.AddImage(x, y, width, height)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(img.Id)
	file, fileHeader, err := r.FormFile("file")
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	f, err := os.OpenFile(img.FileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	io.Copy(f, file)
	fmt.Println(fileHeader)

}
func (s *Server) GetCharta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//fmt.Println(`get id := `, id)
	width, err := strconv.Atoi(r.URL.Query().Get("width"))
	height, err1 := strconv.Atoi(r.URL.Query().Get("height"))
	x, err2 := strconv.Atoi(r.URL.Query().Get("x"))
	y, err3 := strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil || err1 != nil || err2 != nil || err3 != nil || width <= 0 || width > 5000 || height <= 0 || height > 5000 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	idInt, _ := strconv.Atoi(id)
	chart, err := s.repo.GetChart(idInt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}
	var imgIds []int
	for id, img := range chart.Images {
		if (x <= img.X && img.X <= x+width && y <= img.Y && img.Y <= y+height) ||
			(x <= img.X+img.Width && img.X+img.Width <= x+width && y <= img.Y && img.Y <= y+height) ||
			(x <= img.X && img.X <= x+width && y <= img.Y+img.Heidth && img.Y+img.Heidth <= y+height) ||
			(x <= img.X+img.Width && img.X+img.Width <= x+width && y <= img.Y+img.Heidth && img.Y+img.Heidth <= y+height) {
			imgIds = append(imgIds, id)
		}
	}
	//fmt.Println(imgIds)
	fileName := fmt.Sprintf("chart%dx%dy%dwidth%dheight%d.bmp", chart.Id, x, y, width, height)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("файл не создан")
		return
	}
	defer os.Remove(fileName)
	defer file.Close()
	background := image.NewRGBA(image.Rect(x, y, x+width, y+height))
	black := image.NewUniform(color.RGBA{})
	draw.Draw(background, background.Bounds(), black, image.Point{}, draw.Src)

	rectOver := background.Bounds().Intersect(image.Rect(0, 0, chart.Width, chart.Heidth))
	//fmt.Println(rectOver.Bounds())
	for _, id := range imgIds {
		imgFile, _ := os.Open(chart.Images[id].FileName)
		defer imgFile.Close()
		img, _ := bmp.Decode(imgFile)
		r1 := image.Rect(chart.Images[id].X, chart.Images[id].Y, chart.Images[id].X+chart.Images[id].Width, chart.Images[id].Y+chart.Images[id].Heidth)
		r1 = r1.Bounds().Intersect(rectOver)
		draw.Draw(background, r1, img, image.Point{}, draw.Src)
	}

	bmp.Encode(file, background)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	w.Write(fileBytes)
}
func (s *Server) DeleteCharta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	idInt, _ := strconv.Atoi(id)
	chart, err := s.repo.GetChart(idInt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	for _, img := range chart.Images {
		err1 := os.Remove(img.FileName)
		if err1 != nil || !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	w.WriteHeader(http.StatusOK)

}
