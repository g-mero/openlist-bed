package main

import (
	"fmt"
	"os"

	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
)

func getUploadCmd() *cobra.Command {
	var host string
	var apiKey string
	var keepName = false
	var compress = false

	var uploadCmd = &cobra.Command{
		Use:   "upload [file_path] ...",
		Short: "Upload images to the Openlist image bed",
		Long: `Upload one or more image files to the OpenList image bed service.
Example: openlist-bed upload image1.jpg image2.png --host http://localhost:8080 --key your-api-key`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate required flags
			if host == "" {
				return fmt.Errorf("--host is required")
			}
			if apiKey == "" {
				return fmt.Errorf("--key is required")
			}

			// Upload each file
			for _, filePath := range args {
				imageURL, err := uploadFile(filePath, host, apiKey, keepName, compress)
				if err != nil {
					return fmt.Errorf("failed to upload %s: %w", filePath, err)
				}
				// Output image URL for Typora (one URL per line)
				fmt.Println(imageURL)
			}

			return nil
		},
	}

	// Bind flags to local variables
	uploadCmd.Flags().StringVar(&host, "host", "", "OpenList host URL (required)")
	uploadCmd.Flags().StringVar(&apiKey, "key", "", "OpenList API key (required)")
	uploadCmd.Flags().BoolVar(&keepName, "keep", false, "Keep original file names when uploading")
	uploadCmd.Flags().BoolVar(&compress, "compress", false, "Compress images before uploading (not implemented)")

	// Mark flags as required
	_ = uploadCmd.MarkFlagRequired("host")
	_ = uploadCmd.MarkFlagRequired("key")

	return uploadCmd
}

func uploadFile(filePath, host, apiKey string, keepName, compress bool) (string, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("cannot access file: %w", err)
	}

	// Create req client
	client := req.C().SetBaseURL(host).SetCommonHeader("API-Key", apiKey)

	// Prepare response structure
	var result struct {
		URL string `json:"url"`
	}

	// Send upload request
	resp, err := client.R().
		SetFile("image", filePath).
		SetQueryParam("keep_name", fmt.Sprintf("%v", keepName)).
		SetQueryParam("compress", fmt.Sprintf("%v", compress)).
		SetSuccessResult(&result).
		Post("/api/upload")

	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	// Check response
	if !resp.IsSuccessState() {
		return "", fmt.Errorf("server returned status %d: %s", resp.GetStatusCode(), resp.String())
	}

	// Return the image URL
	return result.URL, nil
}
