package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// LocationHandler handles HTTP requests for location data
type LocationHandler struct {
	dataService *DataService
}

// NewLocationHandler creates a new LocationHandler
func NewLocationHandler(dataService *DataService) *LocationHandler {
	return &LocationHandler{
		dataService: dataService,
	}
}

// GetGeographies returns all geographies
func (h *LocationHandler) GetGeographies(c *fiber.Ctx) error {
	geographies := h.dataService.GetGeographies()
	
	return c.JSON(APIResponse{
		Status: "success",
		Data:   geographies,
	})
}

// GetProvinces returns all provinces with optional filtering
func (h *LocationHandler) GetProvinces(c *fiber.Ctx) error {
	geographyIDStr := c.Query("geography_id")
	search := strings.ToLower(c.Query("search"))
	
	provinces := h.dataService.GetProvinces()
	
	// Filter by geography ID if provided
	if geographyIDStr != "" {
		geographyID, err := strconv.Atoi(geographyIDStr)
		if err != nil {
			return c.Status(400).JSON(APIResponse{
				Status: "error",
				Error:  "Invalid geography_id parameter",
			})
		}
		provinces = h.dataService.GetProvincesByGeography(geographyID)
	}
	
	// Filter by search term if provided
	if search != "" {
		filteredProvinces := make([]Province, 0)
		for _, province := range provinces {
			if strings.Contains(strings.ToLower(province.NameTH), search) ||
				strings.Contains(strings.ToLower(province.NameEN), search) {
				filteredProvinces = append(filteredProvinces, province)
			}
		}
		provinces = filteredProvinces
	}
	
	// Add pagination
	page, limit := getPaginationParams(c)
	paginatedData, pagination := paginate(provinces, page, limit)
	
	return c.JSON(PaginatedResponse{
		Status:     "success",
		Data:       paginatedData,
		Pagination: pagination,
	})
}

// GetProvinceByID returns a specific province by ID
func (h *LocationHandler) GetProvinceByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(APIResponse{
			Status: "error",
			Error:  "Invalid province ID",
		})
	}
	
	province, exists := h.dataService.GetProvince(id)
	if !exists {
		return c.Status(404).JSON(APIResponse{
			Status: "error",
			Error:  "Province not found",
		})
	}
	
	// Include geography information
	geography, _ := h.dataService.GetGeography(province.GeographyID)
	result := ProvinceWithGeography{
		Province:  province,
		Geography: &geography,
	}
	
	return c.JSON(APIResponse{
		Status: "success",
		Data:   result,
	})
}

// GetDistrictsByProvinceID returns districts for a specific province
func (h *LocationHandler) GetDistrictsByProvinceID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	provinceID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(APIResponse{
			Status: "error",
			Error:  "Invalid province ID",
		})
	}
	
	// Check if province exists
	_, exists := h.dataService.GetProvince(provinceID)
	if !exists {
		return c.Status(404).JSON(APIResponse{
			Status: "error",
			Error:  "Province not found",
		})
	}
	
	districts := h.dataService.GetDistrictsByProvince(provinceID)
	search := strings.ToLower(c.Query("search"))
	
	// Filter by search term if provided
	if search != "" {
		filteredDistricts := make([]District, 0)
		for _, district := range districts {
			if strings.Contains(strings.ToLower(district.NameTH), search) ||
				strings.Contains(strings.ToLower(district.NameEN), search) {
				filteredDistricts = append(filteredDistricts, district)
			}
		}
		districts = filteredDistricts
	}
	
	// Add pagination
	page, limit := getPaginationParams(c)
	paginatedData, pagination := paginate(districts, page, limit)
	
	return c.JSON(PaginatedResponse{
		Status:     "success",
		Data:       paginatedData,
		Pagination: pagination,
	})
}

