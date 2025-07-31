package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"stream/ent"
	"stream/internal/config"
	"sync"
)

// TranscodeAndProcessVideo é o handler do job que executa todo o pipeline de processamento.
func TranscodeAndProcessVideo(videoID, inputPath string, client *ent.Client, appConfig *config.Config) {
	ctx := context.Background()

	// 1. Atualiza o status para PROCESSING
	_, err := client.File.UpdateOneID(videoID).SetStatus("PROCESSING").Save(ctx)
	if err != nil {
		log.Printf("Error updating status to PROCESSING for video %s: %v", videoID, err)
		return
	}

	// Defer a function to handle panics and update status to FAILED
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in TranscodeAndProcessVideo: %v", r)
			updateStatusToFailed(ctx, client, videoID, fmt.Sprintf("panic: %v", r))
		}
	}()

	outputDir := filepath.Dir(inputPath)

	// 2. Extrair Metadados e Gerar Thumbnail
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		extractAndSaveMetadata(ctx, client, videoID, inputPath)
	}()

	go func() {
		defer wg.Done()
		generateThumbnail(inputPath, outputDir)
	}()

	wg.Wait()

	// 3. Transcodificar para HLS em paralelo
	if err := transcodeToHLSParallel(inputPath, outputDir, appConfig);
	 err != nil {
		log.Printf("Error during HLS transcoding for video %s: %v", videoID, err)
		updateStatusToFailed(ctx, client, videoID, err.Error())
		return
	}

	// 4. Atualiza o status para COMPLETED
	_, err = client.File.UpdateOneID(videoID).SetStatus("COMPLETED").SetVisibility(true).Save(ctx)
	if err != nil {
		log.Printf("Error updating status to COMPLETED for video %s: %v", videoID, err)
		// Mesmo que o status não seja atualizado, o processo foi concluído.
	}

	log.Printf("Successfully processed video %s", videoID)
}

func transcodeToHLSParallel(inputPath, outputDir string, appConfig *config.Config) error {
	var wg sync.WaitGroup
	var variants []string
	mu := &sync.Mutex{}
	errs := make(chan error, len(appConfig.Transcoding.Qualities))

	for _, q := range appConfig.Transcoding.Qualities {
		wg.Add(1)
		go func(quality config.Quality) {
			defer wg.Done()

			outputPath := filepath.Join(outputDir, quality.Name)
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				errs <- fmt.Errorf("failed to create directory for %s: %w", quality.Name, err)
				return
			}

			args := buildFFmpegArgs(inputPath, outputPath, quality, appConfig.Transcoding.FFmpegParams)
			cmd := exec.Command("ffmpeg", args...)

			progressLogger, err := NewProgressLogger(inputPath)
			if err != nil {
				errs <- fmt.Errorf("failed to create progress logger for %s: %w", quality.Name, err)
				return
			}

			stderr, err := cmd.StderrPipe()
			if err != nil {
				errs <- fmt.Errorf("failed to get stderr pipe for %s: %w", quality.Name, err)
				return	
			}

			if err := cmd.Start(); err != nil {
				errs <- fmt.Errorf("failed to start transcoding for %s: %w", quality.Name, err)
				return
			}

			go progressLogger.LogProgress(stderr)

			if err := cmd.Wait(); err != nil {
				errs <- fmt.Errorf("transcoding failed for %s: %w", quality.Name, err)
				return
			}

			mu.Lock()
			variants = append(variants, fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n%s/index.m3u8",
				quality.Bitrate*1000, quality.Width, quality.Height, quality.Name))
			mu.Unlock()
		}(q)
	}

	wg.Wait()
	close(errs)

	// Check for errors from goroutines
	for err := range errs {
		if err != nil {
			return err // Return the first error encountered
		}
	}

	// Create master playlist
	return createMasterPlaylist(outputDir, variants)
}

func buildFFmpegArgs(inputPath, outputPath string, quality config.Quality, baseParams []string) []string {
	args := []string{"-i", inputPath}
	args = append(args, "-vf", fmt.Sprintf("scale=w=%d:h=%d", quality.Width, quality.Height))
	args = append(args, baseParams...)
	args = append(args, "-b:v", fmt.Sprintf("%dk", quality.Bitrate))
	args = append(args, "-maxrate", fmt.Sprintf("%dk", int(float64(quality.Bitrate)*1.07)))
	args = append(args, "-bufsize", fmt.Sprintf("%dk", quality.Bitrate*2))
	args = append(args, "-hls_segment_filename", filepath.Join(outputPath, "segment_%03d.ts"))
	args = append(args, filepath.Join(outputPath, "index.m3u8"))
	return args
}

func createMasterPlaylist(outputDir string, variants []string) error {
	masterPath := filepath.Join(outputDir, "master.m3u8")
	masterFile, err := os.Create(masterPath)
	if err != nil {
		return fmt.Errorf("failed to create master.m3u8: %w", err)
	}
	defer masterFile.Close()

	_, _ = masterFile.WriteString("#EXTM3U\n")
	for _, variant := range variants {
		_, _ = masterFile.WriteString(variant + "\n")
	}
	return nil
}

func extractAndSaveMetadata(ctx context.Context, client *ent.Client, videoID, inputPath string) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", inputPath)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error extracting metadata for video %s: %v", videoID, err)
		return
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(output, &metadata); err != nil {
		log.Printf("Error unmarshalling metadata for video %s: %v", videoID, err)
		return
	}

	_, err = client.File.UpdateOneID(videoID).SetMetadata(metadata).Save(ctx)
	if err != nil {
		log.Printf("Error saving metadata for video %s: %v", videoID, err)
	}
}

func generateThumbnail(inputPath, outputDir string) {
	thumbnailPath := filepath.Join(outputDir, "thumbnail.jpg")
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-ss", "00:00:01.000", "-vframes", "1", thumbnailPath)
	if err := cmd.Run(); err != nil {
		log.Printf("Error generating thumbnail for %s: %v", inputPath, err)
	}
}

func updateStatusToFailed(ctx context.Context, client *ent.Client, videoID, details string) {
	_, err := client.File.UpdateOneID(videoID).SetStatus("FAILED").SetStatusDetails(details).Save(ctx)
	if err != nil {
		log.Printf("CRITICAL: Failed to update status to FAILED for video %s: %v", videoID, err)
	}
}