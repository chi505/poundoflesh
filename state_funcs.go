package main

import "math/rand"
import "time"
import "math"
import "sort"

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
	world.Assets.Organs["kidney"] = append(world.Assets.Organs["kidney"], MeatData{Description: "A bean shaped brownish red chunk of meat."})

	world.PersonSpec["heart"] = MeatSpec{Count: 1, MeanInitMeat: int(MAXMEAT)}
	world.Assets.Organs["heart"] = append(make([]MeatData, 0), MeatData{Description: "A throbbing, beating, dripping, symbolic heart."})

	world.PersonSpec["lung"] = MeatSpec{Count: 2, MeanInitMeat: int(MAXMEAT / 5)}
	world.Assets.Organs["lung"] = append([]MeatData{}, MeatData{Description: "A glistening spongy lung shaped chunk of MEAT."})
	world.Assets.Organs["lung"] = append(world.Assets.Organs["lung"], MeatData{Description: "A lung shaped sponge meat."})

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
	if rand.Float64() < world.Params.UpdateProbPerRound {
		world.Params.JitterParams()
	}

	sort.Slice(world.People, func(i, j int) bool { return len(world.People[i].State.MeatBag) > len(world.People[j].State.MeatBag) })
}

func (params *PoundOFleshParams) JitterParams() {
	params.MeatLossFrac += ClampF64(rand.NormFloat64()*0.02, 1, 0)
	params.PerRoundLossFrac += ClampF64(rand.NormFloat64()*0.05, 1, 0)
}

func (world *WorldState) interact(agent *Person, patient *Person) {
	meat := agent.PullAMeatRequest(patient)
	if patient.WouldAcceptOfferFrom(agent.State, meat) {
		patient.GiveMeatTo(agent, meat)
	}
}

func (agent *Person) PullAMeatRequest(patient *Person) *MeatPiece {
	index := rand.Intn(patient.State.MeatTotal)
	meat, valid := patient.GetMeatByWeight(index)
	if valid {
		return meat
	}
	return nil
}

func (patient *Person) WouldAcceptOfferFrom(as PersonalState, request *MeatPiece) bool {
	return true
}

//insertion can't logically be impossible
func (person *Person) AddMeat(meat *MeatPiece) {
	person.State.MeatBag = append(person.State.MeatBag, meat)
	person.State.MeatTotal += meat.Meat
}

func (person *Person) GetMeatIndex(meat *MeatPiece) (int, bool) {
	for meatIndex := range person.State.MeatBag {
		if person.State.MeatBag[meatIndex] == meat {
			return meatIndex, true
		}
	}
	return 0, false
}

func (person *Person) GetMeatByWeight(weight int) (*MeatPiece, bool) {
	sum := 0
	meatbag := person.State.MeatBag
	for meatIndex := range meatbag {
		sum += meatbag[meatIndex].Meat
		if sum >= weight {
			return meatbag[meatIndex], true
		}
	}
	return nil, false
}

//need return value in case we get misaskedfor meat
func (person *Person) RemoveMeat(meat *MeatPiece) bool {
	for meatIndex := range person.State.MeatBag {
		if person.State.MeatBag[meatIndex] == meat {
			person.State.MeatBag = append(person.State.MeatBag[:meatIndex], person.State.MeatBag[meatIndex+1:]...)
			person.State.MeatTotal -= meat.Meat
			return true
		}
	}
	return false
}

//It's easier to not lose the meat mid-transfer if it's passed in as an argument
func (giver *Person) GiveMeatTo(recip *Person, meat *MeatPiece) bool {
	_, valid := giver.GetMeatIndex()
	if valid {
		recip.AddMeat(meat)
		giver.RemoveMeat(meat) //don't need to check return because GetMeatIndex already does
		return true
	}
	return false
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
			Meat:     rand.Intn(int(world.Params.NewEntrantMeanMeat * 2)),
			Altruism: rand.Intn(world.Params.NewEntrantMeanAltruism * 2),
			MeatBag:  make([]MeatPiece, 0),
			Birthday: world.Count},
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
						Description: assets.Organs[name][rand.Intn(len(assets.Organs[name]))].Description},
					Meat:      int(math.Min(float64(spec.MeanInitMeat/2+rand.Intn(spec.MeanInitMeat)), MAXMEAT)),
					OrigOwner: noob.Name})
		}
	}
}

func MakeNewName() string {
	return GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar() + GetNextNameChar()
}

func GetNextNameChar() string {
	return string(rand.Intn(126-33) + 33)
}

func ClampF64(input float64, upper float64, lower float64) float64 {
	return math.Min(upper, Math.Max(lower, input))
}
