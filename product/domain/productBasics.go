package domain

import (
	"fmt"
	"time"
)

const (
	MediaUsageList   = "list"
	MediaUsageDetail = "detail"
)

type (
	// BasicProduct interface - Need to be implements by all Product Types!
	BasicProduct interface {
		BaseData() BasicProductData
		TeaserData() TeaserData
		//IsSaleable - indicates if that product type can be purchased
		IsSaleable() bool
		SaleableData() Saleable
		Type() string
		GetIdentifier() string
		HasMedia(group string, usage string) bool
		GetMedia(group string, usage string) Media
	}

	// BasicProductData is the basic product model
	BasicProductData struct {
		Title            string
		Attributes       Attributes
		ShortDescription string
		Description      string
		Media            []Media

		MarketPlaceCode string
		RetailerCode    string
		RetailerSku     string

		CreatedAt   time.Time
		UpdatedAt   time.Time
		VisibleFrom time.Time
		VisibleTo   time.Time

		CategoryPath          []string
		CategoryCodes         []string
		CategoryToCodeMapping []string

		Keywords []string
		IsNew    bool
	}

	// Saleable are properties required for being selled
	Saleable struct {
		IsSaleable      bool
		SaleableFrom    time.Time
		SaleableTo      time.Time
		ActivePrice     PriceInfo
		AvailablePrices []PriceInfo
	}

	// PriceInfo holds product price information
	PriceInfo struct {
		Default           float64
		Discounted        float64
		DiscountText      string
		Currency          string
		ActiveBase        float64
		ActiveBaseAmount  float64
		ActiveBaseUnit    string
		IsDiscounted      bool
		CampaignRules     []string
		DenyMoreDiscounts bool
		Context           PriceContext
		TaxClass          string
	}

	// PriceContext defines the scope in which the price was calculated
	PriceContext struct {
		CustomerGroup string
		ChannelCode   string
		Locale        string
	}

	// TeaserData is the teaser-information for product previews
	TeaserData struct {
		ShortTitle       string
		ShortDescription string
		//TeaserPrice is the price that should be shown in teasers (listview)
		TeaserPrice PriceInfo
		//TeaserPriceIsFromPrice - is set to true in cases where a product might have different prices (e.g. configurable)
		TeaserPriceIsFromPrice bool
		//PreSelectedVariantSku - might be set for configurables to give a hint to link to a variant of a configurable (That might be the case if a user filters for an attribute and in the teaser the variant with that attribute is shown)
		PreSelectedVariantSku string
		//Media
		Media []Media
		//The sku that should be used to link from Teasers
		MarketPlaceCode string
	}

	// Media holds product media information
	Media struct {
		Type      string
		MimeType  string
		Usage     string
		Title     string
		Reference string
	}

	// Attributes describe a product attributes map
	Attributes map[string]Attribute

	// Attribute for product attributes
	Attribute struct {
		Code     string
		Label    string
		RawValue interface{}
	}
)

// Value returns the raw value
func (at Attribute) Value() string {
	return fmt.Sprintf("%v", at.RawValue)
}

// IsEnabledValue returns true if the value can be seen as a toogle and is enabled
func (at Attribute) IsEnabledValue() bool {
	if at.RawValue == "Yes" {
		return true
	}
	if at.RawValue == "true" {
		return true
	}
	if at.RawValue == "yes" {
		return true
	}
	if at.RawValue == "1" {
		return true
	}
	if at.RawValue == 1 {
		return true
	}
	return false
}

// IsDisabledValue returns true if the value can be seen as a disable toggle/swicth value
func (at Attribute) IsDisabledValue() bool {
	if at.RawValue == "no" {
		return true
	}
	if at.RawValue == "false" {
		return true
	}
	if at.RawValue == "No" {
		return true
	}
	if at.RawValue == "0" {
		return true
	}
	return false
}

// HasMultipleValues checks for multiple raw values
func (at Attribute) HasMultipleValues() bool {
	_, ok := at.RawValue.([]interface{})
	return ok
}

// Values builds a list of product attribute values
func (at Attribute) Values() []string {
	var result []string
	list, ok := at.RawValue.([]interface{})
	if ok {
		for _, entry := range list {
			result = append(result, fmt.Sprintf("%v", entry))
		}
	}
	return result
}

// HasAttribute check
func (bpd BasicProductData) HasAttribute(key string) bool {
	if _, ok := bpd.Attributes[key]; ok {
		return true
	}
	return false
}

// GetFinalPrice getter for price that should be used in calculations (either discounted or default)
func (p PriceInfo) GetFinalPrice() float64 {
	if p.IsDiscounted {
		return p.Discounted
	}
	return p.Default
}

// GetListMedia returns the product media for listing
func (bpd BasicProductData) GetListMedia() Media {
	return bpd.GetMedia("list")
}

// GetMedia returns the FIRST found product media by usage
func (bpd BasicProductData) GetMedia(usage string) Media {
	var emptyMedia Media
	for _, media := range bpd.Media {
		if media.Usage == usage {
			return media
		}
	}
	return emptyMedia
}

// IsSaleableNow  checks flag and time
func (p Saleable) IsSaleableNow() bool {
	if p.IsSaleable == false {
		return false
	}

	//For some reasons IsZero does not always work - thats why we check for 1970
	if (p.SaleableFrom.IsZero() || p.SaleableFrom.Year() == 1970 || p.SaleableFrom.Before(time.Now())) &&
		(p.SaleableTo.IsZero() || p.SaleableTo.Year() == 1970 || p.SaleableTo.After(time.Now())) {
		return true
	}

	return false
}

func findMediaInProduct(p BasicProduct, group string, usage string) *Media {
	var mediaList []Media
	if group == "teaser" {
		mediaList = p.TeaserData().Media
		for _, media := range mediaList {
			if media.Usage == usage {
				return &media
			}
		}
	}

	mediaList = p.BaseData().Media
	for _, media := range mediaList {
		if media.Usage == usage {
			return &media
		}
	}
	return nil
}
