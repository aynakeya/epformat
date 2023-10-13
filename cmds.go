package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strings"
)

func RenameEpInfo(extractor EpisodeExtractor, value, title, format string) (string, error) {
	value = strings.TrimSpace(value)
	episode := extractor.Extract(value)
	if title != "" {
		episode.Title = title
	}
	return episode.FormatInfo(format)
}

func createFormatCmd() *cobra.Command {
	var verbose bool
	formatCmd := &cobra.Command{
		Use:   "format <name>...",
		Short: "format anime file with proper episode format",
		Long:  "format anime file with proper episode format",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			title := cmd.Flag("title").Value.String()
			format := cmd.Flag("format").Value.String()
			for _, arg := range args {
				renamed, err := RenameEpInfo(MainExtractor, arg, title, format)
				if verbose {
					if err != nil {
						fmt.Printf("- \"%s\" => %s\n", arg, err)
					} else {
						fmt.Printf("- \"%s\" => \"%s\"\n", arg, renamed)
					}
				} else {
					if err != nil {
						fmt.Printf("%s\n", err)
					} else {
						fmt.Printf("%s\n", renamed)
					}
				}

			}
		},
	}
	formatCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose print")
	return formatCmd
}

type fileInfo struct {
	path string
	info os.FileInfo
}

func getAllFiles(dir string) []fileInfo {
	var files []fileInfo = make([]fileInfo, 0)
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			files = append(files, fileInfo{path: path, info: info})
		}
		return nil
	})
	return files
}

func removeSpecialChars(str string) string {
	str = strings.ReplaceAll(str, ":", "")
	str = strings.ReplaceAll(str, "/", "")
	str = strings.ReplaceAll(str, "\\", "")
	str = strings.ReplaceAll(str, "*", "")
	str = strings.ReplaceAll(str, "?", "")
	str = strings.ReplaceAll(str, "\"", "")
	str = strings.ReplaceAll(str, "<", "")
	str = strings.ReplaceAll(str, ">", "")
	str = strings.ReplaceAll(str, "|", "")
	return str
}

func createRenameCmd() *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "rename <file/directory>...",
		Short: "rename anime file with proper episode format",
		Long:  "rename anime file with proper episode format",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			title := cmd.Flag("title").Value.String()
			format := cmd.Flag("format").Value.String()
			files := make([]fileInfo, 0)
			for _, arg := range args {
				files = append(files, getAllFiles(arg)...)
			}
			for _, file := range files {
				fileName := file.info.Name()
				renamed, err := RenameEpInfo(MainExtractor, fileName, title, format)
				if err != nil {
					fmt.Printf("- \"%s\" => %s\n", fileName, err)
					continue
				}
				if !yes {
					fmt.Printf("- \"%s\" => \"%s\" rename? (y/n) ", fileName, renamed)
					var answer string
					_, _ = fmt.Scanln(&answer)
					yes = answer == "y"
				}
				if yes {
					err = os.Rename(file.path, filepath.Join(filepath.Dir(file.path), removeSpecialChars(renamed)))
					if err != nil {
						fmt.Printf("- \"%s\" => \"%s\" (%s)\n", fileName, renamed, err)
					} else {
						fmt.Printf("- \"%s\" => \"%s\" (ok)\n", fileName, renamed)
					}
				}
			}
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "rename file without confirmation")
	return cmd
}

func createRootCmd() *cobra.Command {
	var format string
	var title string
	var configFile string

	rootCmd := &cobra.Command{
		Use:   "epformat",
		Short: "format anime file with proper episode format",
		Long:  "format anime file with proper episode format",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if configFile == "" {
				return nil
			}
			config, err := ini.Load(configFile)
			if err != nil {
				return nil
			}
			tmp := config.Section("").Key("title").String()
			if tmp != "" {
				title = tmp
			}
			tmp = config.Section("").Key("format").String()
			if tmp != "" {
				format = tmp
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", DefaultFormat, "format string")
	rootCmd.PersistentFlags().StringVarP(&title, "title", "t", "", "episode title")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "use config file")
	rootCmd.AddCommand(createFormatCmd(), createRenameCmd())
	return rootCmd
}
