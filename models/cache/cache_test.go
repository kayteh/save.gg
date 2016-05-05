package cache

import (
	_ "encoding/json"
	"log"
	"os"
	dummyCache "save.gg/sgg/models/cache/dummy"
	"testing"
	"time"
)

var c *Cache

type testType struct {
	Success bool `json:"success"`
}

func TestMain(m *testing.M) {

	b, err := dummyCache.NewDummyCache()
	if err != nil {
		log.Fatalln(err)
	}

	i, err := NewCache(b)
	if err != nil {
		log.Fatalln(err)
	}

	c = i

	os.Exit(m.Run())

}

func TestMarshalling(t *testing.T) {

	testData := testType{Success: true}
	err := c.Set("test-marshal", testData, 10*time.Minute)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var data testType
	err = c.Get("test-marshal", &data)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if data.Success != true {
		t.Error("cache missed")
		t.Fail()
	}

}

func TestBasicTTL(t *testing.T) {

	testData := testType{Success: true}
	err := c.Set("test-marshal", testData, 1*time.Second)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	time.Sleep(2 * time.Second)

	var data testType
	err = c.Get("test-marshal", &data)
	if err != nil && err.Error() != "cache miss" {
		t.Error(err)
		t.FailNow()
	}

}
