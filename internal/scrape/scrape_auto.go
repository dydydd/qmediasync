package scrape

import (
	"Q115-STRM/internal/baidupan"
	"Q115-STRM/internal/helpers"
	"Q115-STRM/internal/models"
	"Q115-STRM/internal/openlist"
	"Q115-STRM/internal/v115open"
	"context"
)

// autoScrapeImpl 自动识别模式的刮削实现
// 目录媒体类型为 auto 时，扫描阶段已按文件判断出每个文件是电影还是剧集，
// 这里将单个文件分发给对应的电影或电视剧实现处理
type autoScrapeImpl struct {
	movieScrape  *movieScrapeImpl
	tvshowScrape *tvShowScrapeImpl
}

func NewAutoScrapeImpl(scrapePath *models.ScrapePath, ctx context.Context, v115Client *v115open.OpenClient, openlistClient *openlist.Client, baiduPanClient *baidupan.Client) scrapeImpl {
	return &autoScrapeImpl{
		movieScrape:  NewMovieScrapeImpl(scrapePath, ctx, v115Client, openlistClient, baiduPanClient).(*movieScrapeImpl),
		tvshowScrape: NewTvShowScrapeImpl(scrapePath, ctx, v115Client, openlistClient, baiduPanClient).(*tvShowScrapeImpl),
	}
}

// 先处理电影，再处理电视剧
func (a *autoScrapeImpl) Start() error {
	helpers.AppLogger.Infof("开始自动识别模式刮削：先处理电影，再处理电视剧")
	if err := a.movieScrape.Start(); err != nil {
		helpers.AppLogger.Errorf("自动识别模式电影刮削失败: %v", err)
		return err
	}
	helpers.AppLogger.Infof("电影刮削完成，开始处理电视剧")
	if err := a.tvshowScrape.Start(); err != nil {
		helpers.AppLogger.Errorf("自动识别模式电视剧刮削失败: %v", err)
		return err
	}
	helpers.AppLogger.Infof("自动识别模式刮削全部完成")
	return nil
}

func (a *autoScrapeImpl) Process(mediaFile *models.ScrapeMediaFile) error {
	if mediaFile.MediaType == models.MediaTypeTvShow {
		return a.tvshowScrape.Process(mediaFile)
	}
	return a.movieScrape.Process(mediaFile)
}

func (a *autoScrapeImpl) Scrape(mediaFile *models.ScrapeMediaFile) error {
	if mediaFile.MediaType == models.MediaTypeTvShow {
		return a.tvshowScrape.Scrape(mediaFile)
	}
	return a.movieScrape.Scrape(mediaFile)
}

func (a *autoScrapeImpl) Rollback(mediaFile *models.ScrapeMediaFile) error {
	if mediaFile.MediaType == models.MediaTypeTvShow {
		return a.tvshowScrape.Rollback(mediaFile)
	}
	return a.movieScrape.Rollback(mediaFile)
}
