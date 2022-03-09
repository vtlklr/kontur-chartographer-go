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
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	repo *Repository
}

func NewServer(repo *Repository) *Server {
	return &Server{
		repo: repo,
	}
}

func (s *Server) NewCharta(w http.ResponseWriter, r *http.Request) {

	width, err := strconv.Atoi(r.URL.Query().Get("width"))
	if err != nil {
		fmt.Println("не корректно указана ширина " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	height, err1 := strconv.Atoi(r.URL.Query().Get("height"))
	if err1 != nil {
		fmt.Println("не корректно указана высота " + err1.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if width <= 0 || width > 20000 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("не корректно указана ширина")
		return
	}
	if height <= 0 || height > 50000 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("не корректно указана высота")
		return
	}
	chart := s.repo.AddChart(width, height)
	resp := map[string]uint32{"id": chart.Id}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err2 := json.NewEncoder(w).Encode(resp)
	if err2 != nil {
		fmt.Println(err2.Error())
		return
	}
}
func (s *Server) EditCharta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	width, err := strconv.Atoi(r.URL.Query().Get("width"))
	if err != nil {
		fmt.Println("не корректно указана ширина " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	height, err := strconv.Atoi(r.URL.Query().Get("height"))
	if err != nil {
		fmt.Println("не корректно указана высота " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	x, err := strconv.Atoi(r.URL.Query().Get("x"))
	if err != nil {
		fmt.Println("не корректно указана координата x " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil {
		fmt.Println("не корректно указана координата y " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if width <= 0 || width > 20000 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("не корректно указана ширина")
		return
	}
	if height <= 0 || height > 50000 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("не корректно указана высота")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("не корректно указан id " + err.Error())
		return
	}
	chart, err := s.repo.GetChart(idInt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("не найден id")
		return
	}
	img, err := chart.AddImage(x, y, width, height)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("ошибка открытия файла " + err.Error())
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("ошибка закрытия файла " + err.Error())
			return
		}
	}(file)
	f, err := os.OpenFile(img.FileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("ошибка открытия файла " + err.Error())
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("ошибка закрытия файла" + img.FileName)
		}
	}(f)
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
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
	if err != nil {
		fmt.Println("не корректно указана ширина " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	height, err := strconv.Atoi(r.URL.Query().Get("height"))
	if err != nil {
		fmt.Println("не корректно указана высота " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	x, err := strconv.Atoi(r.URL.Query().Get("x"))
	if err != nil {
		fmt.Println("не корректно указана координата x " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil {
		fmt.Println("не корректно указана координата y " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if width <= 0 || width > 5000 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("не корректно указана ширина")
		return
	}
	if height <= 0 || height > 5000 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("не корректно указана высота")
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("не корректно указан id " + err.Error())
		return
	}
	chart, err := s.repo.GetChart(idInt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("id: %d  не найден\n", idInt)
		return
	}
	if x+width < 0 || x > chart.Width || y+height < 0 || y > chart.Width {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("запрашиваемый фрагмент вне изображения")
		return
	}
	imgIds := getImages(x, y, width, height, chart)

	background := image.NewRGBA(image.Rect(x, y, x+width, y+height))
	black := image.NewUniform(color.RGBA{})
	draw.Draw(background, background.Bounds(), black, image.Point{}, draw.Src)

	rectOver := background.Bounds().Intersect(image.Rect(0, 0, chart.Width, chart.Heidth))
	for _, id := range imgIds {
		imgFile, err := os.Open(chart.Images[id].FileName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("ошибка открытия файла " + err.Error())
			return
		}
		r1 := image.Rect(chart.Images[id].X, chart.Images[id].Y, chart.Images[id].X+chart.Images[id].Width, chart.Images[id].Y+chart.Images[id].Heidth)
		r1 = r1.Bounds().Intersect(rectOver)

		img1, err := bmp.Decode(imgFile)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("ошибка декодирования bmp " + err.Error())
			return
		}
		draw.Draw(background, r1, img1, image.Point{}, draw.Src)
		err = imgFile.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("ошибка закрытия файла " + err.Error())
			return
		}
	}
	buf := bytes.NewBuffer(nil)
	err = bmp.Encode(buf, background)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("ошибка кодирования bmp " + err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	f, err := ioutil.ReadAll(buf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("ошибка чтения " + err.Error())
		return
	}
	_, err = w.Write(f)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("ошибка записи " + err.Error())
		return
	}
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
		fmt.Println("не корректно указан id")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("не корректно указан id " + err.Error())
		return
	}
	chart, err := s.repo.GetChart(idInt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("id не найден")
		return
	}

	for _, img := range chart.Images {
		err := os.Remove(img.FileName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("ошибка удаления файла " + err.Error())
			return
		}
	}
	if err := s.repo.DeleteChart(idInt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("ошибка удаления изображения " + err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)

}
