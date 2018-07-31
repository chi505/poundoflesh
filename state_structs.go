package main

var MAXMEAT = 128.0
var NUMPEOPLE = 20
var MEATDEC = 5

type Person struct {
	Name  string
	State PersonalState
	ID    int
}

type MentalState struct {
	Fear     int
	Hope     int
	Altruism int
	Caprice  int
}

type PersonalState struct {
	MeatTotal int //quantity
	MeatBag   []*MeatPiece
	Birthday  int
	Mind      MentalState
}

type MeatPiece struct {
	Name      string //this is outside of Data because used as key into various things
	Data      MeatData
	Meat      int
	OrigOwner string
}

type MeatData struct {
	Description string
}

type WorldState struct {
	Params     PoundOFleshParams
	People     []*Person
	Count      int
	Assets     TextAssets
	PersonSpec map[string]MeatSpec
}

type MeatSpec struct {
	Count        int
	MeanInitMeat int
}

type PoundOFleshParams struct {
	MeatLossFrac           float64
	PerRoundLossFrac       float64
	NewEntrantMeanMeat     int
	NewEntrantMeanAltruism int
	UpdateProbPerRound     float64
}

type TextAssets struct {
	Organs map[string][]MeatData
}
