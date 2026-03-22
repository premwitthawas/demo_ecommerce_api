package searc_product_handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/search/port"
	"go.opentelemetry.io/otel/trace"
)

type searchProductHandler struct {
	tp      trace.Tracer
	usecase port.SearchUsecase
}

func NewSearchProductHandler(tp trace.Tracer,
	usecase port.SearchUsecase) *searchProductHandler {
	return &searchProductHandler{
		tp:      tp,
		usecase: usecase,
	}
}

type response struct {
	Data  any `json:"data"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}
type responseErr struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (h *searchProductHandler) GetProducts(c fiber.Ctx) error {
	ctx, sp := h.tp.Start(c.Context(), "handler.search.getProducts")
	defer sp.End()
	q := c.Query("q", "")
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	products, total, err := h.usecase.GetProducts(ctx, q, page, limit)
	if err != nil {
		sp.RecordError(err)
		return c.Status(500).JSON(&responseErr{
			Status:  500,
			Message: "ailed to fetch products",
		})
	}
	return c.Status(200).JSON(&response{
		Data:  products,
		Page:  page,
		Limit: limit,
		Total: total,
	})
}
