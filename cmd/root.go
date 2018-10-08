// Copyright © 2018 Chris Tava <chris1tava@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	chart "github.com/wcharczuk/go-chart"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "github-teamwork",
	Short: "a set of commands to foster collaboration on github.com",
	Long:  fmt.Sprintf(`github-teamwork - a set of commands to get answers to your questions about github.com Version: %s Author: Chris Tava <chris1tava@gmail.com>`, VERSION),
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	defaultThreads := runtime.NumCPU()
	if defaultThreads > 2 {
		defaultThreads = 2
	}
	RootCmd.PersistentFlags().IntP("threads", "t", defaultThreads, "number of CPUs. (default value: 1 for single-CPU PC, 2 for others)")
}

func getFlagString(cmd *cobra.Command, flag string) string {
	value, err := cmd.Flags().GetString(flag)
	checkError(err)
	return value
}

func getFlagBool(cmd *cobra.Command, flag string) bool {
	value, err := cmd.Flags().GetBool(flag)
	checkError(err)
	return value
}

func getFlagInt(cmd *cobra.Command, flag string) int {
	value, err := cmd.Flags().GetInt(flag)
	checkError(err)
	return value
}

const numberOfCharactersInDate = 10

func daysIn(m time.Month, y int) int {
	return time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func writeDataSetToFile(fileName string, data []byte) error {
	err := ioutil.WriteFile(fileName, data, os.ModePerm)
	if err != nil {
		return errors.New("Could not write output file")
	}
	return nil
}

func getDataSetFromFile(fileName string) (bytes.Buffer, error) {
	var b bytes.Buffer
	csvBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return b, errors.New("Could not read input file")
	}

	b.Write(csvBytes)
	return b, nil
}

func getTimeSeriesDataForTheMonth(startYear, endYear int, startMonth, endMonth time.Month) ([]time.Time, error) {

	//startTime := time.Date(startYear, startMonth, 1, 0, 0, 0, 0, time.UTC)
	//endTime := time.Date(endYear, endMonth, 1, 0, 0, 0, 0, time.UTC) //not supported
	var timeSeries []time.Time
	for i := 1; i <= daysIn(startMonth, startYear); i++ {
		day := time.Date(startYear, startMonth, i, 0, 0, 0, 0, time.UTC)
		timeSeries = append(timeSeries, day)
	}
	return timeSeries, nil
}

func getCountDataPerDay(startYear int, startMonth time.Month, dataSetFileName string) ([]float64, error) {

	data, err := getDataSetFromFile(dataSetFileName)
	if err != nil {
		return nil, errors.New("Could not get data from file")
	}
	reader := csv.NewReader(bytes.NewReader(data.Bytes()))
	reader.FieldsPerRecord = 1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("Could not read in data records.")
	}

	countsMap := make(map[int]float64)
	layout := "2006-01-02"
	for i := 1; i <= daysIn(startMonth, startYear); i++ {
		day := time.Date(startYear, startMonth, i, 0, 0, 0, 0, time.UTC)
		countsMap[day.Day()] = float64(0)
		for _, each := range records {
			t, err := time.Parse(layout, each[0][0:numberOfCharactersInDate])
			if err != nil {
				return nil, errors.New("Could not parse timestamps")
			}
			if day.Equal(t) {
				dayCount := countsMap[day.Day()]
				countsMap[day.Day()] = dayCount + 1
			}
		}
	}

	var counts []float64
	for _, v := range countsMap {
		counts = append(counts, v)
	}
	return counts, nil
}

func drawChart(startYear, endYear int, startMonth, endMonth time.Month, legend, inputFileName, outputfileName string) error {

	timeSeries, tserr := getTimeSeriesDataForTheMonth(startYear, endYear, startMonth, endMonth)
	if tserr != nil {
		return tserr
	}
	counts, cerr := getCountDataPerDay(startYear, startMonth, inputFileName)
	if cerr != nil {
		return cerr
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{Show: true},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    legend,
				XValues: timeSeries,
				YValues: counts,
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	f, err := os.Create(outputfileName)
	if err != nil {
		return err
	}
	if _, err := f.Write(buffer.Bytes()); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func drawChartWithFourLines(startYear, endYear int, startMonth, endMonth time.Month, legend1, legend2, legend3, legend4, input1FileName, input2FileName, input3FileName, input4FileName, outputfileName string) error {

	timeSeries, tserr := getTimeSeriesDataForTheMonth(startYear, endYear, startMonth, endMonth)
	if tserr != nil {
		return tserr
	}
	counts1, cerr := getCountDataPerDay(startYear, startMonth, input1FileName)
	if cerr != nil {
		return cerr
	}
	counts2, cerr := getCountDataPerDay(startYear, startMonth, input2FileName)
	if cerr != nil {
		return cerr
	}
	counts3, cerr := getCountDataPerDay(startYear, startMonth, input3FileName)
	if cerr != nil {
		return cerr
	}
	counts4, cerr := getCountDataPerDay(startYear, startMonth, input4FileName)
	if cerr != nil {
		return cerr
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{Show: true},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{Name: legend1, XValues: timeSeries, YValues: counts1},
			chart.TimeSeries{Name: legend2, XValues: timeSeries, YValues: counts2},
			chart.TimeSeries{Name: legend3, XValues: timeSeries, YValues: counts3},
			chart.TimeSeries{Name: legend4, XValues: timeSeries, YValues: counts4},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	f, err := os.Create(outputfileName)
	if err != nil {
		return err
	}
	if _, err := f.Write(buffer.Bytes()); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
