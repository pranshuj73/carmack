package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"strings"

	"github.com/pranshuj73/carmack/utils"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	directory string
	editor    string
)

var rootCmd = &cobra.Command{
	Use: "carmack",
	Short: "Carmack is a CLI tool to help you manage your daily plans.",
	Long: `Carmack is a CLI tool to help you manage your daily plans. You can create, edit, and delete plans using this tool. You can also sync your plans with a remote git repository.`,
	Args:  cobra.MaximumNArgs(1),
	Run: handlePlanFile,
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/.carmack.yaml)")
	rootCmd.PersistentFlags().StringVar(&directory, "directory", "", "Set the directory for storing plan files.")
	rootCmd.PersistentFlags().StringVar(&editor, "editor", "nvim", "Set the editor for opening plan files.")

	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))
	viper.BindPFlag("editor", rootCmd.PersistentFlags().Lookup("editor"))

	viper.SetDefault("editor", "nvim")
}

func initConfig() {
	// Determine config file location
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("Error finding home directory:", err)
			os.Exit(1)
		}
		viper.AddConfigPath(filepath.Join(home, ".config"))
		viper.SetConfigName(".carmack")
		viper.SetConfigType("yaml")
		cfgFile = filepath.Join(home, ".config", ".carmack.yaml")
	}

	// Read in existing config file or create a new one
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Creating a new config file...")
			saveConfig()
		} else {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
	}
}

func saveConfig() {
	if err := viper.WriteConfigAs(cfgFile); err != nil {
		fmt.Printf("Error writing config to %s: %v\n", cfgFile, err)
		os.Exit(1)
	}
}

func handlePlanFile(cmd *cobra.Command, args []string) {
	// Update the config if the directory flag is provided
	dir, _ := cmd.Flags().GetString("directory")
	if dir != "" {
		viper.Set("directory", dir)
		saveConfig()
		fmt.Printf("Directory set to: %s\n", dir)
		return
	}

	// Retrieve directory from config
	dir = viper.GetString("directory")
	if dir == "" {
    fmt.Println("Error: Please enter the directory you'd like to use with carmack:")
    fmt.Scanln(&dir)
    viper.Set("directory", dir)
    saveConfig()
    fmt.Printf("Directory set to: %s\n", dir)
		os.Exit(1)
	}

	editor := viper.GetString("editor")

	// Determine the date based on the input argument
	var date string
	now := time.Now()
	if len(args) > 0 {
		switch strings.ToLower(args[0]) {
		case "today":
			date = now.Format("02012006")
		case "yesterday":
			date = now.AddDate(0, 0, -1).Format("02012006")
		default:
      // check if args[0] is a valid date
      _, err := time.Parse("02012006", args[0])
      if err != nil {
        fmt.Println("Error: Please enter a valid date in the format DDMMYYYY.")
        os.Exit(1)
      }
			date = args[0] // Assume the argument is already a valid date
		}
	} else {
		date = now.Format("02012006") // Default to today's date
	}

	// Construct the file path
	expandedDir, err := homedir.Expand(dir)
	if err != nil {
		fmt.Println("Error expanding directory path:", err)
		os.Exit(1)
	}

	filename := filepath.Join(expandedDir, fmt.Sprintf("%s.md", date))
	err = utils.EnsureDirectoryExists(expandedDir)
	if err != nil {
		fmt.Println("Error ensuring directory exists:", err)
		os.Exit(1)
	}

	err = utils.OpenFileWithEditor(editor, filename, date)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
}
