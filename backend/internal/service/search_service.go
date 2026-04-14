package service

import (
	"github.com/ethanqian1990/wechat-mall-backend/internal/repository"
)

type SearchService struct {
	historyRepo *repository.SearchHistoryRepository
}

func NewSearchService() *SearchService {
	return &SearchService{
		historyRepo: repository.NewSearchHistoryRepository(),
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
