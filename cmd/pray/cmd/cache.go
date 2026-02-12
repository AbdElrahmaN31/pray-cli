package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/anashaat/pray-cli/internal/config"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Cache management",
	Long:  `Manage the pray CLI cache for prayer times data.`,
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all cached data",
	Long:  `Remove all cached prayer times data from the cache directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cacheDir, err := config.GetCacheDir()
		if err != nil {
			return fmt.Errorf("failed to get cache directory: %w", err)
		}

		green := color.New(color.FgGreen).SprintFunc()

		// Check if cache directory exists
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			fmt.Println("Cache directory does not exist. Nothing to clear.")
			return nil
		}

		// Get cache size before clearing
		size, count, err := getCacheStats(cacheDir)
		if err != nil {
			return fmt.Errorf("failed to get cache stats: %w", err)
		}

		// Remove all files in cache directory
		entries, err := os.ReadDir(cacheDir)
		if err != nil {
			return fmt.Errorf("failed to read cache directory: %w", err)
		}

		for _, entry := range entries {
			path := filepath.Join(cacheDir, entry.Name())
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("failed to remove %s: %w", path, err)
			}
		}

		fmt.Printf("%s Cache cleared!\n", green("✓"))
		fmt.Printf("  Removed %d files (%s)\n", count, formatSize(size))
		return nil
	},
}

var cacheShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show cache status and size",
	Long:  `Display information about the current cache including size and file count.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cacheDir, err := config.GetCacheDir()
		if err != nil {
			return fmt.Errorf("failed to get cache directory: %w", err)
		}

		cyan := color.New(color.FgCyan).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		fmt.Println()
		fmt.Printf("📦 %s\n", cyan("Cache Status"))
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		// Check if cache directory exists
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			fmt.Printf("  Status:    %s\n", yellow("Empty (no cache directory)"))
			fmt.Printf("  Location:  %s\n", cacheDir)
			fmt.Println()
			return nil
		}

		// Get cache stats
		size, count, err := getCacheStats(cacheDir)
		if err != nil {
			return fmt.Errorf("failed to get cache stats: %w", err)
		}

		// Check if caching is enabled
		cfg := GetConfig()
		cacheEnabled := "Enabled"
		if !cfg.CacheEnabled {
			cacheEnabled = yellow("Disabled")
		}

		fmt.Printf("  Status:    %s\n", cacheEnabled)
		fmt.Printf("  Location:  %s\n", cacheDir)
		fmt.Printf("  Files:     %d\n", count)
		fmt.Printf("  Size:      %s\n", formatSize(size))
		fmt.Println()

		// List cache files if there are any
		if count > 0 && count <= 10 {
			fmt.Println("  Cached files:")
			entries, _ := os.ReadDir(cacheDir)
			for _, entry := range entries {
				if !entry.IsDir() {
					info, _ := entry.Info()
					if info != nil {
						fmt.Printf("    - %s (%s)\n", entry.Name(), formatSize(info.Size()))
					} else {
						fmt.Printf("    - %s\n", entry.Name())
					}
				}
			}
			fmt.Println()
		}

		return nil
	},
}

var cachePathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show cache directory path",
	Long:  `Display the path to the cache directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cacheDir, err := config.GetCacheDir()
		if err != nil {
			return fmt.Errorf("failed to get cache directory: %w", err)
		}
		fmt.Println(cacheDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
	cacheCmd.AddCommand(cacheClearCmd)
	cacheCmd.AddCommand(cacheShowCmd)
	cacheCmd.AddCommand(cachePathCmd)
}

// getCacheStats returns the total size and file count in the cache directory
func getCacheStats(dir string) (int64, int, error) {
	var totalSize int64
	var fileCount int

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
			fileCount++
		}
		return nil
	})

	return totalSize, fileCount, err
}

// formatSize formats bytes as a human-readable string
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}
