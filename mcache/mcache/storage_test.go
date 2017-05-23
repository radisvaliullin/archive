package mcache

import (
	"testing"
	"time"
)

func TestSetString(t *testing.T) {

	store := NewStorage()

	str := "TEST73"
	err := store.Set("key", str, time.Second*3600)
	if err != nil {
		t.Fatal("Set err")
	}
	t.Logf("STORE %+v", store)

	sv := store.Get("key")
	if sv == nil {
		t.Fatal("get empty")
	}
	t.Logf("STORE VAL %+v", sv)

	str, ok := sv.GetString()
	if !ok {
		t.Logf("STORE %+v", store)
		t.Fatal("sv is not string")
	}
	t.Log("str var ", str)

}