// GetDistricts returns all districts with optional filtering
func (h *LocationHandler) GetDistricts(c *fiber.Ctx) error {
	provinceIDStr := c.Query("province_id")
	search := strings.ToLower(c.Query("search"))
	
	districts := h.dataService.GetDistricts()
	
	// Filter by province ID if provided
	if provinceIDStr != "" {
		provinceID, err := strconv.Atoi(provinceIDStr)
		if err != nil {
			return c.Status(400).JSON(APIResponse{
				Status: "error",
				Error:  "Invalid province_id parameter",
			})
		}
		districts = h.dataService.GetDistrictsByProvince(provinceID)
	}
	
	// Filter by search term if provided
	if search != "" {
		filteredDistricts := make([]District, 0)
		for _, district := range districts {
			if strings.Contains(strings.ToLower(district.NameTH), search) ||
				strings.Contains(strings.ToLower(district.NameEN), search) {
				filteredDistricts = append(filteredDistricts, district)
			}
		}
		districts = filteredDistricts
	}
	
	// Add pagination
	page, limit := getPaginationParams(c)
	paginatedData, pagination := paginate(districts, page, limit)
	
	return c.JSON(PaginatedResponse{
		Status:     "success",
		Data:       paginatedData,
		Pagination: pagination,
	})
}

// GetDistrictByID returns a specific district by ID
func (h *LocationHandler) GetDistrictByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(APIResponse{
			Status: "error",
			Error:  "Invalid district ID",
		})
	}
	
	district, exists := h.dataService.GetDistrict(id)
	if !exists {
		return c.Status(404).JSON(APIResponse{
			Status: "error",
			Error:  "District not found",
		})
	}
	
	// Include province information
	province, _ := h.dataService.GetProvince(district.ProvinceID)
	result := DistrictWithProvince{
		District: district,
		Province: &province,
	}
	
	return c.JSON(APIResponse{
		Status: "success",
		Data:   result,
	})
}

// GetSubDistrictsByDistrictID returns sub-districts for a specific district
func (h *LocationHandler) GetSubDistrictsByDistrictID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	districtID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(APIResponse{
			Status: "error",
			Error:  "Invalid district ID",
		})
	}
	
	// Check if district exists
	_, exists := h.dataService.GetDistrict(districtID)
	if !exists {
		return c.Status(404).JSON(APIResponse{
			Status: "error",
			Error:  "District not found",
		})
	}
	
	subDistricts := h.dataService.GetSubDistrictsByDistrict(districtID)
	search := strings.ToLower(c.Query("search"))
	zipCodeStr := c.Query("zip_code")
	
	// Filter by zip code if provided
	if zipCodeStr != "" {
		zipCode, err := strconv.Atoi(zipCodeStr)
		if err != nil {
			return c.Status(400).JSON(APIResponse{
				Status: "error",
				Error:  "Invalid zip_code parameter",
			})
		}
		filteredSubDistricts := make([]SubDistrict, 0)
		for _, subDistrict := range subDistricts {
			if subDistrict.ZipCode == zipCode {
				filteredSubDistricts = append(filteredSubDistricts, subDistrict)
			}
		}
		subDistricts = filteredSubDistricts
	}
	
	// Filter by search term if provided
	if search != "" {
		filteredSubDistricts := make([]SubDistrict, 0)
		for _, subDistrict := range subDistricts {
			if strings.Contains(strings.ToLower(subDistrict.NameTH), search) ||
				strings.Contains(strings.ToLower(subDistrict.NameEN), search) {
				filteredSubDistricts = append(filteredSubDistricts, subDistrict)
			}
		}
		subDistricts = filteredSubDistricts
	}
	
	// Add pagination
	page, limit := getPaginationParams(c)
	paginatedData, pagination := paginate(subDistricts, page, limit)
	
	return c.JSON(PaginatedResponse{
		Status:     "success",
		Data:       paginatedData,
		Pagination: pagination,
	})
}

