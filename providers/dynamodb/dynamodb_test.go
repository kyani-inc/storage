package dynamodb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/kyani-inc/storage/providers/dynamodb"
	"github.com/subosito/gotenv"
)

var (
	ddb       dynamodb.DynamoDB
	envLoaded = false
	region    = ""
	dbtable   = ""
	endpoint  = ""
)

func checkEnv(t *testing.T) {
	if envLoaded == false {
		gotenv.Load(".env")
		region = os.Getenv("AWS_REGION") //"us-west-2"
		dbtable = os.Getenv("DYNAMO_DB_TABLE")
		endpoint = os.Getenv("DYNAMO_DB_ENDPOINT")
		envLoaded = true
	}

	if region == "" || dbtable == "" || endpoint == "" {
		t.Skip("env vars not defined")
	}
}

func TestConnect(t *testing.T) {
	checkEnv(t)

	var err error
	ddb, err = dynamodb.New(region, endpoint, dbtable)

	if err != nil {
		t.Fatal("Failed to establish connection with DynamoDB!")
	} else {
		t.Log("Connected to local DynamoDB server")
	}
}

func TestPut(t *testing.T) {
	checkEnv(t)

	blah := []byte("hello, world!!")
	err := ddb.Put("test1", blah)

	if err != nil {
		t.Error(err.Error())
	}

	blah = []byte(`{"hello":"world"}`)
	err = ddb.Put("test2", blah)

	if err != nil {
		t.Error(err.Error())
	}

	blah = []byte("")
	err = ddb.Put("nodata", blah)

	if err == nil {
		t.Error("An error was expected but passed for some reason..")
	}
}

func TestGet(t *testing.T) {
	checkEnv(t)

	data := ddb.Get("test1")

	if strings.Contains(string(data), "hello, world") == false {
		t.Error("item `test1` does not contain expected values")
	}

	data = ddb.Get("test2")

	if strings.Contains(string(data), `{"hello":"world"}`) == false {
		t.Error("item `test2` does not contain expected values")
	}

	data = ddb.Get("nodata")

	if strings.Contains(string(data), "") == false {
		t.Error("item `nodata` should not contain data")
	}
}

func TestDelete(t *testing.T) {
	checkEnv(t)

	ddb.Delete("test2")

	data := ddb.Get("test2")

	if strings.Contains(string(data), "world") == true {
		t.Error("key test2 was not deleted!")
	}
}

func TestFlush(t *testing.T) {
	checkEnv(t)

	ddb.Flush()

	exists := ddb.TableExists()

	if exists {
		t.Error("DB Table was not 'flushed'. Failed to remove table.")
	}
}
