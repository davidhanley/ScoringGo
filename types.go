package main

import (
  "time"
)

type AthleteRaceResult struct {
	ath    AthleteAndPoints
	race   *Race
	points float32
	rank   int
}

type Athlete struct {
	id        int
	name      string
	age       int
	sex       string
	foreign   bool
	racesDone int
}

type AthleteAndPoints struct {
	athlete *Athlete
	points  float32
}

type Race struct {
	name     string
	points   int
	date     time.Time
	athletes []*Athlete
}

type Foreignicity int

const (
	ALL Foreignicity = iota //starts with zero
	THREE_RACE_FOREIGNERS
	US_ONLY
)

type CategoryResult struct {
	gender         string
	ageLow         int
	ageHigh        int
	includeForeign Foreignicity
	results        map[string][]*AthleteRaceResult
	sortedAthletes []AthleteAndPoints
}
