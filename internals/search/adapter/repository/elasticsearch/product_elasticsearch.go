package search_elasticsearch_repository

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/elastic/go-elasticsearch/v9/esutil"
	search "github.com/premwitthawas/demo_ecommerce_api/internals/search/domain/product"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/search/port"
	"go.opentelemetry.io/otel/trace"
)

type searchElasticSearchRepository struct {
	client *elasticsearch.Client
	tp     trace.Tracer
}

func Ptr[T any](v T) *T {
	return &v
}

// GetProducts implements search.SearchRepository.
func (s *searchElasticSearchRepository) GetProducts(ctx context.Context, q string, page int, limit int) ([]*search.SearchProductMessage, int, error) {
	ctx, sp := s.tp.Start(ctx, "repository.search.GetProducts")
	defer sp.End()
	from := (page - 1) * limit
	var queryMap map[string]interface{}
	if q == "" {
		queryMap = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	} else {
		queryMap = map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     q,
				"fields":    []string{"name", "description"},
				"fuzziness": "AUTO",
			},
		}
	}
	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": queryMap,
	}
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, 0, err
	}
	req := esapi.SearchRequest{
		Index:          []string{"products"},
		Body:           &buf,
		From:           Ptr(from),
		Size:           Ptr(limit),
		TrackTotalHits: true,
	}
	res, err := req.Do(ctx, s.client)
	if err != nil {
		sp.RecordError(err)
		return nil, 0, err
	}
	defer res.Body.Close()
	if res.IsError() {
		sp.RecordError(err)
		return nil, 0, err
	}

	var resBody struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source search.SearchProductMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		sp.RecordError(err)
		return nil, 0, err
	}
	products := make([]*search.SearchProductMessage, 0, len(resBody.Hits.Hits))
	for _, hit := range resBody.Hits.Hits {
		p := hit.Source
		products = append(products, &p)
	}
	return products, resBody.Hits.Total.Value, nil
}

func (s *searchElasticSearchRepository) DeleteProduct(ctx context.Context, productID string) error {
	ctx, sp := s.tp.Start(ctx, "repository.search.delete_product")
	defer sp.End()
	req := esapi.DeleteRequest{
		Index:      "products",
		DocumentID: productID,
		Refresh:    "true",
	}
	res, err := req.Do(ctx, s.client)
	if err != nil {
		sp.RecordError(err)
		return err
	}
	if res.IsError() {
		sp.RecordError(err)
		return err
	}
	defer res.Body.Close()
	return nil
}

func (s *searchElasticSearchRepository) UpsertProduct(ctx context.Context, doc *search.SearchProductMessage) error {
	ctx, sp := s.tp.Start(ctx, "repository.search.upsert_product")
	defer sp.End()
	req := esapi.IndexRequest{
		Index:      "products",
		DocumentID: doc.ID,
		Body:       esutil.NewJSONReader(doc),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, s.client)
	if err != nil {
		sp.RecordError(err)
		return err
	}
	if res.IsError() {
		sp.RecordError(err)
		return err
	}
	defer res.Body.Close()
	return nil
}

func NewSearchElasticSearchRepository(client *elasticsearch.Client, tp trace.Tracer) port.SearchRepository {
	return &searchElasticSearchRepository{
		client: client,
		tp:     tp,
	}
}
