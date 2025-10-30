package main

import (
	"time"
)

// Geography represents a geographic region in Thailand
type Geography struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Province represents a province in Thailand
type Province struct {
	ID          int        `json:"id"`
	NameTH      string     `json:"name_th"`
	NameEN      string     `json:"name_en"`
	GeographyID int        `json:"geography_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// District represents a district (amphoe) in Thailand
type District struct {
	ID         int        `json:"id"`
	NameTH     string     `json:"name_th"`
	NameEN     string     `json:"name_en"`
	ProvinceID int        `json:"province_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

// SubDistrict represents a sub-district (tambon) in Thailand
type SubDistrict struct {
	ID         int        `json:"id"`
	ZipCode    int        `json:"zip_code"`
	NameTH     string     `json:"name_th"`
	NameEN     string     `json:"name_en"`
	DistrictID int        `json:"district_id"`
	Lat        *float64   `json:"lat"`
	Long       *float64   `json:"long"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

// API Response structures
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Extended structures with relationships
type ProvinceWithGeography struct {
	Province
	Geography *Geography `json:"geography,omitempty"`
}

type DistrictWithProvince struct {
	District
	Province *Province `json:"province,omitempty"`
}

type SubDistrictWithDistrict struct {
	SubDistrict
	District *District `json:"district,omitempty"`
	Province *Province `json:"province,omitempty"`
}