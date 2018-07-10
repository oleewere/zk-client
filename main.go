package main

import "os"
import (
	"flag"
	"fmt"
	"github.com/oleewere/go-solr-client/solr"
	"github.com/satori/go.uuid"
)

var Version string
var GitRevString string

func main() {

	var collection string
	var url string
	var krb5Path string
	var keytabPath string
	var principal string
	var realm string
	var isVersionCheck bool

	flag.BoolVar(&isVersionCheck, "version", false, "Print application version and git revision if available")
	flag.StringVar(&url, "url", "http://localhost:8983", "URL name for Solr or Solr proxy")
	flag.StringVar(&collection, "collection", "hadoop_logs", "Collection name for the Solr client")
	flag.StringVar(&krb5Path, "krb-conf-path", "", "Kerberos config location")
	flag.StringVar(&keytabPath, "keytab-path", "", "Kerberos keytab location")
	flag.StringVar(&principal, "principal", "", "Kerberos principal")
	flag.StringVar(&realm, "realm", "", "Kerberos Realm e.g.: EXAMPLE.COM")

	flag.Parse()

	if isVersionCheck {
		if GitRevString == "" {
			fmt.Println("version:", Version)
		} else {
			fmt.Printf("version: %s (git revision: %s)\n", Version, GitRevString)
		}
		os.Exit(0)
	}

	fmt.Print("Start Solr Cloud Client ...\n")

	securityConfig := solr.InitSecurityConfig(krb5Path, keytabPath, principal, realm)

	solrConfig := solr.SolrConfig{url, "hadoop_logs", &securityConfig, "/solr",
		solr.TLSConfig{}, true, 60}
	solrClient, err := solr.NewSolrClient(url, collection, &solrConfig)

	_, response, _ := solrClient.Query(nil)
	docs := response.Response.Docs
	for _, doc := range docs {
		fmt.Printf("----------------------")
		for k, v := range doc {
			fmt.Print("key: ", k)
			fmt.Println(" , value: ", v)
		}
		fmt.Printf("----------------------")
	}

	putDocs := solr.SolrDocuments{
		solr.SolrDocument{
			"id":          uuid.NewV4(),
			"log_message": "oleewere@gmail.com",
			"seq_num":     100,
			"level":       "FATAL",
			"logtime":     "2018-07-03T15:55:47.396Z",
		},
		solr.SolrDocument{
			"id":          uuid.NewV4(),
			"log_message": "oleewere@gmail.com",
			"seq_num":     1000,
			"level":       "FATAL",
			"logtime":     "2018-07-03T15:55:47.396Z",
		},
	}

	solrClient.Update(putDocs, nil, false)

	if err != nil {
		fmt.Print(err)
	}

	os.Exit(0)
}
