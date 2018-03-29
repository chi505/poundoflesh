package main

import "math/rand"
import "time"

var MAXMEAT = 512.0
var NUMPEOPLE = 2
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
	world.loadMeatMaps()
	for i := 0; i < NUMPEOPLE; i++ {
		world.People = append(world.People, world.MakeNewPerson(i))
	}
}

func (world *WorldState) loadMeatMaps() {
	world.PersonSpec["kidney"] = MeatSpec{Count: 2, MeanInitMeat: int(MAXMEAT / 10)}
	world.Assets.Organs["kidney"] = append([]MeatData{}, MeatData{Description: "A glistening reddish brown bean shaped chunk of MEAT."})
	world.PersonSpec["heart"] = MeatSpec{Count: 1, MeanInitMeat: int(MAXMEAT)}
	world.Assets.Organs["heart"] = append(make([]MeatData, 0), MeatData{Description: "A throbbing, beating, dripping, symbolic heart."})
}

func (world *WorldState) updateState() {
	permVec := rand.Perm(NUMPEOPLE)
	for i := 0; i < NUMPEOPLE/2; i++ {
		world.interact(world.People[permVec[2*i]], world.People[permVec[2*i+1]])
		world.Count++
	}
	for i, person := range world.People {
		if world.MassageMeat(person) == 0 {
			world.People[i] = world.MakeNewPerson(i) // could do this in MassageMeat but making replacement more explicit
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

func (world *WorldState) MassageMeat(p *Person) int {
	n := len(p.State.MeatBag) - 1
	count := n + 1
	for i := range p.State.MeatBag {
		p.State.MeatBag[n-i].Meat -= MEATDEC
		if p.State.MeatBag[n-i].Meat <= 0 {
			p.State.MeatBag = append(p.State.MeatBag[:n-i], p.State.MeatBag[n-i+1:]...)
			count--
		}
	}
	return count
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

func (noob *Person) InsertMeat(assets TextAssets, specs map[string]MeatSpec) {
	for name, spec := range specs {
		for i := 0; i < spec.Count; i++ {
			noob.State.MeatBag = append(noob.State.MeatBag,
				MeatPiece{
					Name: name,
					Data: MeatData{
						Description: assets.Organs[name][rand.Intn(len(assets.Organs[name]))].Description + " It originally belonged to " + noob.Name + "."},
					Meat: spec.MeanInitMeat/2 + rand.Intn(spec.MeanInitMeat)})
		}
	}
}

func MakeNewName() string {
	return GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar()
}

func GetNextNameChar() string {
	return string(rand.Intn(126-33) + 33)
}
