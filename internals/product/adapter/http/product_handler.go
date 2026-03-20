package product_http

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	product_dto "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/http/dto"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/product/port"
	pkg_error_response "github.com/premwitthawas/demo_ecommerce_api/pkgs/error_handler"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type productHttpHandler struct {
	tp             trace.Tracer
	productUsecase port.ProductUsecase
}

type ProductHttpHandlerEventtye string

const (
	Created ProductHttpHandlerEventtye = "handler.product.create"
	GetByID ProductHttpHandlerEventtye = "handler.product.get_by_id"
)

func NewHttpHandler(tp trace.Tracer, productUsecase port.ProductUsecase) *productHttpHandler {
	return &productHttpHandler{
		tp:             tp,
		productUsecase: productUsecase,
	}
}

func (h *productHttpHandler) CreateProduct(c fiber.Ctx) error {
	ctx, sp := h.tp.Start(c.Context(), string(Created))
	defer sp.End()
	req := new(product_dto.ProductCreateReq)
	if err := c.Bind().JSON(req); err != nil {
		sp.RecordError(err)
		sp.SetStatus(codes.Error, err.Error())
		return c.Status(400).JSON(&pkg_error_response.ErrorResponse{
			Status:  400,
			Message: "body request error.",
		})
	}
	res, err := h.productUsecase.CreateProduct(ctx, &port.ProductCreateDTO{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
	})
	if err != nil {
		sp.RecordError(err)
		sp.SetStatus(codes.Error, err.Error())
		return product_dto.ProductTranslateError(c, err)
	}
	return c.Status(201).JSON(product_dto.TranslateProduct(res))
}

func (h *productHttpHandler) GetProductByID(c fiber.Ctx) error {
	ctx, sp := h.tp.Start(c.Context(), string(GetByID))
	defer sp.End()
	id := c.Params("id")
	if id == "" {
		sp.RecordError(fmt.Errorf("parameter id is invalid or empty"))
		sp.SetStatus(codes.Error, "parameter id is invalid or empty")
		return c.Status(400).JSON(&pkg_error_response.ErrorResponse{
			Status:  400,
			Message: "parameter id is invalid or empty",
		})
	}
	res, err := h.productUsecase.GetProductByID(ctx, id)
	if err != nil {
		sp.RecordError(err)
		sp.SetStatus(codes.Error, err.Error())
		return product_dto.ProductTranslateError(c, err)
	}
	return c.Status(200).JSON(product_dto.TranslateProduct(res))
}

func (h *productHttpHandler) DeleteProductByID(c fiber.Ctx) error {
	ctx, sp := h.tp.Start(c.Context(), string(GetByID))
	defer sp.End()
	id := c.Params("id")
	if id == "" {
		sp.RecordError(fmt.Errorf("parameter id is invalid or empty"))
		sp.SetStatus(codes.Error, "parameter id is invalid or empty")
		return c.Status(400).JSON(&pkg_error_response.ErrorResponse{
			Status:  400,
			Message: "parameter id is invalid or empty",
		})
	}
	res, err := h.productUsecase.DeleteProductByID(ctx, id)
	if err != nil {
		sp.RecordError(err)
		sp.SetStatus(codes.Error, err.Error())
		return product_dto.ProductTranslateError(c, err)
	}
	return c.Status(200).JSON(product_dto.TranslateProduct(res))
}
