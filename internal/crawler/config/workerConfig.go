package config

import (
	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
	"github.com/google/uuid"
)

type WorkerConfig struct {
	CrawlerInstanceID string
	CrawlerWorker     int
}

func LoadWorkerConfig() (workerConfig WorkerConfig, err error) {
	instanceId := configHelper.GetEnv("CRAWLER_INSTANCE_ID", uuid.NewString())
	workerCount := configHelper.GetEnvInt("CRAWLER_WORKERS", 1)

	return WorkerConfig{
		CrawlerInstanceID: instanceId,
		CrawlerWorker:     workerCount,
	}, nil
}
