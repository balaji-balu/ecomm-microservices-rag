package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/olivere/elastic/v7"
    "github.com/go-resty/resty/v2"

    // gptcache "github.com/zilliztech/gptcache"
    // "github.com/zilliztech/gptcache/local"
    // "github.com/zilliztech/gptcache/local/storage/sqlite"
)

var (
    esClient *elastic.Client
    apiKey    = "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"
//    cache      gptcache.Cache
)

func init() {
    // Elasticsearch address and configuration
    esAddress := "http://localhost:9200" // Update this if your Elasticsearch is hosted elsewhere

    var err error
    esClient, err = elastic.NewClient(
        elastic.SetBasicAuth("elastic", "yourpassword"),
        elastic.SetURL(esAddress),
        elastic.SetSniff(false),
        elastic.SetHealthcheck(false),
    )
    if err != nil {
        log.Fatalf("Error creating Elasticsearch client: %v", err)
    }

    // Initialize GPTCache using SQLite for storage
    //store := sqlite.NewStorage("gptcache.db") // Using SQLite as storage
    //cache = local.NewCache(store)
}

func retrieveInventoryData(query string) ([]map[string]interface{}, error) {

    // SingleMatchQuery
    // SingleMatchQuery := elastic.NewMatchQuery("category", query)
    //
    // MultiMatchQuery to search across multiple fields with fuzziness
    multiMatchQuery := elastic.NewMultiMatchQuery(query, "category", "product_name", "description").
        Type("best_fields").      // Choose the most relevant fields
        Fuzziness("AUTO").        // Enable fuzziness for handling spelling mistakes
        Operator("and").           // Ensure that all terms must match
        MinimumShouldMatch("75%") // Ensure at least 75% of the terms match

    res, err := esClient.Search().
        Index("inventory").
        Query(multiMatchQuery).
        Size(100).
        Do(context.Background())
    if err != nil {
        return nil, err
    }

    fmt.Printf("res: %v", res)
    var results []map[string]interface{}
    for _, hit := range res.Hits.Hits {
        var source map[string]interface{}
        if err := json.Unmarshal(hit.Source, &source); err != nil {
            return nil, err
        }
        results = append(results, source)
    }
    fmt.Printf("Inventory data:%v", results)
    return results, nil
}

func generateSummary(inventoryData []map[string]interface{}) (string, error) {
    prompt := "Summarize the following inventory data:\n"
    for _, item := range inventoryData {
        prompt += fmt.Sprintf("Product ID: %v, Stock Level: %v\n", item["product_id"], item["stock_level"])
    }

    fmt.Printf("Prompt: %v", prompt)

    // Check cache for prompt
    // cachedSummary, err := cache.Get(prompt)
    // if err == nil && cachedSummary != "" {
    //     fmt.Println("Cache hit!")
    //     return cachedSummary.(string), nil
    // }

    client := resty.New()
    resp, err := client.R().
        SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)).
        SetHeader("Content-Type", "application/json").
        SetBody(map[string]interface{}{
            "model":"gpt-3.5-turbo-instruct", //"text-davinci-003",
            "prompt": prompt,
            "max_tokens": 100,
        }).
        Post("https://api.openai.com/v1/completions")

    if err != nil {
        return "", err
    }

    fmt.Printf("openai results: %v", resp)

    var result map[string]interface{}
    if err := json.Unmarshal(resp.Body(), &result); err != nil {
        return "", err
    }

    //
    //TBD: error handling if openai returns error
    //
    choices := result["choices"].([]interface{})
    if len(choices) > 0 {
        summary := choices[0].(map[string]interface{})["text"].(string)

        // // Store in cache
        // err = cache.Set(prompt, summary)
        // if err != nil {
        //     fmt.Printf("Failed to store in cache: %v\n", err)
        // }
        // fmt.Println("Cache miss, storing result!")

        return summary, nil
    }

    return "", fmt.Errorf("no choices found")
}

func main() {
    query := "clothing" //"high demand"
    inventoryData, err := retrieveInventoryData(query)
    if err != nil {
        log.Fatalf("Error retrieving inventory data: %v", err)
    }

    summary, err := generateSummary(inventoryData)
    if err != nil {
        log.Fatalf("Error generating summary: %v", err)
    }

    fmt.Println(summary)
}
