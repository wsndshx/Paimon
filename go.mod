module github.com/wsndshx/Paimon

go 1.16

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/robfig/cron/v3 v3.0.1
	github.com/wit-ai/wit-go/v2 v2.0.2
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/wit-ai/wit-go/v2 v2.0.2 => github.com/wsndshx/wit-go/v2 v2.0.3
