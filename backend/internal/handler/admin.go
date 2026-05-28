package handler

import (
	"net/http"
	"strconv"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/service"
	"mini-12306/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	admin *service.AdminService
}

func NewAdminHandler(admin *service.AdminService) *AdminHandler {
	return &AdminHandler{admin: admin}
}

func (h *AdminHandler) PublicStations(c *gin.Context) {
	var query dto.StationListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "站点查询参数不正确")
		return
	}
	result, err := h.admin.PublicStations(query)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) ListStations(c *gin.Context) {
	var query dto.StationListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "站点查询参数不正确")
		return
	}
	result, err := h.admin.ListStations(query)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) CreateStation(c *gin.Context) {
	var req dto.SaveStationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "站点信息不完整")
		return
	}
	result, err := h.admin.CreateStation(req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) UpdateStation(c *gin.Context) {
	id, ok := parseIDParam(c, "stationId", "站点编号不正确")
	if !ok {
		return
	}
	var req dto.SaveStationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "站点信息不完整")
		return
	}
	result, err := h.admin.UpdateStation(id, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) DisableStation(c *gin.Context) {
	id, ok := parseIDParam(c, "stationId", "站点编号不正确")
	if !ok {
		return
	}
	result, err := h.admin.DisableStation(id)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) ListTrains(c *gin.Context) {
	var query dto.AdminTrainListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "车次查询参数不正确")
		return
	}
	result, err := h.admin.ListTrains(query)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) CreateTrain(c *gin.Context) {
	var req dto.SaveTrainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "车次信息不完整")
		return
	}
	result, err := h.admin.CreateTrain(req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) UpdateTrain(c *gin.Context) {
	id, ok := parseIDParam(c, "trainId", "车次编号不正确")
	if !ok {
		return
	}
	var req dto.SaveTrainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "车次信息不完整")
		return
	}
	result, err := h.admin.UpdateTrain(id, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) DeleteTrain(c *gin.Context) {
	id, ok := parseIDParam(c, "trainId", "车次编号不正确")
	if !ok {
		return
	}
	result, err := h.admin.DeleteTrain(id)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) SellableStats(c *gin.Context) {
	var query dto.SellableTrainStatsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "请填写出发站和到达站")
		return
	}
	result, err := h.admin.SellableStats(query)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) ListStops(c *gin.Context) {
	trainID, ok := parseIDParam(c, "trainId", "车次编号不正确")
	if !ok {
		return
	}
	result, err := h.admin.ListStops(trainID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) SaveStops(c *gin.Context) {
	trainID, ok := parseIDParam(c, "trainId", "车次编号不正确")
	if !ok {
		return
	}
	var req dto.SaveTrainStopsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "经停数据不完整")
		return
	}
	result, err := h.admin.SaveStops(trainID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) ListInventories(c *gin.Context) {
	var query dto.InventoryListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "票额查询参数不正确")
		return
	}
	result, err := h.admin.ListInventories(query)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) SaveInventory(c *gin.Context) {
	var req dto.SaveInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "票额信息不完整")
		return
	}
	result, err := h.admin.SaveInventory(req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) QuoteStats(c *gin.Context) {
	var query dto.InventoryQuoteStatsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "请填写车次编号")
		return
	}
	result, err := h.admin.QuoteStats(query)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) FlowInventory(c *gin.Context) {
	var req dto.InventoryFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "票额流转信息不完整")
		return
	}
	result, err := h.admin.FlowInventory(req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func parseIDParam(c *gin.Context, name string, message string) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, message)
		return 0, false
	}
	return id, true
}
