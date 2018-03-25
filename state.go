package main

import "math/rand"

var MAXMEAT = 512

type Person struct {
    Name string
    State PersonalState
    
}

type PersonalState struct {
    Meat int
}


func initializeState(){
    count.Value = 0
    
    People = append(People, Person{"David", PersonalState{rand.Intn(MAXMEAT)}})
    People = append(People, Person{"Taniqua", PersonalState{rand.Intn(MAXMEAT)}})
}

func updateState() {
    count.Value++
}