package main

import (
	"testing"
	"time"

	"github.com/robin-samuel/sensor"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

func TestSimulator(t *testing.T) {
	start := time.Now()
	end := start.Add(10 * time.Second)
	sim := sensor.NewSimulator(start, end, 0.5)
	var positions []sensor.Position
	var orientations []sensor.Orientation

	// pull every 66 milliseconds
	for t := start; t.Before(end); t = t.Add(time.Millisecond * 66) {
		positions = append(positions, sim.Position(t))
		orientations = append(orientations, sim.Orientation(t))
	}

	p := plot.New()
	p.Title.Text = "Position"

	var x, y, z plotter.XYs
	for i, pos := range positions {
		x = append(x, plotter.XY{X: float64(i), Y: pos.Values[0]})
		y = append(y, plotter.XY{X: float64(i), Y: pos.Values[1]})
		z = append(z, plotter.XY{X: float64(i), Y: pos.Values[2]})
	}

	lineX, err := plotter.NewLine(x)
	if err != nil {
		t.Fatal(err)
	}
	lineY, err := plotter.NewLine(y)
	if err != nil {
		t.Fatal(err)
	}
	lineZ, err := plotter.NewLine(z)
	if err != nil {
		t.Fatal(err)
	}
	lineX.Color = plotutil.Color(0)
	lineY.Color = plotutil.Color(1)
	lineZ.Color = plotutil.Color(2)

	p.Add(lineX, lineY, lineZ)

	// Save the plot to a PNG file.
	if err := p.Save(1200, 800, "position.png"); err != nil {
		t.Fatal(err)
	}

	p = plot.New()
	p.Title.Text = "Orientation"

	var pitch, roll, yaw plotter.XYs
	for i, ori := range orientations {
		pitch = append(pitch, plotter.XY{X: float64(i), Y: ori.Values[0]})
		roll = append(roll, plotter.XY{X: float64(i), Y: ori.Values[1]})
		yaw = append(yaw, plotter.XY{X: float64(i), Y: ori.Values[2]})
	}

	linePitch, err := plotter.NewLine(pitch)
	if err != nil {
		t.Fatal(err)
	}
	lineRoll, err := plotter.NewLine(roll)
	if err != nil {
		t.Fatal(err)
	}
	lineYaw, err := plotter.NewLine(yaw)
	if err != nil {
		t.Fatal(err)
	}
	linePitch.Color = plotutil.Color(0)
	lineRoll.Color = plotutil.Color(1)
	lineYaw.Color = plotutil.Color(2)

	p.Add(linePitch, lineRoll, lineYaw)

	// Save the plot to a PNG file.
	if err := p.Save(1200, 800, "orientation.png"); err != nil {
		t.Fatal(err)
	}
}

func TestSensor(t *testing.T) {
	start := time.Now()
	end := start.Add(10 * time.Second)
	s := sensor.NewManager(start, end, 0.3)
	var aEvents []sensor.Event
	var gEvents []sensor.Event
	var mEvents []sensor.Event

	// pull every 66 milliseconds
	for t := start; t.Before(end); t = t.Add(time.Millisecond * 66) {
		aEvents = append(aEvents, s.Get(sensor.Accelerometer, t))
		gEvents = append(gEvents, s.Get(sensor.Gyroscope, t))
		// mEvents = append(mEvents, s.Get(sensor.Magnetometer, t))
	}

	p := plot.New()
	p.Title.Text = "Accelerometer"

	var x, y, z plotter.XYs
	for i, event := range aEvents {
		x = append(x, plotter.XY{X: float64(i), Y: event.Data[0]})
		y = append(y, plotter.XY{X: float64(i), Y: event.Data[1]})
		z = append(z, plotter.XY{X: float64(i), Y: event.Data[2]})
	}

	lineX, err := plotter.NewLine(x)
	if err != nil {
		t.Fatal(err)
	}
	lineY, err := plotter.NewLine(y)
	if err != nil {
		t.Fatal(err)
	}
	lineZ, err := plotter.NewLine(z)
	if err != nil {
		t.Fatal(err)
	}
	lineX.Color = plotutil.Color(0)
	lineY.Color = plotutil.Color(1)
	lineZ.Color = plotutil.Color(2)

	p.Add(lineX, lineY, lineZ)

	// Save the plot to a PNG file.
	if err := p.Save(1200, 800, "accelerometer.png"); err != nil {
		t.Fatal(err)
	}

	p = plot.New()
	p.Title.Text = "Gyroscope"

	var pitch, roll, yaw plotter.XYs
	for i, event := range gEvents {
		pitch = append(pitch, plotter.XY{X: float64(i), Y: event.Data[0]})
		roll = append(roll, plotter.XY{X: float64(i), Y: event.Data[1]})
		yaw = append(yaw, plotter.XY{X: float64(i), Y: event.Data[2]})
	}

	linePitch, err := plotter.NewLine(pitch)
	if err != nil {
		t.Fatal(err)
	}
	lineRoll, err := plotter.NewLine(roll)
	if err != nil {
		t.Fatal(err)
	}
	lineYaw, err := plotter.NewLine(yaw)
	if err != nil {
		t.Fatal(err)
	}
	linePitch.Color = plotutil.Color(0)
	lineRoll.Color = plotutil.Color(1)
	lineYaw.Color = plotutil.Color(2)

	p.Add(linePitch, lineRoll, lineYaw)

	// Save the plot to a PNG file.
	if err := p.Save(1200, 800, "gyroscope.png"); err != nil {
		t.Fatal(err)
	}

	p = plot.New()
	p.Title.Text = "Magnetometer"

	x, y, z = plotter.XYs{}, plotter.XYs{}, plotter.XYs{}
	for i, event := range mEvents {
		x = append(x, plotter.XY{X: float64(i), Y: event.Data[0]})
		y = append(y, plotter.XY{X: float64(i), Y: event.Data[1]})
		z = append(z, plotter.XY{X: float64(i), Y: event.Data[2]})
	}

	lineX, err = plotter.NewLine(x)
	if err != nil {
		t.Fatal(err)
	}
	lineY, err = plotter.NewLine(y)
	if err != nil {
		t.Fatal(err)
	}
	lineZ, err = plotter.NewLine(z)
	if err != nil {
		t.Fatal(err)
	}
	lineX.Color = plotutil.Color(0)
	lineY.Color = plotutil.Color(1)
	lineZ.Color = plotutil.Color(2)

	p.Add(lineX, lineY, lineZ)

	// Save the plot to a PNG file.
	if err := p.Save(1200, 800, "magnetometer.png"); err != nil {
		t.Fatal(err)
	}
}
