package handler

import (
	"net/http"
	"os"

	"github.com/Komilov31/sales-tracker/internal/dto"
	_ "github.com/Komilov31/sales-tracker/internal/model"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// GetAllItems godoc
//
//	@Summary		Get all items
//	@Description	Retrieve a list of all items, optionally sorted by specified fields
//	@Tags			items
//	@Produce		json
//
// @Param sort_by query []string false "Sort fields (e.g., date,amount)"
//
//	@Success		200		{array}		dto.ItemWithoutAggregated	"List of items"
//	@Failure		400		{object}	map[string]string			"Invalid query parameters"
//	@Failure		500		{object}	map[string]string			"Internal server error"
//	@Router			/items [get]
func (h *Handler) GetAllItems(c *ginext.Context) {
	sortBy := c.QueryArray("sort_by")

	if err := validateGetParams(sortBy); err != nil {
		zlog.Logger.Error().Msg("invalid query parameter: " + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	getItemsParams := dto.GetItemsParams{SortBy: sortBy}
	items, err := h.service.GetAllItems(h.ctx, getItemsParams)
	if err != nil {
		zlog.Logger.Error().Msg("could not get items: " + err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfylly handled GET request and returned items")
	c.JSON(http.StatusOK, itemsWithoutAggregated(items))
}

// GetAggregated godoc
//
//	@Summary		Get aggregated analytics
//	@Description	Retrieve aggregated statistics for items within a date range
//	@Tags			analytics
//	@Produce		json
//	@Param			from	query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			to		query		string	false	"End date (YYYY-MM-DD)"
//	@Success		200		{array}		model.Item	"Aggregated items"
//	@Failure		400		{object}	map[string]string	"Invalid date parameters"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/analytics [get]
func (h *Handler) GetAggregated(c *ginext.Context) {
	from := c.Query("from")
	to := c.Query("to")

	if from != "" && to != "" {
		if err := validateDate(from, to); err != nil {
			zlog.Logger.Error().Msg(err.Error())
			c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
			return
		}
	}

	items, err := h.service.GetAggregated(h.ctx, from, to)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfylly handled GET request and returned aggregated items")
	c.JSON(http.StatusOK, items)
}

// GetAggregatedCSV godoc
//
//	@Summary		Export aggregated analytics as CSV
//	@Description	Download CSV file with aggregated statistics for a date range
//	@Tags			analytics
//	@Param			from	query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			to		query		string	false	"End date (YYYY-MM-DD)"
//	@Success		200		{file}		application/octet-stream	"aggregated_data.csv"
//	@Failure		400		{object}	map[string]string	"Invalid date parameters"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/analytics/csv [get]
func (h *Handler) GetAggregatedCSV(c *ginext.Context) {
	from := c.Query("from")
	to := c.Query("to")

	if from != "" && to != "" {
		if err := validateDate(from, to); err != nil {
			zlog.Logger.Error().Msg(err.Error())
			c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
			return
		}
	}

	path, err := h.service.CSVAggregated(h.ctx, from, to)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfylly handled GET request and returned csv file with data")
	c.FileAttachment(path, "aggregated_data.csv")

	if err := os.Remove(path); err != nil {
		zlog.Logger.Info().Msg("could not delete temporary csv file: " + err.Error())
	}
}

// GetFilteredCSV godoc
//
//	@Summary		Export filtered items as CSV
//	@Description	Download CSV file with filtered and sorted items
//	@Tags			items
//	@Produce		application/octet-stream
//
// @Param sort_by query []string false "Sort fields (e.g., date,amount)"
// @Success		200		{file}		application/octet-stream	"filtered_data.csv"
// @Failure		400		{object}	map[string]string	"Invalid query parameters"
// @Failure		500		{object}	map[string]string	"Internal server error"
// @Router			/items/csv [get]
func (h *Handler) GetFilteredCSV(c *ginext.Context) {
	sortBy := c.QueryArray("sort_by")

	if err := validateGetParams(sortBy); err != nil {
		zlog.Logger.Error().Msg("invalid query parameter: " + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	getItemsParams := dto.GetItemsParams{SortBy: sortBy}
	path, err := h.service.CSVAllItems(h.ctx, getItemsParams)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfylly handled GET request and returned csv file with data")
	c.FileAttachment(path, "filtered_data.csv")

	if err := os.Remove(path); err != nil {
		zlog.Logger.Info().Msg("could not delete temporary csv file: " + err.Error())
	}
}

// GetMainPage godoc
// @Summary      Get main page
// @Description  Get the main HTML page of the application
// @Tags         pages
// @Accept       json
// @Produce      html
// @Success      200  {string} string "HTML page content"
// @Router       / [get]
func (h *Handler) GetMainPage(c *ginext.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