// GetSubDistricts returns all sub-districts with optional filtering
func (h *LocationHandler) GetSubDistricts(c *fiber.Ctx) error {
	districtIDStr := c.Query("district_id")
	search := strings.ToLower(c.Query("search"))
	zipCodeStr := c.Query("zip_code")
	
	subDistricts := h.dataService.GetSubDistricts()
	
	// Filter by district ID if provided
	if districtIDStr != "" {
		districtID, err := strconv.Atoi(districtIDStr)
		if err != nil {
			return c.Status(400).JSON(APIResponse{
				Status: "error",
				Error:  "Invalid district_id parameter",
			})
		}
		subDistricts = h.dataService.GetSubDistrictsByDistrict(districtID)
	}
	
	// Filter by zip code if provided
	if zipCodeStr != "" {
		zipCode, err := strconv.Atoi(zipCodeStr)
		if err != nil {
			return c.Status(400).JSON(APIResponse{
				Status: "error",
				Error:  "Invalid zip_code parameter",
			})
		}
		filteredSubDistricts := make([]SubDistrict, 0)
		for _, subDistrict := range subDistricts {
			if subDistrict.ZipCode == zipCode {
				filteredSubDistricts = append(filteredSubDistricts, subDistrict)
			}
		}
		subDistricts = filteredSubDistricts
	}
	
	// Filter by search term if provided
	if search != "" {
		filteredSubDistricts := make([]SubDistrict, 0)
		for _, subDistrict := range subDistricts {
			if strings.Contains(strings.ToLower(subDistrict.NameTH), search) ||
				strings.Contains(strings.ToLower(subDistrict.NameEN), search) {
				filteredSubDistricts = append(filteredSubDistricts, subDistrict)
			}
		}
		subDistricts = filteredSubDistricts
	}
	
	// Add pagination
	page, limit := getPaginationParams(c)
	paginatedData, pagination := paginate(subDistricts, page, limit)
	
	return c.JSON(PaginatedResponse{
		Status:     "success",
		Data:       paginatedData,
		Pagination: pagination,
	})
}

// GetSubDistrictByID returns a specific sub-district by ID
func (h *LocationHandler) GetSubDistrictByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(APIResponse{
			Status: "error",
			Error:  "Invalid sub-district ID",
		})
	}
	
	subDistrict, exists := h.dataService.GetSubDistrict(id)
	if !exists {
		return c.Status(404).JSON(APIResponse{
			Status: "error",
			Error:  "Sub-district not found",
		})
	}
	
	// Include district and province information
	district, _ := h.dataService.GetDistrict(subDistrict.DistrictID)
	province, _ := h.dataService.GetProvince(district.ProvinceID)
	
	result := SubDistrictWithDistrict{
		SubDistrict: subDistrict,
		District:    &district,
		Province:    &province,
	}
	
	return c.JSON(APIResponse{
		Status: "success",
		Data:   result,
	})
}

// Helper functions

// getPaginationParams extracts pagination parameters from query string
func getPaginationParams(c *fiber.Ctx) (page, limit int) {
	page = 1
	limit = 20
	
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	
	return page, limit
}

// paginate applies pagination to a slice of data
func paginate(data interface{}, page, limit int) (interface{}, Pagination) {
	var total int
	var result interface{}
	
	switch v := data.(type) {
	case []Geography:
		total = len(v)
		start := (page - 1) * limit
		end := start + limit
		if start > total {
			result = []Geography{}
		} else {
			if end > total {
				end = total
			}
			result = v[start:end]
		}
	case []Province:
		total = len(v)
		start := (page - 1) * limit
		end := start + limit
		if start > total {
			result = []Province{}
		} else {
			if end > total {
				end = total
			}
			result = v[start:end]
		}
	case []District:
		total = len(v)
		start := (page - 1) * limit
		end := start + limit
		if start > total {
			result = []District{}
		} else {
			if end > total {
				end = total
			}
			result = v[start:end]
		}
	case []SubDistrict:
		total = len(v)
		start := (page - 1) * limit
		end := start + limit
		if start > total {
			result = []SubDistrict{}
		} else {
			if end > total {
				end = total
			}
			result = v[start:end]
		}
	default:
		result = data
	}
	
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}
	
	pagination := Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
	
	return result, pagination
}