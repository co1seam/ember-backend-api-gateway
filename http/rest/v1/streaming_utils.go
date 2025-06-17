package v1

import (
	"os/exec"
	"path/filepath"
)

func convertToHLSVP9(inputPath, outputPath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c:v", "libvpx-vp9", "-crf", "30", "-b:v", "0", "-row-mt", "1", "-threads", "8",
		"-c:a", "libopus", "-b:a", "128k",
		"-f", "hls",
		"-hls_time", "10",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(outputPath, "segment_%03d.m4s"),
		filepath.Join(outputPath, "playlist.m3u8"))
	return cmd.Run()
}
