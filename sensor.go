package sensor

import (
	"math"
	"math/rand"
	"time"

	"github.com/robin-samuel/spline"
)

// Type represents a sensor type.
type Type int

var sensorNames = map[Type]string{
	Accelerometer: "Accelerometer",
	Gyroscope:     "Gyroscope",
	Magnetometer:  "Magnetometer",
}

// String returns the string representation of the sensor type.
func (t Type) String() string {
	if n, ok := sensorNames[t]; ok {
		return n
	}
	return "Unknown sensor"
}

const (
	Accelerometer Type = iota
	Gyroscope
	Magnetometer
)

type Event struct {
	// Sensor is the type of the sensor the event is coming from.
	Sensor Type

	// Timestamp is a event time in nanoseconds.
	Timestamp time.Time

	// Data is the event data.
	//
	// If the event source is Accelerometer,
	//  - Data[0]: acceleration force in x axis in m/s^2
	//  - Data[1]: acceleration force in y axis in m/s^2
	//  - Data[2]: acceleration force in z axis in m/s^2
	//
	// If the event source is Gyroscope,
	//  - Data[0]: rate of rotation around the x axis in rad/s
	//  - Data[1]: rate of rotation around the y axis in rad/s
	//  - Data[2]: rate of rotation around the z axis in rad/s
	//
	// If the event source is Magnetometer,
	//  - Data[0]: force of gravity along the x axis in m/s^2
	//  - Data[1]: force of gravity along the y axis in m/s^2
	//  - Data[2]: force of gravity along the z axis in m/s^2
	//
	Data []float64
}

type Manager struct {
	start time.Time
	end   time.Time
	sim   *Simulator
}

func NewManager(start time.Time, end time.Time, activity float64) *Manager {
	return &Manager{
		start: start,
		end:   end,
		sim:   NewSimulator(start.Add(-time.Second), end.Add(time.Second), activity),
	}
}

func (s *Manager) Start() time.Time {
	return s.start
}

func (s *Manager) End() time.Time {
	return s.end
}

func (s *Manager) Orientation(t time.Time) Orientation {
	return s.sim.Orientation(t)
}

func (s *Manager) Position(t time.Time) Position {
	return s.sim.Position(t)
}

func (s *Manager) Get(st Type, t time.Time) Event {
	pos0 := s.sim.Position(t.Add(-time.Millisecond * 66))
	pos1 := s.sim.Position(t)
	ori0 := s.sim.Orientation(t.Add(-time.Millisecond * 66))
	ori1 := s.sim.Orientation(t)

	switch st {
	case Accelerometer:
		// calculate position delta and convert millimeters to meters
		positionD := []float64{
			(pos1.Values[0] - pos0.Values[0]) / 1000,
			(pos1.Values[1] - pos0.Values[1]) / 1000,
			(pos1.Values[2] - pos0.Values[2]) / 1000,
		}
		// calculate acceleration
		timeD := pos1.Timestamp.Sub(pos0.Timestamp).Seconds()
		acceleration := []float64{
			positionD[0] / timeD / timeD,
			positionD[1] / timeD / timeD,
			positionD[2] / timeD / timeD,
		}
		return Event{
			Sensor:    Accelerometer,
			Timestamp: pos1.Timestamp,
			Data:      acceleration,
		}
	case Gyroscope:
		// calculate orientation delta
		oientationD := []float64{
			ori1.Values[0] - ori0.Values[0],
			ori1.Values[1] - ori0.Values[1],
			ori1.Values[2] - ori0.Values[2],
		}
		// calculate angular velocity
		timeD := ori1.Timestamp.Sub(ori0.Timestamp).Seconds()
		angularVelocity := []float64{
			oientationD[0] / timeD,
			oientationD[1] / timeD,
			oientationD[2] / timeD,
		}
		return Event{
			Sensor:    Gyroscope,
			Timestamp: ori1.Timestamp,
			Data:      angularVelocity,
		}
	default:
		return Event{}
	}
}

// Position is a 3D vector representing the position of the device.
type Position struct {
	// Timestamp is a event time in nanoseconds.
	Timestamp time.Time

	// Values is the position data.
	//  - Values[0]: position along the x axis in millimeters
	//  - Values[1]: position along the y axis in millimeters
	//  - Values[2]: position along the z axis in millimeters
	Values []float64
}

// Orientation is a 3D vector representing the orientation of the device.
type Orientation struct {
	// Timestamp is a event time in nanoseconds.
	Timestamp time.Time

	// Values is the orientation data.
	//  - Values[0]: rotation around the x axis in radians (pitch)
	//  - Values[1]: rotation around the y axis in radians (roll)
	//  - Values[2]: rotation around the z axis in radians (yaw)
	Values []float64
}

// Simulator simulates the device's position and orientation.
type Simulator struct {
	start time.Time
	end   time.Time

	activity      float64
	activityCurve spline.Spline

	positionCurveX   spline.Spline
	positionCurveY   spline.Spline
	positionCurveZ   spline.Spline
	positionInterval int

	orientationCurveX   spline.Spline
	orientationCurveY   spline.Spline
	orientationCurveZ   spline.Spline
	orientationInterval int
}

