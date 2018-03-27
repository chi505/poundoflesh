package main

import "math/rand"
import "time"

var MAXMEAT = 512
var NUMPEOPLE = 512

type Person struct {
	Name  string
	State PersonalState
	ID    int
}

type PersonalState struct {
	Meat     float64 //quantity
	Altruism int     //Will be generous with probability (1 + exp(-altruism - meatdelta))^-1
}

type WorldState struct {
	MeatLossFrac           float64
	PerRoundLossFrac       float64
	NewEntrantMeanMeat     int
	NewEntrantMeanAltruism int
	UpdateProbPerRound     float64
	People                 []*Person
}

func (world *WorldState) initializeState() {
	now := time.Now()
	var seed int64
	seed = now.Unix()
	rand.Seed(seed)
	for i := 0; i < NUMPEOPLE; i++ {
		world.People = append(world.People, world.MakeNewPerson(i))
	}
}

func (world *WorldState) updateState() {
	permVec := rand.Perm(NUMPEOPLE)
	for i := 0; i < NUMPEOPLE; i++ {
		world.interact(world.People[permVec[i]], world.People[permVec[i+1]])
	}
}

func (world *WorldState) interact(agent *Person, patient *Person) {
	meatAmount := agent.PullARequestAmount(patient.State)
	if patient.WouldAcceptOfferFrom(agent.State, meatAmount) {
		patient.State.Meat -= meatAmount
		agent.State.Meat += meatAmount * (1 - world.MeatLossFrac)
	}
	agent.State.Meat -= float64(MAXMEAT) * world.PerRoundLossFrac
	patient.State.Meat -= float64(MAXMEAT) * world.PerRoundLossFrac

	if agent.State.Meat < 0 {
		world.People[agent.ID] = world.MakeNewPerson(agent.ID)
	}

	if patient.State.Meat < 0 {
		world.People[patient.ID] = world.MakeNewPerson(patient.ID)
	}
}

func (agent *Person) PullARequestAmount(ps PersonalState) float64 {
	return float64(rand.Intn(int(ps.Meat)))
}

func (patient *Person) WouldAcceptOfferFrom(as PersonalState, amount float64) bool {
	return true
}

func (world WorldState) MakeNewPerson(id int) *Person {
	return &Person{
		Name:  MakeNewName(),
		State: PersonalState{Meat: float64(rand.Intn(world.NewEntrantMeanMeat * 2)), Altruism: rand.Intn(world.NewEntrantMeanAltruism * 2)},
		ID:    id}

}

func MakeNewName() string {
	return GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar()
}

func GetNextNameChar() string {
	return string(rand.Intn(126-33) + 33)
}
