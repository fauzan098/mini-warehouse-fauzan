package cmd

import (
	"fmt"
	_ "os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	// userLicense string

	rootCmd = &cobra.Command{
		Use:   "warehouse-service",
		Short: "Warehouse service CLI",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Run(startCmd, args)
			// runServer()
		},
	}
)

// Execute executes the root command.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// func runServer() {
//     port := viper.GetString("PORT")
//     if port == "" {
//         port = "8080" // default port
//     }

//     fmt.Printf("Starting User Service on port %s...\n", port)
    
//     // Masukkan kode untuk inisialisasi database, router, atau server kamu di sini
// }

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")

	// rootCmd.AddCommand(addCmd)
	// rootCmd.AddCommand(initCmd)

	rootCmd.Flags().BoolP("toggle", "t", false, "help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		viper.SetConfigFile(`.env`)
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)

		// // Search config in home directory with name ".cobra" (without extension).
		// viper.AddConfigPath(home)
		// viper.SetConfigType("yaml")
		// viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
