package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"time"
)

type Product struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func main() {
	es, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		// Handle error
		panic(err)
	}
	defer es.Stop()

	//p1 := `{"title" : "华为", "description" : "新品上市","price":3888.88 }`
	//product1, err := es.Index().
	//	Index("web_log").
	//	Id("3").
	//	BodyString(p1).
	//	Do(context.Background())
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}
	//fmt.Printf("Indexed tweet %s to index %s, type %s\n", product1.Id, product1.Index, product1.Type)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	p2 := Product{Title: "小米", Description: "are you ok!", Price: 2338.11}
	product2, err := es.Index().
		Index("web_log").
		BodyJson(p2).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", product2.Id, product2.Index, product2.Type)
	//
	//// Get tweet with specified ID
	//get1, err := es.Get().
	//	Index("web").
	//	Id("1").
	//	Do(context.Background())
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}
	//if get1.Found {
	//	fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	//}

}
