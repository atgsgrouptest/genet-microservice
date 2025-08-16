package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
)

type QdrantPoint struct {
	ID      int                 `json:"id"`
	Vector  []float64              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

type QdrantUpsertRequest struct {
	Points []QdrantPoint `json:"points"`
}
func StoreInQdrant(docs []EmbeddedDocument, collectionName string) error {
	const batchSize = 500 // adjust based on your setup
	for i := 0; i < len(docs); i += batchSize {
		end := i + batchSize
		if end > len(docs) {
			end = len(docs)
		}
		batch := docs[i:end]

		var points []QdrantPoint
		for j, doc := range batch {
			points = append(points, QdrantPoint{
				ID:     i + j, // keep unique ID across batches
				Vector: doc.Embedding,
				Payload: map[string]interface{}{
					"text": doc.Content,
				},
			})
		}

		body := QdrantUpsertRequest{Points: points}
		bodyData, _ := json.Marshal(body)

		url := fmt.Sprintf("http://localhost:6333/collections/%s/points?wait=true", collectionName)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyData))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to insert points, status: %s, body: %s", resp.Status, string(b))
		}
		fmt.Printf("Inserted batch %d–%d into Qdrant\n", i, end)
	}
	return nil
}
