package dynamodb_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kyani-inc/storage/providers/dynamodb"
	"github.com/subosito/gotenv"
)

var (
	ddb      dynamodb.DynamoDB
	region   = ""
	dbtable  = ""
	endpoint = ""
)

func TestMain(m *testing.M) {
	gotenv.Load(".env")

	region = os.Getenv("AWS_REGION") //"us-west-2"
	dbtable = os.Getenv("DYNAMO_DB_TABLE")
	endpoint = os.Getenv("DYNAMO_DB_ENDPOINT")
	os.Exit(m.Run())
}

func TestConnect(t *testing.T) {
	if region == "" || dbtable == "" || endpoint == "" {
		t.Fatal("Missing required env vars!")
	}

	var err error
	ddb, err = dynamodb.New(region, endpoint, "test_table")

	if err != nil {
		t.Fatal("Failed to establish connection with DynamoDB!")
	} else {
		fmt.Println("Connected to local DynamoDB server")
	}
}

func TestPut(t *testing.T) {
	blah := []byte("hello, world!!")
	_ = ddb.Put("test1", blah)

	blah = []byte(`{"hello":"world"}`)
	_ = ddb.Put("test2", blah)

	blah = []byte{}
	_ = ddb.Put("nodata", blah)
}

func TestGet(t *testing.T) {
	data := ddb.Get("test1")
	fmt.Println(string(data))

	data = ddb.Get("test2")
	fmt.Println(string(data))

	data = ddb.Get("nodata")
	fmt.Println(string(data))
}

func TestDelete(t *testing.T) {
	ddb.Delete("test2")

	data := ddb.Get("test2")

	if strings.Contains(string(data), "world") == true {
		t.Error("key test2 was not deleted!")
	}
}

func TestFlush(t *testing.T) {
	ddb.Flush()

	exists := ddb.TableExists()

	if exists {
		t.Error("DB Table was not 'flushed'. Failed to remove table.")
	}
}
