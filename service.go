package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// DataService handles loading and providing access to geographic data
type DataService struct {
	mu            sync.RWMutex
	geographies   []Geography
	provinces     []Province
	districts     []District
	subDistricts  []SubDistrict
	
	// Index maps for fast lookups
	geographyMap   map[int]Geography
	provinceMap    map[int]Province
	districtMap    map[int]District
	subDistrictMap map[int]SubDistrict
	
	// Relationship indexes
	provincesByGeography    map[int][]Province
	districtsByProvince     map[int][]District
	subDistrictsByDistrict  map[int][]SubDistrict
}

// NewDataService creates a new DataService and loads data from JSON files
func NewDataService(dataPath string) (*DataService, error) {
	ds := &DataService{
		geographyMap:            make(map[int]Geography),
		provinceMap:             make(map[int]Province),
		districtMap:             make(map[int]District),
		subDistrictMap:          make(map[int]SubDistrict),
		provincesByGeography:    make(map[int][]Province),
		districtsByProvince:     make(map[int][]District),
		subDistrictsByDistrict:  make(map[int][]SubDistrict),
	}

	if err := ds.loadData(dataPath); err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return ds, nil
}

// loadData loads all geographic data from JSON files
func (ds *DataService) loadData(dataPath string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// Load geographies
	if err := ds.loadGeographies(filepath.Join(dataPath, "geographies.json")); err != nil {
		return fmt.Errorf("failed to load geographies: %w", err)
	}

	// Load provinces
	if err := ds.loadProvinces(filepath.Join(dataPath, "provinces.json")); err != nil {
		return fmt.Errorf("failed to load provinces: %w", err)
	}

	// Load districts
	if err := ds.loadDistricts(filepath.Join(dataPath, "districts.json")); err != nil {
		return fmt.Errorf("failed to load districts: %w", err)
	}

	// Load sub-districts
	if err := ds.loadSubDistricts(filepath.Join(dataPath, "sub_districts.json")); err != nil {
		return fmt.Errorf("failed to load sub-districts: %w", err)
	}

	// Build indexes
	ds.buildIndexes()

	return nil
}

// loadGeographies loads geography data from JSON file
func (ds *DataService) loadGeographies(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &ds.geographies); err != nil {
		return err
	}

	return nil
}

// loadProvinces loads province data from JSON file
func (ds *DataService) loadProvinces(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &ds.provinces); err != nil {
		return err
	}

	return nil
}

// loadDistricts loads district data from JSON file
func (ds *DataService) loadDistricts(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &ds.districts); err != nil {
		return err
	}

	return nil
}

// loadSubDistricts loads sub-district data from JSON file
func (ds *DataService) loadSubDistricts(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &ds.subDistricts); err != nil {
		return err
	}

	return nil
}

// buildIndexes creates index maps for fast lookups and relationships
func (ds *DataService) buildIndexes() {
	// Build geography map
	for _, geography := range ds.geographies {
		ds.geographyMap[geography.ID] = geography
	}

	// Build province map and geography relationships
	for _, province := range ds.provinces {
		ds.provinceMap[province.ID] = province
		ds.provincesByGeography[province.GeographyID] = append(
			ds.provincesByGeography[province.GeographyID], province)
	}

	// Build district map and province relationships
	for _, district := range ds.districts {
		ds.districtMap[district.ID] = district
		ds.districtsByProvince[district.ProvinceID] = append(
			ds.districtsByProvince[district.ProvinceID], district)
	}

	// Build sub-district map and district relationships
	for _, subDistrict := range ds.subDistricts {
		ds.subDistrictMap[subDistrict.ID] = subDistrict
		ds.subDistrictsByDistrict[subDistrict.DistrictID] = append(
			ds.subDistrictsByDistrict[subDistrict.DistrictID], subDistrict)
	}
}

// GetGeographies returns all geographies
func (ds *DataService) GetGeographies() []Geography {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.geographies
}

// GetGeography returns a geography by ID
func (ds *DataService) GetGeography(id int) (Geography, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	geography, exists := ds.geographyMap[id]
	return geography, exists
}

// GetProvinces returns all provinces
func (ds *DataService) GetProvinces() []Province {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.provinces
}

// GetProvince returns a province by ID
func (ds *DataService) GetProvince(id int) (Province, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	province, exists := ds.provinceMap[id]
	return province, exists
}

// GetProvincesByGeography returns provinces by geography ID
func (ds *DataService) GetProvincesByGeography(geographyID int) []Province {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.provincesByGeography[geographyID]
}

// GetDistricts returns all districts
func (ds *DataService) GetDistricts() []District {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.districts
}

// GetDistrict returns a district by ID
func (ds *DataService) GetDistrict(id int) (District, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	district, exists := ds.districtMap[id]
	return district, exists
}

// GetDistrictsByProvince returns districts by province ID
func (ds *DataService) GetDistrictsByProvince(provinceID int) []District {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.districtsByProvince[provinceID]
}

// GetSubDistricts returns all sub-districts
func (ds *DataService) GetSubDistricts() []SubDistrict {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.subDistricts
}

// GetSubDistrict returns a sub-district by ID
func (ds *DataService) GetSubDistrict(id int) (SubDistrict, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	subDistrict, exists := ds.subDistrictMap[id]
	return subDistrict, exists
}

// GetSubDistrictsByDistrict returns sub-districts by district ID
func (ds *DataService) GetSubDistrictsByDistrict(districtID int) []SubDistrict {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.subDistrictsByDistrict[districtID]
}