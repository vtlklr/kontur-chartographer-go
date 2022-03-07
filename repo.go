package main

import "fmt"

type Repository struct {
	charts  map[int]*Chart
	counter int
}

type Chart struct {
	Heidth  int
	Width   int
	Id      int
	Images  []Image
	counter int
}

type Image struct {
	X        int
	Y        int
	Heidth   int
	Width    int
	Id       int
	FileName string
}

func New() *Repository {
	return &Repository{charts: make(map[int]*Chart, 0), counter: 0}
}

func (r *Repository) AddChart(width, heidth int) *Chart {
	c := &Chart{Heidth: heidth, Width: width, Id: r.nextCount()}
	r.charts[c.Id] = c
	return c
}

func (r *Repository) GetChart(id int) (*Chart, error) {
	if c, ok := r.charts[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("нет id")
}

func (c *Chart) AddImage(x, y, width, heigth int) (Image, error) {
	if x > c.Width || y > c.Heidth || x+width < 0 || y+heigth < 0 {
		return Image{}, fmt.Errorf("не правильные координиаты")
	}
	img := Image{X: x, Y: y, Width: width, Heidth: heigth, Id: c.nextCount()}
	img.FileName = fmt.Sprintf("chart%dimg%d.bmp", c.Id, img.Id)
	c.Images = append(c.Images, img)

	return img, nil
}

func (r *Repository) nextCount() int {
	r.counter = r.counter + 1
	return r.counter
}

func (c *Chart) nextCount() int {
	c.counter = c.counter + 1
	return c.counter
}
