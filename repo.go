package main

import (
	"fmt"
	"sync/atomic"
)

type Repository struct {
	charts  map[int]*Chart
	counter uint32
}

type Chart struct {
	Heidth  int
	Width   int
	Id      uint32
	Images  []Image
	counter uint32
}

type Image struct {
	X        int
	Y        int
	Heidth   int
	Width    int
	FileName string
}

func New() *Repository {
	return &Repository{charts: make(map[int]*Chart, 0), counter: 0}
}

func (r *Repository) AddChart(width, heidth int) *Chart {
	c := &Chart{Heidth: heidth, Width: width, Id: r.nextCount()}
	r.charts[int(c.Id)] = c
	return c
}

func (r *Repository) GetChart(id int) (*Chart, error) {
	if c, ok := r.charts[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("нет id")
}
func (r *Repository) DeleteChart(id int) error {
	if _, ok := r.charts[id]; ok {
		delete(r.charts, id)
		return nil
	}
	return fmt.Errorf("нет id")
}

func (c *Chart) AddImage(x, y, width, heigth int) (Image, error) {
	if x > c.Width || y > c.Heidth || x+width < 0 || y+heigth < 0 {
		return Image{}, fmt.Errorf("не правильные координиаты")
	}
	c.nextCount()
	img := Image{X: x, Y: y, Width: width, Heidth: heigth}
	img.FileName = fmt.Sprintf("chart%dimg%d.bmp", c.Id, c.counter)
	c.Images = append(c.Images, img)

	return img, nil
}

func (r *Repository) nextCount() uint32 {
	//r.counter = r.counter + 1
	atomic.AddUint32(&r.counter, 1)
	return r.counter
}

func (c *Chart) nextCount() uint32 {
	//c.counter = c.counter + 1
	atomic.AddUint32(&c.counter, 1)
	return c.counter
}
