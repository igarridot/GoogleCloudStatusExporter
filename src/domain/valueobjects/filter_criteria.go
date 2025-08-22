package valueobjects

import "strings"

type FilterCriteria struct {
	Zones    []string
	Products []string
}

func NewFilterCriteria(zones, products string) FilterCriteria {
	var zoneList, productList []string

	if zones != "" {
		zoneList = strings.Split(zones, ",")
	}

	if products != "" {
		productList = strings.Split(products, ",")
	}

	return FilterCriteria{
		Zones:    zoneList,
		Products: productList,
	}
}

func (fc FilterCriteria) HasZoneFilter() bool {
	return len(fc.Zones) > 0
}

func (fc FilterCriteria) HasProductFilter() bool {
	return len(fc.Products) > 0
}
