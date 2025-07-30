package tools

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"stream/ent"
)

func TranscodeToHLS(inputPath, outputDir string, client *ent.Client) {
	qualities := []struct {
		Name    string
		Width   int
		Height  int
		Bitrate int
	}{
		{"1080p", 1920, 1080, 5000},
		{"720p", 1280, 720, 3000},
		{"480p", 854, 480, 1000},
	}

	var variants []string

	for _, q := range qualities {
		outputPath := filepath.Join(outputDir, q.Name)
		err := os.MkdirAll(outputPath, 0755)
		if err != nil {
			log.Printf("Erro ao criar pasta para %s: %v", q.Name, err)
			continue
		}

		cmd := exec.Command("ffmpeg",
			"-i", inputPath,
			"-vf", fmt.Sprintf("scale=w=%d:h=%d", q.Width, q.Height),
			"-c:a", "aac",
			"-ar", "48000",
			"-c:v", "h264",
			"-profile:v", "main",
			"-crf", "20",
			"-sc_threshold", "0",
			"-g", "48",
			"-keyint_min", "48",
			"-hls_time", "4",
			"-hls_playlist_type", "vod",
			"-b:v", fmt.Sprintf("%dk", q.Bitrate),
			"-maxrate", fmt.Sprintf("%dk", int(float64(q.Bitrate)*1.07)),
			"-bufsize", fmt.Sprintf("%dk", q.Bitrate*2),
			"-b:a", "128k",
			"-hls_segment_filename", filepath.Join(outputPath, "segment_%03d.ts"),
			filepath.Join(outputPath, "index.m3u8"),
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Erro ao transcodificar para %s: %v\nSaída: %s", q.Name, err, string(output))
			continue
		}

		variants = append(variants, fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n%s/index.m3u8",
			q.Bitrate*1000, q.Width, q.Height, q.Name))
	}

	// Cria o master.m3u8
	masterPath := filepath.Join(outputDir, "master.m3u8")
	masterFile, err := os.Create(masterPath)
	if err != nil {
		log.Printf("Erro ao criar master.m3u8: %v", err)
		return
	}
	defer masterFile.Close()

	_, _ = masterFile.WriteString("#EXTM3U\n")

	for _, variant := range variants {
		_, _ = masterFile.WriteString(variant + "\n")
	}

	log.Println("Processamento HLS concluído com sucesso para:", inputPath)
}
