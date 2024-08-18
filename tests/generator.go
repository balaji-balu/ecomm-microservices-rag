package main

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "os"
    "time"
)

type Product struct {
    ProductID   string  `json:"product_id"`
    ProductName string  `json:"product_name"`
    Category    string  `json:"category"`
    StockLevel  int     `json:"stock_level"`
    Price       float64 `json:"price"`
    Description string  `json:"description"`
}

func generateSampleData(numItems int) []Product {
    categories := []string{"Electronics", "Clothing", "Home & Kitchen", "Books", "Toys"}
    data := make([]Product, numItems)

    rand.Seed(time.Now().UnixNano())
    for i := 0; i < numItems; i++ {
        data[i] = Product{
            ProductID:   fmt.Sprintf("P%d", i+1),
            ProductName: fmt.Sprintf("Product %d", i+1),
            Category:    categories[rand.Intn(len(categories))],
            StockLevel:  rand.Intn(101),
            Price:       float64(rand.Intn(490) + 10), // Price between 10 and 500
            Description: fmt.Sprintf("Description for product %d", i+1),
        }
    }

    return data
}

func main() {
    // Generate 100 sample items
    sampleData := generateSampleData(100)

    // Save to a JSON file
    file, err := os.Create("sample_inventory_data.json")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(sampleData); err != nil {
        fmt.Println("Error encoding JSON:", err)
    }

    fmt.Println("Sample data generated and saved to 'sample_inventory_data.json'.")
}
