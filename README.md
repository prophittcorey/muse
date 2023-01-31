# Muse

Your music, in your browser. Muse is a minimal, self-hostable music player
written in Go with zero dependencies.

![A screenshot demonstrating the app running in a browser.](.github/screenshot.png)

## Installation & Tool Usage

```bash
# Install the latest tool.
go install github.com/prophittcorey/muse/cmd/muse@latest

# Serve your music.
muse --dir "~/Music" --host "0.0.0.0" --port "3000"

# Open http://0.0.0.0:3000 with any browser on your network.
```

## How it Works

Muse locates all mp3 files within the specified directory and all subdirectories.

Each mp3 file's internal ID3 tags are parsed and used for each track's artist,
title and album artwork.

## License

The source code for this repository is licensed under the MIT license, which you can
find in the [LICENSE](LICENSE.md) file.
