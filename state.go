package main

import "math/rand"
import "time"

var MAXMEAT = 512.0
var NUMPEOPLE = 16
var MEATDEC = 5

type Person struct {
	Name  string
	State PersonalState
	ID    int
}

type PersonalState struct {
	Meat     float64 //quantity
	Altruism int     //Will be generous with probability (1 + exp(-altruism - meatdelta))^-1
	MeatBag  []MeatPiece
}

type MeatPiece struct {
	Name string
	Data MeatData
	Meat int
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
	NewEntrantMeanMeat     float64
	NewEntrantMeanAltruism int
	UpdateProbPerRound     float64
}

type TextAssets struct {
	Organs map[string][]MeatData
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

func (world *WorldState) loadMeatMaps() {
	world.PersonSpec["kidney"] = MeatSpec{Count: 2, MeanInitMeat: int(MAXMEAT / 10)}
	world.Assets.Organs["kidney"] = append(make([]MeatData, 0), MeatData{Description: "A glistening reddish brown bean shaped chunk of MEAT"})
	world.PersonSpec["heart"] = MeatSpec{Count: 2, MeanInitMeat: int(MAXMEAT)}
	world.Assets.Organs["heart"] = append(make([]MeatData, 0), MeatData{Description: "A throbbing, beating, dripping, symbolic heart"})
}

func (world *WorldState) updateState() {
	permVec := rand.Perm(NUMPEOPLE)
	for i := 0; i < NUMPEOPLE/2; i++ {
		world.interact(world.People[permVec[2*i]], world.People[permVec[2*i+1]])
		world.Count++
	}
	for i, person := range world.People {
		world.MassageMeat(person)
		if len(person.State.MeatBag) == 0 {
			person = world.MakeNewPerson(i) // could do this in MassageMeat but making replacement more explicit
		}
	}
}

func (world *WorldState) interact(agent *Person, patient *Person) {
	meatIndex := agent.PullAMeatRequest(patient.State)
	if patient.WouldAcceptOfferFrom(agent.State, &patient.State.MeatBag[meatIndex]) {
		agent.State.MeatBag = append(agent.State.MeatBag, patient.State.MeatBag[meatIndex])
		//THESE LINES MUST BE IN THIS ORDER
		patient.State.MeatBag = append(patient.State.MeatBag[:meatIndex], patient.State.MeatBag[meatIndex+1:]...)
	}
}

func (agent *Person) PullAMeatRequest(ps PersonalState) int {
	return rand.Intn(len(ps.MeatBag))
}

func (patient *Person) WouldAcceptOfferFrom(as PersonalState, request *MeatPiece) bool {
	return true
}

func (world *WorldState) MassageMeat(p *Person) {
	for i := range p.State.MeatBag {
		p.State.MeatBag[i].Meat -= MEATDEC
		if p.State.MeatBag[i].Meat <= 0 {
			p.State.MeatBag = append(p.State.MeatBag[:i], p.State.MeatBag[i+1:]...)
		}
	}
}

func (world WorldState) MakeNewPerson(id int) *Person {
	noob := &Person{
		Name: MakeNewName(),
		State: PersonalState{
			Meat:     float64(rand.Intn(int(world.Params.NewEntrantMeanMeat * 2))),
			Altruism: rand.Intn(world.Params.NewEntrantMeanAltruism * 2),
			MeatBag:  make([]MeatPiece, 0)},
		ID: id}
	noob.InsertMeat(world.Assets, world.PersonSpec)
	return noob

}

func (noob *Person) InsertMeat(assets Assets, spec map[string]MeatSpec) {
	for name, spec := range spec {
		for i := 0; i < spec.Count; i++ {
			noob.State.MeatBag = append(noob.State.MeatBag,
				MeatPiece{
					Name: "name",
					Data: MeatData{
						Description: assets.TextAssets.Organs[name][rand.Intn(assets.TextAssets.Organs[name])],
						Meat:        spec[name].MeanInitMeat/2 + rand.Intn(spec[name].MeanInitMeat)}})
		}
	}
}

func MakeNewName() string {
	return GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar()
}

func GetNextNameChar() string {
	return string(rand.Intn(126-33) + 33)
}
