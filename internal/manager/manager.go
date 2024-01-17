package manager

import (
	"github.com/quipham98/cdn-testing/internal/config"
	"github.com/quipham98/cdn-testing/internal/processor"
	"github.com/quipham98/cdn-testing/internal/report"
	"github.com/quipham98/cdn-testing/internal/service"
	"github.com/spf13/cobra"
	"github.com/trustwallet/assets-go-libs/file"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

var configPath, root string

func InitCommands() {
	rootCmd.Flags().StringVar(&configPath, "config", ".github/assets.config.yaml",
		"config file (default is $HOME/.github/assets.config.yaml)")
	rootCmd.Flags().StringVar(&root, "root", ".", "root path to files")

	rootCmd.AddCommand(checkCmd)
}

var (
	rootCmd = &cobra.Command{
		Use:   "assets",
		Short: "",
		Long:  "",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: " Execute validation checks",
		Run: func(cmd *cobra.Command, args []string) {
			assetsService := InitAssetsService()
			assetsService.RunJob(assetsService.Check)
		},
	}
)

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func InitAssetsService() *service.Service {
	setup()

	paths, err := file.ReadLocalFileStructure(root, config.Default.ValidatorsSettings.RootFolder.SkipFiles)
	if err != nil {
		log.WithError(err).Fatal("Failed to load file structure.")
	}

	paths = filter(paths, func(path string) bool {
		for _, dir := range config.Default.ValidatorsSettings.RootFolder.SkipDirs {
			if strings.Contains(path, dir) {
				return false
			}
		}
		return true
	})

	fileService := file.NewService(paths...)
	validatorsService := processor.NewService(fileService)
	reportService := report.NewService()

	return service.NewService(fileService, validatorsService, reportService, paths)
}

func setup() {
	if err := config.SetConfig(configPath); err != nil {
		log.WithError(err).Fatal("Failed to set config.")
	}

	logLevel, err := log.ParseLevel(config.Default.App.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("Failed to parse log level.")
	}

	log.SetLevel(logLevel)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
