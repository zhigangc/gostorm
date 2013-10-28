package gostorm

import "testing"
import "os"

type VV map[string]interface{}

func _Test1(t *testing.T) {
	comp := &Component{}
	file, err := os.Open("input.json")
	if err != nil {
		t.Fatalf("err: %s\n", err.Error())
	}
	setupInfo := comp.ReadValues(file)
    pidDir, _ := setupInfo.GetString("pidDir")
    comp.SendPid(pidDir)
    comp.Conf, _ = setupInfo.GetValues("conf")
    comp.Context, _ = setupInfo.GetValues("context")
}

func Test2(t *testing.T) {
	input := `{"id":"2362548963208114831","stream":"default","comp":"spout","tuple":["the cow jumped over the moon"],"task":4}`
	vals, err := ParseValues([]byte(input))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if val, ok := vals.GetInt("task"); !ok || val != 4 {
		t.Fatalf("expected an int (4)")
	}
}
