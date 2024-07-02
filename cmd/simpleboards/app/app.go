package app

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/posilva/simpleboards/cmd/simpleboards/config"
	"github.com/posilva/simpleboards/internal/adapters/input/handler"
	"github.com/posilva/simpleboards/internal/adapters/output/configprovider"
	"github.com/posilva/simpleboards/internal/adapters/output/logging"
	"github.com/posilva/simpleboards/internal/adapters/output/repository"
	"github.com/posilva/simpleboards/internal/adapters/output/scoreboard"
	"github.com/posilva/simpleboards/internal/core/services"
)

func Run() {
	r := gin.Default()

	service, err := createService()
	if err != nil {
		panic(fmt.Errorf("failed to create service instance: %v", err))
	}

	httpHandler := handler.NewHTTPHandler(service)
	r.GET("/", httpHandler.Handle)
	api := r.Group("api/v1")

	api.PUT("/score/:leaderboard", httpHandler.HandlePutScore)
	api.GET("/scores/:leaderboard", httpHandler.HandleGetScores)

	err = r.Run(config.GetAddr())
	if err != nil {
		panic(fmt.Errorf("failed to start the server %v", err))
	}

}

func createService() (*services.LeaderboardsService, error) {
	var cfg aws.Config
	if config.IsLocal() {
		fmt.Println("Running in local mode")
		cfg = repository.DefaultLocalAWSClientConfig()
	} else {
		cfg = *aws.NewConfig()
	}
	
	settings := repository.DynamoDBSettings{
		Client: dynamodb.NewFromConfig(cfg),
		Logger: logging.NewSimpleLogger(),
		Table:  config.GetDynamoDBTableName(),
	}

	repo, err := repository.NewDynamoDBRepository(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamodb repository: %v", err)
	}

	configProvider := configprovider.NewDynamoConfigProvider(repo, settings.Logger)

	scoreboard, err := scoreboard.NewRedisScoreboard(config.GetRedisAddr())
	if err != nil {
		return nil, fmt.Errorf("failed to create redis scoreboard: %v", err)
	}
	return services.NewLeaderboardsService(repo, scoreboard, configProvider), nil
}
