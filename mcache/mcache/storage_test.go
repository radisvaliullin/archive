package mcache

import (
	"testing"
	"time"
	"reflect"
)

func TestSetString(t *testing.T) {

	store := NewStorage()

	setStr := "TEST STRING"

	err := store.Set("test", setStr, time.Second*3600)
	if err != nil {
		t.Fatal("set err")
	}

	sv := store.Get("test")
	if sv == nil {
		t.Fatal("get store value, empty")
	}

	getStr, ok := sv.GetString()
	if !ok {
		t.Fatal("store value is not string")
	}

	if setStr == getStr {
		t.Log("setStr equal getStr")
	} else {
		t.Fatal("setStr not equal getStr")
	}
}

//
func TestSetMap(t *testing.T) {

	store := NewStorage()

	setMap := map[string]string{"test":"TEST STRING"}

	err := store.Set("test", setMap, time.Second*3600)
	if err != nil {
		t.Fatal("set err")
	}

	sv := store.Get("test")
	if sv == nil {
		t.Fatal("get store value, empty")
	}

	getMap, ok := sv.GetMap()
	if !ok {
		t.Fatal("store value is not map")
	}

	eq := reflect.DeepEqual(setMap, getMap)
	if eq {
		t.Log("setMap equal getMap")
	} else {
		t.Fatal("setMap not equal getMap")
	}
}
