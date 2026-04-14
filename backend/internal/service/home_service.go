package service

import (
	"encoding/json"
	"fmt"

	"github.com/ethanqian1990/wechat-mall-backend/internal/config"
	"github.com/ethanqian1990/wechat-mall-backend/internal/model"
	"github.com/ethanqian1990/wechat-mall-backend/internal/repository"
)

type HomeService struct {
	cfg          *config.Config
	homeCfgRepo  *repository.HomeConfigRepository
	productRepo  *repository.ProductRepository
}

func NewHomeService(cfg *config.Config) *HomeService {
	return &HomeService{
		cfg:         cfg,
		homeCfgRepo: repository.NewHomeConfigRepository(),
		productRepo: repository.NewProductRepository(),
	}
}

type HomeBanner struct {
	ImageURL  string `json:"image_url"`
	JumpType  string `json:"jump_type"` // product|category|search|none
	JumpValue string `json:"jump_value"`
}

type FeaturedConfig struct {
	Mode      string `json:"mode"`       // AUTO|MANUAL
	ProductID string `json:"product_id"` // MANUAL 时必填
}

type ListConfig struct {
	ListSize int `json:"list_size"`
}

func (s *HomeService) clampListSize(n int) int {
	minN := s.cfg.Home.MinListSize
	maxN := s.cfg.Home.MaxListSize
	defN := s.cfg.Home.DefaultListSize
	if defN <= 0 {
		defN = 12
	}
	if minN <= 0 {
		minN = 1
	}
	if maxN <= 0 {
		maxN = 50
	}
	if n <= 0 {
		return defN
	}
	if n < minN {
		return minN
	}
	if n > maxN {
		return maxN
	}
	return n
}

func (s *HomeService) GetBanners(regionCode string) ([]HomeBanner, error) {
	cfg, err := s.homeCfgRepo.FindByRegionAndType(regionCode, "banners")
	if err != nil {
		return nil, err
	}
	if cfg == nil || cfg.ConfigValue == "" {
		return []HomeBanner{}, nil
	}
	var banners []HomeBanner
	if err := json.Unmarshal([]byte(cfg.ConfigValue), &banners); err != nil {
		return nil, fmt.Errorf("banner 配置解析失败: %w", err)
	}
	return banners, nil
}

func (s *HomeService) GetFeatured(regionCode, categoryID string) (*model.Product, error) {
	cfg, err := s.homeCfgRepo.FindByRegionAndType(regionCode, "featured")
	if err != nil {
		return nil, err
	}

	var fc FeaturedConfig
	if cfg != nil && cfg.ConfigValue != "" {
		_ = json.Unmarshal([]byte(cfg.ConfigValue), &fc)
	}

	// MANUAL：强制按指定商品
	if fc.Mode == "MANUAL" && fc.ProductID != "" {
		p, err := s.productRepo.FindByID(fc.ProductID)
		if err != nil {
			return nil, err
		}
		// 仍要满足候选池规则：上架、首页展示、可售区域、分类(若传)
		if p == nil {
			return nil, nil
		}
		if p.IsOnShelf != 1 || p.ShowOnHome != 1 {
			return nil, nil
		}
		if categoryID != "" && p.CategoryID != categoryID {
			return nil, nil
		}
		if regionCode != "" {
			candidates, err := s.productRepo.FindHomeCandidates(regionCode, categoryID, 1, nil)
			if err != nil {
				return nil, err
			}
			if len(candidates) == 0 || candidates[0].ID != p.ID {
				// 候选池校验失败：不返回
				return nil, nil
			}
		}
		return p, nil
	}

	// AUTO/默认：销量最高的候选商品
	list, err := s.productRepo.FindHomeCandidates(regionCode, categoryID, 1, nil)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

func (s *HomeService) GetHomeProducts(regionCode, categoryID string) ([]model.Product, int, error) {
	listCfgRec, err := s.homeCfgRepo.FindByRegionAndType(regionCode, "list")
	if err != nil {
		return nil, 0, err
	}
	var lc ListConfig
	if listCfgRec != nil && listCfgRec.ConfigValue != "" {
		_ = json.Unmarshal([]byte(listCfgRec.ConfigValue), &lc)
	}
	n := s.clampListSize(lc.ListSize)

	featured, err := s.GetFeatured(regionCode, categoryID)
	if err != nil {
		return nil, 0, err
	}
	exclude := []string{}
	if featured != nil {
		exclude = append(exclude, featured.ID)
	}

	list, err := s.productRepo.FindHomeCandidates(regionCode, categoryID, n, exclude)
	if err != nil {
		return nil, 0, err
	}
	return list, n, nil
}

func (s *HomeService) UpdateRegionBanners(regionCode string, banners []HomeBanner) error {
	bs, err := json.Marshal(banners)
	if err != nil {
		return err
	}
	return s.homeCfgRepo.Upsert(regionCode, "banners", string(bs))
}

func (s *HomeService) UpdateFeatured(regionCode string, fc FeaturedConfig) error {
	bs, err := json.Marshal(fc)
	if err != nil {
		return err
	}
	return s.homeCfgRepo.Upsert(regionCode, "featured", string(bs))
}

func (s *HomeService) UpdateList(regionCode string, lc ListConfig) error {
	bs, err := json.Marshal(lc)
	if err != nil {
		return err
	}
	return s.homeCfgRepo.Upsert(regionCode, "list", string(bs))
}

type HomeAllConfig struct {
	RegionCode string         `json:"region_code"`
	Banners    []HomeBanner   `json:"banners"`
	Featured   FeaturedConfig `json:"featured"`
	List       ListConfig     `json:"list"`
}

func (s *HomeService) GetAllConfig(regionCode string) (*HomeAllConfig, error) {
	banners, err := s.GetBanners(regionCode)
	if err != nil {
		return nil, err
	}

	featuredRec, err := s.homeCfgRepo.FindByRegionAndType(regionCode, "featured")
	if err != nil {
		return nil, err
	}
	var fc FeaturedConfig
	if featuredRec != nil && featuredRec.ConfigValue != "" {
		_ = json.Unmarshal([]byte(featuredRec.ConfigValue), &fc)
	}

	listRec, err := s.homeCfgRepo.FindByRegionAndType(regionCode, "list")
	if err != nil {
		return nil, err
	}
	var lc ListConfig
	if listRec != nil && listRec.ConfigValue != "" {
		_ = json.Unmarshal([]byte(listRec.ConfigValue), &lc)
	}

	return &HomeAllConfig{
		RegionCode: regionCode,
		Banners:    banners,
		Featured:   fc,
		List:       lc,
	}, nil
}

