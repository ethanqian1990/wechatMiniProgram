package service

import (
	wechatmall "github.com/ethanqian1990/wechat-mall-backend/internal/model"
	wechatrepo "github.com/ethanqian1990/wechat-mall-backend/internal/repository"
	"github.com/google/uuid"
)

type SearchService struct {
	searchHistoryRepo *wechatrepo.SearchHistoryRepository
}

func NewSearchService() *SearchService {
	return &SearchService{
		searchHistoryRepo: wechatrepo.NewSearchHistoryRepository(),
	}
}

// SaveHistory 保存搜索历史
func (s *SearchService) SaveHistory(userID, keyword string) error {
	if keyword == "" {
		return nil
	}
	
	// 先删除同用户同关键词的旧记录
	s.searchHistoryRepo.DeleteByUserAndKeyword(userID, keyword)
	
	// 保存新记录
	history := &wechatmall.SearchHistory{
		ID:      "sh_" + uuid.New().String()[:8],
		UserID:  userID,
		Keyword: keyword,
	}
	return s.searchHistoryRepo.Create(history)
}

// GetHistory 获取用户搜索历史
func (s *SearchService) GetHistory(userID string) ([]wechatmall.SearchHistory, error) {
	return s.searchHistoryRepo.FindByUserID(userID)
}

// ClearHistory 清除用户搜索历史
func (s *SearchService) ClearHistory(userID string) error {
	return s.searchHistoryRepo.DeleteByUserID(userID)
}

// GetHotKeywords 获取热门搜索词
func (s *SearchService) GetHotKeywords(limit int) ([]string, error) {
	return s.searchHistoryRepo.FindHotKeywords(limit)
}
