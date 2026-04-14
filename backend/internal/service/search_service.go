package service

import (
	"github.com/ethanqian1990/wechat-mall-backend/internal/repository"
)

type SearchService struct {
	historyRepo *repository.SearchHistoryRepository
	productRepo *repository.ProductRepository
}

func NewSearchService() *SearchService {
	return &SearchService{
		historyRepo: repository.NewSearchHistoryRepository(),
		productRepo: repository.NewProductRepository(),
	}
}

func (s *SearchService) ListHistory(userID string) ([]string, error) {
	histories, err := s.historyRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(histories))
	for _, h := range histories {
		out = append(out, h.Keyword)
	}
	return out, nil
}

func (s *SearchService) ClearHistory(userID string) error {
	return s.historyRepo.DeleteByUserID(userID)
}

func (s *SearchService) Suggest(keyword string, limit int) ([]string, error) {
	if keyword == "" {
		return []string{}, nil
	}
	return s.productRepo.SuggestNamesByPrefix(keyword, limit)
}

func (s *SearchService) Hot(limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.historyRepo.FindHotKeywords(limit)
}
