package rediscache

import (
	redis "gopkg.in/redis.v3"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var c *RedisCache

func TestMain(m *testing.M) {

	i, err := NewRedisCache(&redis.Options{
		Addr: "localhost:6379",
	})
	if err != nil {
		log.Fatalln(err)
	}

	c = i

	c.Redis.FlushAll()

	os.Exit(m.Run())

}

func TestSet(t *testing.T) {
	err := c.Set("test-set", []byte(`{"success":true}`), 10*time.Minute)
	if err != nil {
		t.Errorf("set error: %v", err)
		t.FailNow()
	}
}

func TestGet(t *testing.T) {
	testData := []byte(`{"success":true}`)

	err := c.Set("test-get", testData, 10*time.Minute)
	if err != nil {
		t.Errorf("set error: %v", err)
		t.Fail()
	}

	data, err := c.Get("test-get")
	if err != nil {
		t.Errorf("get error: %v", err)
		t.Fail()
	}

	if !cmp(testData, data) {
		t.Errorf("expected `%s`, got `%s`", testData, data)
		t.Fail()
	}
}

func TestDel(t *testing.T) {
	err := c.Set("test-del", []byte(`{"success":false}`), 10*time.Minute)
	if err != nil {
		t.Errorf("set error: %v", err)
		t.Fail()
	}

	err = c.Del("test-del")
	if err != nil {
		t.Errorf("del error: %v", err)
		t.Fail()
	}

	data, err := c.Get("test-del")
	if err != nil {
		t.Errorf("get error: %v", err)
		t.Fail()
	}

	if data != nil {
		t.Errorf("expected nil, got %v", data)
		t.Fail()
	}
}

func TestTTL(t *testing.T) {
	err := c.Set("test-ttl", []byte(`{"success":false}`), 2*time.Second)
	if err != nil {
		t.Errorf("set error: %v", err)
		t.Fail()
	}

	time.Sleep(3 * time.Second)

	data, err := c.Get("test-ttl")
	if err != nil {
		t.Errorf("get error: %v", err)
		t.Fail()
	}

	if data != nil {
		t.Errorf("expected nil, got `%s`", data)
		t.Fail()
	}
}

func TestTTL2(t *testing.T) {
	testData := []byte(`{"success":false}`)
	err := c.Set("test-ttl2", testData, 4*time.Second)
	if err != nil {
		t.Errorf("set error: %v", err)
		t.Fail()
	}

	time.Sleep(2 * time.Second)

	data1, err := c.Get("test-ttl2")
	if err != nil {
		t.Errorf("get error: %v", err)
		t.Fail()
	}

	if !cmp(testData, data1) {
		t.Errorf("expected `%s`, got `%s`", testData, data1)
		t.Fail()
	}

	time.Sleep(3 * time.Second)

	data2, err := c.Get("test-ttl2")
	if err != nil {
		t.Errorf("get error: %v", err)
		t.Fail()
	}

	if data2 != nil {
		t.Errorf("expected nil, got `%s`", data2)
		t.Fail()
	}

}

func cmp(s1, s2 []byte) bool {
	return strings.Compare(string(s1), string(s2)) == 0
}
