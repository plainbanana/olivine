package cmd

import (
	"log"
	"os"
	"testing"

	"github.com/gocarina/gocsv"
	"github.com/plainbanana/olivine/entities"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func TestPlotHistory1(t *testing.T) {
	// TODO : add testdata csv
	clientsFile, err := os.OpenFile("../testdata/input/test.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	csv := []*entities.Result01{}
	if err := gocsv.UnmarshalFile(clientsFile, &csv); err != nil {
		panic(err)
	}
	csv2 := append([]*entities.Result01{}, csv...)

	for i, v := range csv {
		log.Println("orig", v, v.JobName)
		if i == 10 {
			break
		}
	}

	// set X asis min max
	slice := plotHistory(csv)
	minX := slice[0].StartTime
	maxX := slice[len(slice)-1].FinishTime

	// make map from csv. key is hostname
	res := plotHistories(csv2, FEachHost)
	log.Println("Test len", len(csv), len(res))

	rows := len(res)
	cols := 1

	plots := make([][]*plot.Plot, rows)

	j := 0
	for h, v := range res {
		plots[j] = make([]*plot.Plot, cols)
		for i := 0; i < cols; i++ {
			p := drawAsLinesPoints(v)

			p.X.Min = tof64(minX)
			p.X.Max = tof64(maxX)
			p.Title.Text = h

			p.Add(plotter.NewGrid())

			plots[j][i] = p
		}
		j++
	}
	for h, v := range res {
		log.Println("cat", h, len(v))

	}

	img := vgimg.New(15*vg.Inch, 20*vg.Inch)
	dc := draw.New(img)

	tt := draw.Tiles{
		Rows: rows,
		Cols: cols,
	}

	canvases := plot.Align(plots, tt, dc)
	for j := 0; j < rows; j++ {
		for i := 0; i < cols; i++ {
			if plots[j][i] != nil {
				plots[j][i].Draw(canvases[j][i])
			}
		}
	}

	w, err := os.Create("../testdata/plot/aligned_taskline.png")
	if err != nil {
		panic(err)
	}

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}
}
