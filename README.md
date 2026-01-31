# soundplay

Multiplatform CLI app to play sounds. Good for CLI scripting.

## Features

- 🎵 **Multiple formats**: MP3, OGG (Vorbis), WAV, FLAC
- 🌐 **URL streaming**: Play audio from HTTP/HTTPS URLs
- 🚀 **Single binary**: No runtime dependencies
- 💻 **Cross-platform**: Works on Linux, macOS, Windows
- 🔧 **Simple**: Just `soundplay <file-or-url>`

## Installation

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/michalbcz/soundplay/releases):

```bash
# Linux (amd64)
wget https://github.com/michalbcz/soundplay/releases/latest/download/soundplay-linux-amd64
chmod +x soundplay-linux-amd64
sudo mv soundplay-linux-amd64 /usr/local/bin/soundplay

# macOS (arm64)
curl -LO https://github.com/michalbcz/soundplay/releases/latest/download/soundplay-darwin-arm64
chmod +x soundplay-darwin-arm64
sudo mv soundplay-darwin-arm64 /usr/local/bin/soundplay

# Windows
# Download soundplay-windows-amd64.exe and add to PATH
```

### Build from Source

Requirements:
- Go 1.21 or later
- ALSA development files (Linux only): `sudo apt-get install libasound2-dev`

```bash
git clone https://github.com/michalbcz/soundplay.git
cd soundplay
go build -ldflags "-s -w" -o soundplay
sudo mv soundplay /usr/local/bin/
```

## Usage

### Basic Usage

```bash
# Play a local file
soundplay music.mp3

# Play from home directory
soundplay ~/sounds/notification.wav

# Play from URL
soundplay https://example.com/audio.ogg

# Show help
soundplay -h

# Show version
soundplay -v
```

### Examples

```bash
# Play an MP3
soundplay song.mp3

# Play a WAV file
soundplay ~/sounds/beep.wav

# Play OGG Vorbis
soundplay podcast.ogg

# Play FLAC
soundplay high-quality.flac

# Stream from URL
soundplay https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3
```

## Comparison

| Feature | soundplay | ffplay | mpg123 | aplay |
|---------|-----------|--------|--------|-------|
| Single binary | ✅ | ❌ | ❌ | ✅ |
| Multiple formats | ✅ | ✅ | ❌ | ❌ |
| URL support | ✅ | ✅ | ✅ | ❌ |
| No GUI | ✅ | ❌ | ✅ | ✅ |
| Cross-platform | ✅ | ✅ | ✅ | ❌ |
| Binary size | ~6MB | ~70MB | ~2MB | <1MB |

## Build Requirements

- **Go**: 1.21 or later
- **Linux**: ALSA development files (`libasound2-dev`)
- **macOS**: No additional requirements
- **Windows**: No additional requirements

### Build Commands

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install

# Clean build artifacts
make clean

# Show binary sizes
make sizes
```

## Supported Formats

- **MP3**: MPEG Layer 3 audio
- **OGG**: Ogg Vorbis audio
- **WAV**: Waveform Audio File Format
- **FLAC**: Free Lossless Audio Codec

## Technical Details

- Uses [gopxl/beep/v2](https://github.com/gopxl/beep) for audio decoding and playback
- Downloads entire audio file into memory for URL streaming (decoders require seeking)
- HTTP client timeout: 30 seconds
- Speaker initialization uses format's native sample rate

## License

MIT License - see [LICENSE](LICENSE) file for details.

Copyright (c) 2026 michalbcz

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
