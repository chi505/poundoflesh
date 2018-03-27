package main

import "math/rand"

var MAXMEAT = 512
var NUMPEOPLE = 512

type Person struct {
	Name  string
	State PersonalState
	ID    int
}

type PersonalState struct {
	Meat     int //quantity
	Altruism int //Will be generous with probability (1 + exp(-altruism - meatdelta))^-1
}

type WorldState struct {
	MeatLossFrac           float64
	PerRoundLossFrac       float64
	NewEntrantMeanMeat     int
	NewEntrantMeanAltruism int
	UpdateProbPerRound     float64
	People                 []Person
}

func (world WorldState) initializeState() {
	for i := 0; i < NUMPEOPLE; i++ {
		world.People = append(world.People, MakeNewPerson())
	}
}

func (world WorldState) updateState() {
	count.Value++
}

func (world WorldState) interact(agent *Person, patient *Person) {
	meatAmount := agent.PullARequestAmount(patient.State)
	if patient.WouldAcceptOfferFrom(agent.State, meatAmount) {
		patient.State.Meat -= meatAmount
		agent.State.Meat += meatAmount * (1 - MeatLossFrac)
	}
	agent.State.Meat = agent.State.Meat * (1 - PerRoundLossFrac)
	patient.State.Meat = patient.State.Meat * (1 - PerRoundLossFrac)

	if agent.State.Meat < 0 {
		world.People[agent.ID] = world.MakeNewPerson(agent.ID)
	}

	if patient.State.Meat < 0 {
		world.People[patient.ID] = world.MakeNewPerson(patient.ID)
	}
}

func (agent *Person) PullARequestAmount(ps PersonalState) int {
	return rand.Intn(ps.Meat)
}

func (patient *Person) WouldAcceptOfferFrom(as PersonalState, amount int) bool {
	return true
}

func (world WorldState) MakeNewPerson(id int) Person {
	return Person{
		Name:  MakeNewName(),
		State: PersonalState{Meat: rand.Intn(NewEntrantMeanMeat * 2), Altruism: rand.Intn(NewEntrantMeanAltruism * 2)},
		ID:    id}

}

func MakeNewName() string {
	return GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar()
}

func GetNextNameChar() {
	return string(rand.Intn(126-33) + 33)
}
