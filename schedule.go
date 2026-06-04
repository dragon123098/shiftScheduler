package main

import (
	"fmt"
)

// Schedule struct to hold the schedule for each day of the week
type Schedule struct {
	Monday [2]float64
	Tuesday [2]float64
	Wednesday [2]float64
	Thursday [2]float64
	Friday [2]float64
} 


var shifts = map[string]Schedule{} // Map to store the schedule for each user

// SetSchedule method to set the schedule for a specific day
func (s *Schedule) SetSchedule(d string, start float64, end float64) [2]float64 {
	switch d {
	case "Monday":
		s.Monday[0] = start
		s.Monday[1] = end
		return s.Monday
	case "Tuesday":
		s.Tuesday[0] = start
		s.Tuesday[1] = end
		return s.Tuesday	
	case "Wednesday":
		s.Wednesday[0] = start
		s.Wednesday[1] = end
		return s.Wednesday	
	case "Thursday":
		s.Thursday[0] = start
		s.Thursday[1] = end
		return s.Thursday
	case "Friday":
		s.Friday[0] = start
		s.Friday[1] = end
		return s.Friday
	default:
		fmt.Println("Invalid day")
		return [2]float64{0, 0}
	}
}

func (s *Schedule) PrintSchedule() {
	fmt.Printf("Monday: %.1f - %.1f\n", s.Monday[0], s.Monday[1])
	fmt.Printf("Tuesday: %.1f - %.1f\n", s.Tuesday[0], s.Tuesday[1])
	fmt.Printf("Wednesday: %.1f - %.1f\n", s.Wednesday[0], s.Wednesday[1])
	fmt.Printf("Thursday: %.1f - %.1f\n", s.Thursday[0], s.Thursday[1])
	fmt.Printf("Friday: %.1f - %.1f\n", s.Friday[0], s.Friday[1])
}