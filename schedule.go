package main

import (
	"fmt"
)

// Schedule struct to hold the schedule for each day of the week
type Schedule struct {
	Monday [2]int
	Tuesday [2]int
	Wednesday [2]int
	Thursday [2]int
	Friday [2]int
} 


var shifts = map[string]Schedule{} // Map to store the schedule for each user

// SetSchedule method to set the schedule for a specific day
func (s *Schedule) SetSchedule(d string, start int, end int) [2]int {
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
		return [2]int{0, 0}
	}
}

func (s *Schedule) GetSchedule(d string) [2]int {
	switch d {
	case "Monday":
		fmt.Printf("Monday: %d - %d\n", s.Monday[0], s.Monday[1])
		return s.Monday
	case "Tuesday":
		fmt.Printf("Tuesday: %d - %d\n", s.Tuesday[0], s.Tuesday[1])
		return s.Tuesday
	case "Wednesday":
		fmt.Printf("Wednesday: %d - %d\n", s.Wednesday[0], s.Wednesday[1])
		return s.Wednesday
	case "Thursday":
		fmt.Printf("Thursday: %d - %d\n", s.Thursday[0], s.Thursday[1])
		return s.Thursday
	case "Friday":
		fmt.Printf("Friday: %d - %d\n", s.Friday[0], s.Friday[1])
		return s.Friday
	default:
		fmt.Println("Invalid day")
		return [2]int{0, 0}
	}
}