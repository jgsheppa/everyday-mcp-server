package tools

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type CompressArgs struct {
	SourcePath      string `json:"source_path,omitempty"`
	Base64Content   string `json:"base64_content,omitempty"`
	OutputFilename  string `json:"output_filename"`
	DestinationPath string `json:"destination_path,omitempty"`
	CompressionType string `json:"compression_type,omitempty"`
	Quality         int    `json:"quality,omitempty"`
}

func compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gw := gzip.NewWriter(dstFile)
	defer gw.Close()

	_, err = io.Copy(gw, srcFile)
	return err
}

func detectImageFormat(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	_, format, err := image.DecodeConfig(reader)
	if err != nil {
		return "", err
	}
	return format, nil
}

func compressJPEGFromBytes(data []byte, outputPath string, quality int) error {
	// Decode image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Set default quality if not provided
	if quality <= 0 || quality > 100 {
		quality = 85 // Default JPEG quality
	}

	// Encode with specified quality
	options := &jpeg.Options{Quality: quality}
	return jpeg.Encode(outFile, img, options)
}

func compressJPEGFromFile(srcPath, dstPath string, quality int) error {
	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Decode image
	img, _, err := image.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Create output file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer dstFile.Close()

	// Set default quality if not provided
	if quality <= 0 || quality > 100 {
		quality = 85 // Default JPEG quality
	}

	// Encode with specified quality
	options := &jpeg.Options{Quality: quality}
	return jpeg.Encode(dstFile, img, options)
}

func compressBase64Content(base64Content, outputPath, compressionType string, quality int) error {
	// Decode base64 content
	data, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return fmt.Errorf("failed to decode base64 content: %w", err)
	}

	// Auto-detect compression type if not specified
	if compressionType == "" {
		if format, err := detectImageFormat(data); err == nil {
			if format == "jpeg" {
				compressionType = "jpeg"
			} else {
				compressionType = "gzip" // fallback to gzip for other formats
			}
		} else {
			compressionType = "gzip" // fallback to gzip if not an image
		}
	}

	// Use appropriate compression method
	switch compressionType {
	case "jpeg":
		return compressJPEGFromBytes(data, outputPath, quality)
	case "gzip":
		// Create output file
		dstFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer dstFile.Close()

		// Create gzip writer
		gw := gzip.NewWriter(dstFile)
		defer gw.Close()

		// Compress the data
		_, err = io.Copy(gw, bytes.NewReader(data))
		return err
	default:
		return fmt.Errorf("unsupported compression type: %s", compressionType)
	}
}

func CompressFile(
	_ context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[CompressArgs],
) (*mcp.CallToolResult, error) {
	sourcePath := strings.TrimSpace(params.Arguments.SourcePath)
	base64Content := strings.TrimSpace(params.Arguments.Base64Content)
	outputFilename := strings.TrimSpace(params.Arguments.OutputFilename)
	destinationPath := strings.TrimSpace(params.Arguments.DestinationPath)
	compressionType := strings.ToLower(strings.TrimSpace(params.Arguments.CompressionType))
	quality := params.Arguments.Quality

	// Validate input: must have either source_path or base64_content, and output_filename
	if sourcePath == "" && base64Content == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Either source_path or base64_content must be provided",
				},
			},
			IsError: true,
		}, nil
	}

	if sourcePath != "" && base64Content != "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Cannot provide both source_path and base64_content. Choose one input method",
				},
			},
			IsError: true,
		}, nil
	}

	if outputFilename == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "output_filename is required",
				},
			},
			IsError: true,
		}, nil
	}

	// Determine final destination path
	var finalDestinationPath string
	if destinationPath != "" {
		finalDestinationPath = filepath.Join(destinationPath, outputFilename+".gz")
	} else {
		// Use current working directory if no destination specified
		cwd, err := os.Getwd()
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("failed to get current working directory: %v", err),
					},
				},
				IsError: true,
			}, nil
		}
		finalDestinationPath = filepath.Join(cwd, outputFilename+".gz")
	}

	// Check if destination already exists
	if _, err := os.Stat(finalDestinationPath); err == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("destination file already exists: %s", finalDestinationPath),
				},
			},
			IsError: true,
		}, nil
	}

	var originalSize int64
	var compressionErr error

	// Handle file path input
	if sourcePath != "" {
		// Check if source file exists
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("source file does not exist: %s", sourcePath),
					},
				},
				IsError: true,
			}, nil
		}

		// Get source file info for size reporting
		sourceInfo, err := os.Stat(sourcePath)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("failed to get source file info: %v", err),
					},
				},
				IsError: true,
			}, nil
		}
		originalSize = sourceInfo.Size()

		// Perform compression
		compressionErr = compressFile(sourcePath, finalDestinationPath)
	} else {
		// Handle base64 content input
		// Estimate original size from base64 string (approximate)
		decodedData, err := base64.StdEncoding.DecodeString(base64Content)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("failed to decode base64 content: %v", err),
					},
				},
				IsError: true,
			}, nil
		}
		originalSize = int64(len(decodedData))

		// Perform compression
		compressionErr = compressBase64Content(base64Content, finalDestinationPath, compressionType, quality)
	}

	if compressionErr != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to compress content: %v", compressionErr),
				},
			},
			IsError: true,
		}, nil
	}

	// Get compressed file info for size reporting
	compressedInfo, err := os.Stat(finalDestinationPath)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("compression completed but failed to get compressed file info: %v", err),
				},
			},
		}, nil
	}
	compressedSize := compressedInfo.Size()

	compressionRatio := float64(compressedSize) / float64(originalSize) * 100

	sourceDescription := sourcePath
	if sourcePath == "" {
		sourceDescription = "uploaded content"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Successfully compressed content:\nSource: %s (%d bytes)\nDestination: %s (%d bytes)\nCompression ratio: %.1f%%", 
					sourceDescription, originalSize, finalDestinationPath, compressedSize, compressionRatio),
			},
		},
	}, nil
}

func NewCompressTool() *Config[CompressArgs] {
	inputSchema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"source_path": {
				Type:        "string",
				Description: "The path to the file to compress (alternative to base64_content)",
			},
			"base64_content": {
				Type:        "string",
				Description: "Base64-encoded content to compress (alternative to source_path). Use this for uploaded files",
			},
			"output_filename": {
				Type:        "string",
				Description: "The filename for the compressed output (without .gz extension)",
			},
			"destination_path": {
				Type:        "string",
				Description: "Optional directory path where the compressed file will be saved. If not provided, uses current working directory",
			},
		},
		Required: []string{"output_filename"},
	}

	definition := mcp.Tool{
		Name:        "compress_file",
		Description: "Compresses a file using gzip compression. Supports both file paths and base64-encoded content (for uploaded files)",
		Title:       "File Compression",
		InputSchema: inputSchema,
	}

	err := ValidateToolName(definition.Name)
	if err != nil {
		return &Config[CompressArgs]{
			Error: err,
		}
	}

	return &Config[CompressArgs]{
		Definition: &definition,
		Call:       CompressFile,
	}
}

func AddCompressTool(server *mcp.Server) {
	tool := NewCompressTool()
	mcp.AddTool(server, tool.Definition, tool.Call)
}
