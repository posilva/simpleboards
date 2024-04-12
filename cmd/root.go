package cmd

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/posilva/simpleboards/cmd/leaderboards/config"
	"github.com/posilva/simpleboards/internal/adapters/input/handler"
	"github.com/posilva/simpleboards/internal/adapters/output/configprovider"
	"github.com/posilva/simpleboards/internal/adapters/output/repository"
	"github.com/posilva/simpleboards/internal/adapters/output/scoreboard"
	"github.com/posilva/simpleboards/internal/core/services"
	"github.com/posilva/simpleboards/internal/testutil"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simpleboards",
	Short: "Service for leaderboards",
	Long: ` Service to handle leaderboards manager
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		r.Use(gin.Logger())
		r.Use(gin.Recovery())

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
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.simpleboards.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createService() (*services.LeaderboardsService, error) {
	// TODO: For local testing
	settings := testutil.NewDefaultDynamoDBSettings()
	repo, err := repository.NewDynamoDBRepository(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create repositort %v", err)
	}

	config := configprovider.NewSimpleConfigProvider(repo, settings.Logger)
	scoreboard, err := scoreboard.NewRedisScoreboard("localhost:6379")
	if err != nil {
		return nil, fmt.Errorf("failed to create redis scoreboard: %v", err)
	}
	return services.NewLeaderboardsService(repo, scoreboard, config), nil
}