// NewSimulator returns a new simulator.
func NewSimulator(start, end time.Time, activity float64) *Simulator {
	duration := end.Sub(start)

	s0, _ := spline.NewSpline(spline.Bspline, randomControlPointsActivity(activity, duration))

	sim := &Simulator{start: start, end: end, activity: activity, activityCurve: s0, positionInterval: 50, orientationInterval: 100}

	s1, _ := spline.NewSpline(spline.CatmullRom, sim.randomControlPointsPosition(0, duration, -5.0, 5.0))
	s2, _ := spline.NewSpline(spline.CatmullRom, sim.randomControlPointsPosition(0, duration, -5.0, 5.0))
	s3, _ := spline.NewSpline(spline.CatmullRom, sim.randomControlPointsPosition(0, duration, -5.0, 5.0))

	s4, _ := spline.NewSpline(spline.CatmullRom, sim.randomControlPointsOrientation(rand.Float64()*1.5, duration, -0.5*math.Pi, 0.5*math.Pi))
	s5, _ := spline.NewSpline(spline.CatmullRom, sim.randomControlPointsOrientation(rand.Float64()-0.5, duration, -0.25*math.Pi, 0.25*math.Pi))
	s6, _ := spline.NewSpline(spline.CatmullRom, sim.randomControlPointsOrientation((rand.Float64()-0.5)*math.Pi, duration, -math.Pi, math.Pi))

	sim.positionCurveX = s1
	sim.positionCurveY = s2
	sim.positionCurveZ = s3

	sim.orientationCurveX = s4
	sim.orientationCurveY = s5
	sim.orientationCurveZ = s6

	return sim
}

func randomControlPointsActivity(activity float64, d time.Duration) []spline.Point {
	var points []spline.Point

	for i := 0; i < int(d.Milliseconds()); i += rand.Intn(1000) + 10 {
		value := rand.Float64() * activity
		points = append(points, spline.Point{
			X: float64(i),
			Y: value,
		})
	}
	return points
}

func (s *Simulator) randomControlPointsPosition(offset float64, d time.Duration, min, max float64) []spline.Point {
	var points []spline.Point

	delay := rand.Float64()
	intensity := rand.Float64()*4 + 1

	for i := 0; i < int(d.Milliseconds()); i += s.positionInterval {
		if rand.Float64() < s.activityCurve.At(float64(i)).Y {
			offset += (rand.Float64() - 0.5) * intensity
		}
		noise := rand.Float64()
		value := 0.005*math.Sin(float64(i)/float64(s.positionInterval)+delay+noise) + offset
		points = append(points, spline.Point{
			X: float64(i),
			Y: between(value, min, max),
		})
	}
	return points
}

func (s *Simulator) randomControlPointsOrientation(offset float64, d time.Duration, min, max float64) []spline.Point {
	var points []spline.Point

	delay := rand.Float64()
	intensity := rand.Float64()
	invert := rand.Float64() < 0.5

	for i := 0; i < int(d.Milliseconds()); i += s.orientationInterval {
		if rand.Float64() < s.activityCurve.At(float64(i)).Y {
			offset += (rand.Float64() - 0.5) * intensity
		}
		noise := rand.Float64()
		var value float64
		if invert {
			value = 0.0005*(math.Sin(float64(i)/float64(s.orientationInterval)+delay)+noise) + offset
		} else {
			value = 0.0005*(math.Cos(float64(i)/float64(s.orientationInterval)+delay)+noise) + offset
		}
		points = append(points, spline.Point{
			X: float64(i),
			Y: between(value, min, max),
		})
	}
	return points
}

func (s *Simulator) Position(t time.Time) Position {
	if t.Before(s.start) {
		t = s.start
	}
	if t.After(s.end) {
		t = s.end
	}
	i := float64(t.Sub(s.start).Milliseconds()) / float64(s.positionInterval)
	return Position{
		Timestamp: t,
		Values: []float64{
			float64(s.positionCurveX.At(i).Y),
			float64(s.positionCurveY.At(i).Y),
			float64(s.positionCurveZ.At(i).Y),
		},
	}
}

func (s *Simulator) Orientation(t time.Time) Orientation {
	if t.Before(s.start) {
		t = s.start
	}
	if t.After(s.end) {
		t = s.end
	}
	i := float64(t.Sub(s.start).Milliseconds()) / float64(s.orientationInterval)
	return Orientation{
		Timestamp: t,
		Values: []float64{
			betweenPi(float64(s.orientationCurveX.At(i).Y)),
			betweenPi(float64(s.orientationCurveY.At(i).Y)),
			betweenPi(float64(s.orientationCurveZ.At(i).Y)),
		},
	}
}

func betweenPi(value float64) float64 {
	for value > math.Pi {
		value -= 2 * math.Pi
	}
	for value < -math.Pi {
		value += 2 * math.Pi
	}
	return value
}

func between(value, min, max float64) float64 {
	if value < min {
		value += (min - value) * 2
	}
	if value > max {
		value -= (value - max) * 2
	}
	return value
}
