# Muse

Your music, in your browser. Muse is a minimal, self-hostable music player that
serves any folder containing mp3 files.

Written in Go with zero dependencies.

![A screenshot demonstrating the app running in a browser.](.github/screenshot.png)

## Installation & Tool Usage

```bash
# Install the latest tool.
go install github.com/prophittcorey/muse/cmd/muse@latest

# Serve your music.
muse --dir "~/Music" --host "0.0.0.0" --port "3000"
```

## License

The source code for this repository is licensed under the MIT license, which you can
find in the [LICENSE](LICENSE.md) file.
