package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

const version = "1.0.0"

func main() {
	// Define flags
	helpFlag := flag.Bool("h", false, "Show help message")
	helpLongFlag := flag.Bool("help", false, "Show help message")
	versionFlag := flag.Bool("v", false, "Show version")
	versionLongFlag := flag.Bool("version", false, "Show version")

	flag.Parse()

	// Handle flags
	if *helpFlag || *helpLongFlag {
		showHelp()
		return
	}

	if *versionFlag || *versionLongFlag {
		fmt.Printf("soundplay version %s\n", version)
		return
	}

	// Check for audio source argument
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No audio source specified")
		fmt.Fprintln(os.Stderr, "Usage: soundplay <file-or-url>")
		fmt.Fprintln(os.Stderr, "Try 'soundplay -h' for more information.")
		os.Exit(1)
	}

	source := args[0]

	// Play the audio
	if err := play(source); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("soundplay - Multiplatform CLI app to play sounds")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  soundplay <file-or-url>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help     Show this help message")
	fmt.Println("  -v, --version  Show version information")
	fmt.Println()
	fmt.Println("Supported formats:")
	fmt.Println("  MP3, OGG (Vorbis), WAV, FLAC")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  soundplay music.mp3")
	fmt.Println("  soundplay ~/sounds/beep.wav")
	fmt.Println("  soundplay https://example.com/sound.ogg")
}

func play(source string) error {
	// Open the audio source
	streamer, format, err := openAudio(source)
	if err != nil {
		return fmt.Errorf("failed to open audio: %w", err)
	}
	defer streamer.Close()

	// Initialize speaker
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Create done channel
	done := make(chan bool)

	// Play the audio
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for playback to finish
	<-done

	return nil
}

func openAudio(source string) (beep.StreamSeekCloser, beep.Format, error) {
	// Check if source is a URL
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return openURL(source)
	}

	// Otherwise, treat as file path
	return openFile(source)
}

func openFile(path string) (beep.StreamSeekCloser, beep.Format, error) {
	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, beep.Format{}, fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(home, path[2:])
	}

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("failed to open file: %w", err)
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(path))

	// Decode based on extension
	streamer, format, err := decode(file, ext)
	if err != nil {
		file.Close()
		return nil, beep.Format{}, err
	}

	return streamer, format, nil
}

func openURL(url string) (beep.StreamSeekCloser, beep.Format, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Fetch the URL
	resp, err := client.Get(url)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, beep.Format{}, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	// Read entire content into memory (needed for seeking)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("failed to read response: %w", err)
	}

	// Get extension from URL
	ext := strings.ToLower(filepath.Ext(url))
	// Remove query parameters from extension
	if idx := strings.Index(ext, "?"); idx != -1 {
		ext = ext[:idx]
	}

	// Create a ReadSeekCloser from bytes
	reader := newBytesReadCloser(data)

	// Decode based on extension
	streamer, format, err := decode(reader, ext)
	if err != nil {
		reader.Close()
		return nil, beep.Format{}, err
	}

	return streamer, format, nil
}

func decode(reader io.ReadSeeker, ext string) (beep.StreamSeekCloser, beep.Format, error) {
	// Convert reader to appropriate interface based on decoder requirements
	// MP3 and Vorbis need io.ReadCloser, WAV and FLAC need io.Reader
	switch ext {
	case ".mp3":
		// mp3.Decode expects io.ReadCloser
		rc, ok := reader.(io.ReadCloser)
		if !ok {
			// Wrap with nopCloser if not already a ReadCloser
			rc = io.NopCloser(reader)
		}
		return mp3.Decode(rc)
	case ".ogg":
		// vorbis.Decode expects io.ReadCloser
		rc, ok := reader.(io.ReadCloser)
		if !ok {
			// Wrap with nopCloser if not already a ReadCloser
			rc = io.NopCloser(reader)
		}
		return vorbis.Decode(rc)
	case ".wav":
		// wav.Decode expects io.Reader
		return wav.Decode(reader)
	case ".flac":
		// flac.Decode expects io.Reader
		return flac.Decode(reader)
	default:
		return nil, beep.Format{}, fmt.Errorf("unsupported audio format: %s (supported: .mp3, .ogg, .wav, .flac)", ext)
	}
}

// bytesReadCloser wraps a bytes.Reader to implement io.ReadSeekCloser
type bytesReadCloser struct {
	*bytes.Reader
}

func newBytesReadCloser(data []byte) *bytesReadCloser {
	return &bytesReadCloser{
		Reader: bytes.NewReader(data),
	}
}

func (b *bytesReadCloser) Close() error {
	return nil
}
