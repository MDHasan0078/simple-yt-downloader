# Simple YT Downloader

[![Latest release](https://img.shields.io/github/v/release/MDHasan0078/simple-yt-downloader?label=latest&style=flat-square)](https://github.com/MDHasan0078/simple-yt-downloader/releases)
[![License](https://img.shields.io/github/license/MDHasan0078/simple-yt-downloader?style=flat-square)](LICENSE)
[![Last commit](https://img.shields.io/github/last-commit/MDHasan0078/simple-yt-downloader?style=flat-square)](https://github.com/MDHasan0078/simple-yt-downloader)
[![Platform](https://img.shields.io/badge/platform-Linux%20%E2%80%A2%20macOS%20%E2%80%A2%20Windows-blue?style=flat-square)](https://github.com/MDHasan0078/simple-yt-downloader/releases)

A fast, lightweight, open-source YouTube downloader. Every quality from
**144p to 10K**, custom resolution support, selectable **30/60/90 fps**,
playlist support, and pause/resume — all in a single persistent window with
no popup-dialog chains.

> The Linux edition is built with Python + GTK3 and runs on Debian/Ubuntu-based
> distros (tested on Linux Mint). macOS and Windows builds are Flutter apps
> that share the same download engine.

## Screenshots

Click any thumbnail for the full-size image.

| | |
|:---:|:---:|
| [![Ongoing downloads (dark)](screenshots/main-ongoing-dark.png)](screenshots/main-ongoing-dark.png) | [![Completed downloads (dark)](screenshots/main-completed-dark.png)](screenshots/main-completed-dark.png) |
| Ongoing tab — dark theme | Completed tab — dark theme |
| [![Playlist with per-video progress](screenshots/playlist-dark.png)](screenshots/playlist-dark.png) | [![Settings screen](screenshots/settings-dark.png)](screenshots/settings-dark.png) |
| Playlist with per-video progress | Settings screen |
| [![Ongoing downloads (light theme)](screenshots/main-ongoing-light.png)](screenshots/main-ongoing-light.png) | |
| Ongoing tab — light theme | |

## Features

- **Video or audio** — per-download format and quality selection
- **Per-quality file size estimates** right in the dropdown, for both video
  and audio (e.g. `720p (410.9 MB)`)
- **Paste multiple links at once** (space/comma/newline separated) — each
  downloads independently
- **Playlist support** — detected automatically, every video shown in a
  scrollable list (not just the first few) with per-video progress
- **Pause / resume** any download mid-transfer, including videos still queued
  in a playlist (mark them "Held" before their turn comes up)
- **Retry** a failed or cancelled download without re-adding the link
- **Per-row clear** — dismiss individual rows from the Completed tab, or use
  one-click "Clear All"
- **Parallel downloads** — only the CPU-heavy merge/convert step is serialized
  to avoid ffmpeg conflicts when two videos finish at once
- **Ongoing / Completed tabs** with inline collapsible, auto-scrolling logs
- **One-click dependency fixes** — a built-in checker that installs `yt-dlp`
  and `ffmpeg` via `pkexec` when they're missing
- **Check for Updates** — queries the GitHub releases API in the background
  and offers to **Download & Install** the latest version for your platform
  (the downloaded file is cleaned up afterwards), or open the release page

## Install

Download the latest package for your platform from
[Releases](https://github.com/MDHasan0078/simple-yt-downloader/releases).

| Platform | Package | Install |
|----------|---------|---------|
| Linux | `simple-yt-downloader_<version>_all.deb` | `sudo dpkg -i simple-yt-downloader_*.deb` then `sudo apt install -f -y` |
| Windows | `simple-yt-downloader-<version>-Setup.exe` | Run the installer |
| macOS | `simple-yt-downloader-<version>.dmg` | Open the DMG and drag the app to Applications |

### Dependencies

- `python3`, `python3-gi`, `gir1.2-gtk-3.0` (GTK3 + Python bindings)
- [`yt-dlp`](https://github.com/yt-dlp/yt-dlp) — the actual download engine
- `ffmpeg` — required for merging video+audio streams and audio extraction
- `policykit-1` — used for the in-app "Fix Dependencies" installer prompt

These are declared in `packaging/control`'s `Depends:` line, so a normal
`dpkg -i` + `apt install -f` pulls in anything missing automatically.

## Running from source

```bash
git clone https://github.com/MDHasan0078/simple-yt-downloader.git
cd simple-yt-downloader
python3 run.py
```

Requires `python3-gi` and `gir1.2-gtk-3.0` installed via your distro's package
manager (there's no pip package for GTK3 bindings).

## Building the .deb

```bash
./scripts/build_deb.sh
```

Produces `simple-yt-downloader_<version>_all.deb` in the repo root.

## Project structure

```
simple_yt_downloader/     the actual application (Python + GTK3)
  app.py                main window, headerbar, tabs, multi-URL add flow,
                         first-run flow, Light/Dark theme forcing
  row_widgets.py        VideoRow / PlaylistRow / quality dropdown bar,
                         retry/hold/pause logic
  download_task.py      wraps a single yt-dlp subprocess: probing,
                         format-size estimation, download, pause/resume
                         (process-group signaling), merge serialization
  settings_view.py      the Settings screen (dependency checking +
                         Light/Dark theme options)
  dependencies.py       checks/installs yt-dlp + ffmpeg via pkexec
  config.py             settings persistence (~/.config/simple-yt-downloader)
  style.py              the app's CSS (rounded cards, accent color, etc.)
packaging/              .deb control file, desktop entry, icons
scripts/build_deb.sh    rebuilds the .deb from source
run.py                  dev entry point (no install needed)
```

## Notable implementation details

A few things that took real debugging to get right, worth knowing if you're
reading the source:

- **Pause/resume signals the whole process group**, not just the yt-dlp
  process. yt-dlp commonly hands work off to `ffmpeg` as a real child
  process (merging streams, or as an external downloader for some protocols)
  — signaling only the parent left that child running untouched. See
  `_signal_group()` in `download_task.py`.
- **Concurrent merge conflicts are prevented via a detect-and-gate
  approach**: the moment yt-dlp's output shows any postprocessor step
  starting (`[Merger]`, `[ExtractAudio]`, etc.), that process is frozen
  (same process-group SIGSTOP) until a global lock is free, then resumed.
  Actual downloading stays fully parallel; only the CPU-bound ffmpeg phase
  serializes.
- **Dark/Light theme forces `gtk-theme-name` to `"Adwaita"`** rather than
  just setting the "prefer dark" property. Linux Mint's default theme
  (Mint-Y) doesn't implement that property the way Adwaita does — Mint-Y
  and Mint-Y-Dark are two separate named themes, not one theme with a toggle.
- **The root `Gtk.Stack` has `hhomogeneous`/`vhomogeneous` set to `False`**,
  and Settings is wrapped in its own `Gtk.ScrolledWindow`. Without both of
  these, opening Settings (a tall screen) permanently inflates the whole
  window's size even after navigating back to the main view.
- **Dependency checks run on a background thread.** Checking yt-dlp's
  version means spawning a full Python interpreter as a subprocess — doing
  that synchronously during Settings' `__init__` blocked the screen from
  opening for a couple of seconds.

## Known limitations

- Playlists over 8 videos estimate total size from a sample of the first 8
  (scaled up), rather than querying every single video — shown with a `~`
  prefix.
- Audio quality sizes are estimated (`bitrate × duration`), since a file at
  an arbitrary target bitrate doesn't exist until it's actually encoded.
- Downloads within one playlist run sequentially, not concurrently — an
  intentional choice to avoid rate-limit risk.
- Tested primarily on Linux Mint (Cinnamon); should work on other GTK3-based
  Debian/Ubuntu desktops but unverified there.

## License

MIT — see [LICENSE](LICENSE).

## Credits

Built on [yt-dlp](https://github.com/yt-dlp/yt-dlp) and
[ffmpeg](https://ffmpeg.org/). Not affiliated with YouTube.
