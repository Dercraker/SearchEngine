package obs

type Event string

const (
	RunStart   Event = "crawl:run:start"
	RunSummary Event = "crawl:run:summary"
	RunEnd     Event = "crawl:run:end"

	URLStart     Event = "crawler:url:start"
	URLEnd       Event = "crawler:url:end"
	URLEndFailed Event = "crawler:url:end.failed"
	URLRetry     Event = "crawler:url:retry"

	QueueUpsert Event = "crawler:queue:upsert"
)
