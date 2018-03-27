package main

import "math/rand"
import "time"

var MAXMEAT = 512.0
var NUMPEOPLE = 16

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
	NewEntrantMeanMeat     float64
	NewEntrantMeanAltruism int
	UpdateProbPerRound     float64
	People                 []*Person
	Count                  int
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
	for i := 0; i < NUMPEOPLE/2; i++ {
		world.interact(world.People[permVec[2*i]], world.People[permVec[2*i+1]])
		world.Count++
	}
}

func (world *WorldState) interact(agent *Person, patient *Person) {
	meatAmount := agent.PullARequestAmount(patient.State)
	if patient.WouldAcceptOfferFrom(agent.State, meatAmount) {
		patient.State.Meat -= meatAmount
		agent.State.Meat += meatAmount * (1.0 - world.MeatLossFrac)
	}
	agent.State.Meat -= MAXMEAT * world.PerRoundLossFrac
	patient.State.Meat -= MAXMEAT * world.PerRoundLossFrac

	if agent.State.Meat < 0 {
		id := agent.ID
		world.People[id] = world.MakeNewPerson(id)
	}

	if patient.State.Meat < 0 {
		id := patient.ID
		world.People[id] = world.MakeNewPerson(id)
	}
}

func (agent *Person) PullARequestAmount(ps PersonalState) float64 {
	return float64(rand.Intn(int(ps.Meat)) + 1)
}

func (patient *Person) WouldAcceptOfferFrom(as PersonalState, amount float64) bool {
	return true
}

func (world WorldState) MakeNewPerson(id int) *Person {
	return &Person{
		Name:  MakeNewName(),
		State: PersonalState{Meat: float64(rand.Intn(int(world.NewEntrantMeanMeat * 2))), Altruism: rand.Intn(world.NewEntrantMeanAltruism * 2)},
		ID:    id}

}

func MakeNewName() string {
	return GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar()
}

func GetNextNameChar() string {
	return string(rand.Intn(126-33) + 33)
}
