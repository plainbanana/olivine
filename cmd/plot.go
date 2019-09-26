package cmd

import (
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/plainbanana/olivine/entities"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

var plotCmd = &cobra.Command{
	Use:              "plot",
	Short:            "Visualize data.",
	Long:             `Plot data which read from a csv file.`,
	PersistentPreRun: setLogMinLevel,
	Run:              plotRun,
}

func plotRun(cmd *cobra.Command, args []string) {
	file, err := os.OpenFile(config.FileInput, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("alert: something happen when open a file.")
		panic(err)
	}
	defer file.Close()

	csv := []*entities.Result01{}
	if err := gocsv.UnmarshalFile(file, &csv); err != nil {
		panic(err)
	}
	csv2 := append([]*entities.Result01{}, csv...)

	for i, v := range csv {
		log.Println("trace: orig", v, v.JobName)
		if i == 10 {
			break
		}
	}

	// set X asis min max
	slice := plotHistory(csv)
	minX := slice[0].StartTime
	maxX := slice[len(slice)-1].FinishTime

	// make map from csv. key is hostname
	log.Println("warn: TODO Full duplicate check")
	res := plotHistories(csv2, FEachHost)

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
		log.Println("trace: histories map", h, len(v))
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

	w, err := os.Create(config.FileOutput)
	if err != nil {
		panic(err)
	}

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		log.Println("alert: can not save image.")
		panic(err)
	}
	log.Println("info: success to save image.")
}

func plotHistories(inputOrig []*entities.Result01, flag ...int) map[string][]*entities.Result01 {
	input := append([]*entities.Result01{}, inputOrig...)

	m := make(map[string][]*entities.Result01)
	// sort by time
	sort.SliceStable(input, func(i, j int) bool {
		return input[i].StartTime < input[j].StartTime
	})

	if sort.SliceIsSorted(input, func(i, j int) bool {
		return input[i].StartTime < input[j].StartTime
	}) {
		log.Println("debug: sorted!!!")
	} else {
		log.Println("debug: NOTsorted!!!")
	}

	for _, v := range flag {
		// plot data for each host
		if v&FEachHost != 0 {
			tmp := make(map[string][]*entities.Result01)
			for _, v := range input {
				tmp[v.Hostname] = append(tmp[v.Hostname], v)
			}
			log.Println("trace: map1", input)
			log.Println("trace: map", len(tmp))
			for host, val := range tmp {
				m[host] = plotHistory(val)
			}
			return m
		}
	}
	return m
}

func plotHistory(input []*entities.Result01) []*entities.Result01 {
	// 開始時間昇順にソート
	sort.SliceStable(input, func(i, j int) bool {
		return input[i].StartTime < input[j].StartTime
	})

	if sort.SliceIsSorted(input, func(i, j int) bool {
		return input[i].StartTime < input[j].StartTime
	}) {
		log.Println("debug: sorted!!!")
	} else {
		log.Println("debug: NOTsorted!!!")
	}

	for i, v := range input {
		log.Println("trace: sorted", v, v.JobName)
		if i == 10 {
			break
		}
	}

	i := 0
	for {
		if len(input) == 0 {
			log.Println("warn:", ErrInputRange)
		}
		// 最後の要素
		if i == len(input)-1 {
			break
		}

		// 2要素目以降
		if i > 0 && i != len(input)-1 {
			// 隣接するtaskの実行時間が重複しているとき
			if input[i-1].FinishTime >= input[i].StartTime {
				input[i].StartTime = input[i-1].StartTime
				if input[i].FinishTime >= input[i-1].FinishTime {
				} else {
					input[i].FinishTime = input[i-1].FinishTime
				}
				input = unsetSliceContent(input, i-1)
				// Sliceの要素が一つ減るのでデクリメント
				i--
			}
		}

		i++
	}

	return input
}

func unsetSliceContent(s []*entities.Result01, i int) []*entities.Result01 {
	if i >= len(s) {
		return s
	}
	return append(s[:i], s[i+1:]...)
}

func drawAsLinesPoints(input []*entities.Result01) *plot.Plot {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}

	p.Title.Text = "History"
	p.X.Label.Text = "time(ns)"
	p.HideY()

	c := 0
	lines := []interface{}{}
	for _, v := range input {
		var dataArr plotter.XYs

		var d1, d2 plotter.XY
		d1.X = tof64(v.StartTime)
		d1.Y = intTof64(0)
		dataArr = append(dataArr, d1)

		d2.X = tof64(v.FinishTime)
		d2.Y = intTof64(0)
		dataArr = append(dataArr, d2)

		lines = append(lines, v.Hostname)
		lines = append(lines, dataArr)
		c++
	}

	err = plotutil.AddLinePoints(p, lines...)
	if err != nil {
		log.Panic(err)
	}

	return p
}

func tof64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func intTof64(i int) float64 {
	s := strconv.Itoa(i)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
