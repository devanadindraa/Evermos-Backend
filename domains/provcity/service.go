package provcity

import (
	"fmt"

	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/go-resty/resty/v2"
)

type Provcity interface {
	GetListProvince() ([]Province, error)
	GetListCity(provID string) ([]City, error)
	GetProvinceByID(provID string) (*Province, error)
	GetCityByID(cityID string) (*City, error)
}

type provcity struct {
	resty *resty.Client
}

func NewEmsiaClient(conf *config.Config) Provcity {
	EmsifaConf := conf.Emsifa
	c := resty.New().
		SetBaseURL(EmsifaConf.BaseUrl)

	return &provcity{
		resty: c,
	}
}

func (c *provcity) GetListProvince() ([]Province, error) {
	var res []Province

	resp, err := c.resty.R().
		SetResult(&res).
		Get("/provinces.json")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, err
	}

	return res, nil
}

func (c *provcity) GetListCity(provID string) ([]City, error) {
	var res []City

	url := fmt.Sprintf("/regencies/%s.json", provID)

	resp, err := c.resty.R().
		SetResult(&res).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, err
	}

	return res, nil
}

func (c *provcity) GetProvinceByID(provID string) (*Province, error) {
	var provinces []Province

	resp, err := c.resty.R().
		SetResult(&provinces).
		Get("/provinces.json")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to fetch provinces")
	}

	for _, prov := range provinces {
		if prov.ID == provID {
			return &prov, nil
		}
	}

	return nil, fmt.Errorf("province with id %s not found", provID)
}

func (c *provcity) GetCityByID(cityID string) (*City, error) {
	var provinces []Province

	_, err := c.resty.R().
		SetResult(&provinces).
		Get("/provinces.json")
	if err != nil {
		return nil, err
	}

	for _, prov := range provinces {
		var cities []City
		url := fmt.Sprintf("/regencies/%s.json", prov.ID)
		_, err := c.resty.R().
			SetResult(&cities).
			Get(url)
		if err != nil {
			continue
		}
		for _, city := range cities {
			if city.ID == cityID {
				return &city, nil
			}
		}
	}

	return nil, fmt.Errorf("city with id %s not found", cityID)
}
