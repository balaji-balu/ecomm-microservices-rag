package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"

    "github.com/olivere/elastic/v7"
)

// Define a struct to match the JSON structure
type Product struct {
    ProductID   string  `json:"product_id"`
    ProductName string  `json:"product_name"`
    Category    string  `json:"category"`
    StockLevel  int     `json:"stock_level"`
    Price       float64 `json:"price"`
    Description string  `json:"description"`
}

func main() {
    // Elasticsearch configuration
    esAddress := "http://localhost:9200"
    client, err := elastic.NewClient(
		elastic.SetBasicAuth("elastic", "yourpassword"),
		elastic.SetURL(esAddress), 
		elastic.SetSniff(false), 
		elastic.SetHealthcheck(false),
	)
    if err != nil {
        log.Fatalf("Error creating Elasticsearch client: %v", err)
    }

    // Open the JSON file
    file, err := os.Open("sample_inventory_data.json")
    if err != nil {
        log.Fatalf("Error opening JSON file: %v", err)
    }
    defer file.Close()

    // Read and parse JSON file
    byteValue, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatalf("Error reading JSON file: %v", err)
    }

    var products []Product
    if err := json.Unmarshal(byteValue, &products); err != nil {
        log.Fatalf("Error parsing JSON data: %v", err)
    }

	fmt.Printf("%v %v", context.Background(), client)
/*
	// Create a bulk service
	bulkRequest := client.Bulk()
	for _, product := range products {
		req := elastic.NewBulkIndexRequest().
			Index("inventory").
			Id(product.ProductID).
			Doc(product)
		bulkRequest = bulkRequest.Add(req)
	}
	
	// Execute the bulk request
	bulkResponse, err := bulkRequest.Do(context.Background())
	if err != nil {
		log.Fatalf("Error executing bulk request: %v", err)
	}
	
	// Check for any failures
	if bulkResponse.Errors {
		for _, item := range bulkResponse.Items {
			if item.Index.Error != nil {
				log.Printf("Error indexing item: %v", item.Index.Error)
			}
		}
	} else {
		fmt.Println("Bulk indexing completed successfully.")
	}
*/
	// Index each product
    for _, product := range products {

		fmt.Printf("Indexed product %v\n", product)
		
        _, err := client.Index().
            Index("inventory").
            Id(product.ProductID).
            BodyJson(product).
            Do(context.Background())
        if err != nil {
            log.Printf("Error indexing product %s: %v", product.ProductID, err)
        } else {
            fmt.Printf("Indexed product %s\n", product.ProductID)
        }
	}

    fmt.Println("Finished indexing JSON data.")
}
