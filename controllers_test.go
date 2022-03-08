package main

import (
	"reflect"
	"testing"
)

func TestGetImages(t *testing.T) {
	chart := Chart{
		Id:     1,
		Width:  50,
		Heidth: 50,
		Images: []Image{
			{Heidth: 10, Width: 10, X: -5, Y: -5, Id: 1},
			{Heidth: 20, Width: 20, X: 35, Y: 35, Id: 2},
			{Heidth: 10, Width: 30, X: 20, Y: 10, Id: 3},
			{Heidth: 5, Width: 5, X: 10, Y: 30, Id: 4}}}
	type data struct {
		x      int
		y      int
		width  int
		height int
	}
	tests := []struct {
		name string
		data data
		want []int
	}{
		{name: "1", data: data{x: 0, y: 0, width: 50, height: 50}, want: []int{0, 1, 2, 3}},
		{name: "2", data: data{x: 10, y: 0, width: 50, height: 5}, want: nil},
		{name: "3", data: data{x: 10, y: 10, width: 20, height: 25}, want: []int{2, 3}},
		{name: "4", data: data{x: -10, y: -10, width: 25, height: 50}, want: []int{0, 3}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			x := test.data.x
			y := test.data.y
			width := test.data.width
			height := test.data.height
			got := getImages(x, y, width, height, &chart)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("%s: got %d, want %d", test.name, got, test.want)
			}
		})
	}
}
