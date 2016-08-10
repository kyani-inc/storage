package elasticsearch

import (
	"encoding/json"
	"github.com/smartystreets/go-aws-auth"
	"gopkg.in/olivere/elastic.v3"
	"net/http"
)

type ElasticSearch struct {
	client    *elastic.Client
	index     string
	namespace string
}

type awsSigningTransport struct {
	HTTPClient  *http.Client
	Credentials awsauth.Credentials
}

// New creates a new instance of an ES client for interaction
// with an AWS provied ES server
func New(host, scheme, index, namespace, awsKey, awsSecret string) (ElasticSearch, error) {
	var e ElasticSearch

	signingTransport := awsSigningTransport{
		Credentials: awsauth.Credentials{
			AccessKeyID:     awsKey,
			SecretAccessKey: awsSecret,
		},
		HTTPClient: http.DefaultClient,
	}

	signingClient := http.Client{Transport: http.RoundTripper(signingTransport)}

	client, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetMaxRetries(10),
		elastic.SetScheme(scheme),
		elastic.SetHttpClient(&signingClient),
	)

	if err != nil {
		return e, err
	}

	// the namespace becomes the doctype for ease of use
	e.client = client
	e.index = index
	e.namespace = namespace

	return e, err
}

// Put overwrites or creates a new ES entry based on key, index, and doctype
func (e ElasticSearch) Put(key string, data []byte) error {
	// check for the index, create if does not exist
	exists, err := e.client.IndexExists(e.index).Do()

	if err != nil {
		return err
	}

	if !exists {
		createIndex, err := e.client.CreateIndex(e.index).Do()

		if err != nil {
			return err
		}

		// it is possible that no error is thrown the index creation failed
		if !createIndex.Acknowledged {
			return err
		}
	}

	var container interface{}

	err = json.Unmarshal(data, &container)

	if err != nil {
		return err
	}

	_, err = e.client.Index().
		Index(e.index).
		Type(e.namespace).
		Id(key).
		BodyJson(container).
		Do()

	if err != nil {
		return err
	}

	return nil
}

// Get attempts to grab ES entries by key, index and doctype
func (e ElasticSearch) Get(key string) []byte {
	data, err := e.client.Get().
		Index(e.index).
		Type(e.namespace).
		Id(key).
		Do()

	if err != nil || data.Id != key {
		return []byte{}
	}

	res, err := json.Marshal(data.Source)

	if err != nil {
		return []byte{}
	}

	return res
}

// Delete removes an ES entry by key
func (e ElasticSearch) Delete(key string) {
	e.client.Delete().
		Index(e.index).
		Type(e.namespace).
		Id(key).
		Do()
}

// Flush tells ES to free memory from the index and flush data to disk.
func (e ElasticSearch) Flush() {
	e.client.Flush().
		Index(e.index).
		AllowNoIndices(true).
		IgnoreUnavailable(true).
		Do()
}

// RoundTrip implementation
func (a awsSigningTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return a.HTTPClient.Do(awsauth.Sign4(req, a.Credentials))
}
