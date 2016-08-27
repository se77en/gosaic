package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

var (
	maxRepeats         int
	mosaicBuildMacroId int
	mosaicBuildType    string
)

func init() {
	addLocalIntFlag(&mosaicBuildMacroId, "macro_id", "", 0, "Id of macro to use to build mosaic", MosaicBuildCmd)
	addLocalIntFlag(&maxRepeats, "max_repeats", "", -1, "Number of times an index image can be repeated in the mosaic, 0 indicates unlimited, -1 is the minimum number", MosaicBuildCmd)
	addLocalFlag(&mosaicBuildType, "type", "", "random", "Mosaic build type, either 'best' or 'random'", MosaicBuildCmd)
	RootCmd.AddCommand(MosaicBuildCmd)
}

var MosaicBuildCmd = &cobra.Command{
	Use:   "mosaic_build NAME",
	Short: "Build mosaic",
	Long:  "Build mosaic",
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			Env.Fatalln("Mosaic name is required")
		}

		if args[0] == "" {
			Env.Fatalln("Mosaic name is required")
		}

		if mosaicBuildMacroId == 0 {
			Env.Fatalln("Macro id is required")
		}

		if mosaicBuildType != "best" && mosaicBuildType != "random" {
			Env.Fatalln("type must be either 'best' or 'random'")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.MosaicBuild(Env, args[0], mosaicBuildType, int64(mosaicBuildMacroId), maxRepeats)
	},
}
