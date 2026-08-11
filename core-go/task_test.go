package main

import (
	"reflect"
	"testing"
)

// Golden-string tests mirroring simple_yt_downloader/test_download_task.py so
// the Python and Go engines stay behaviorally in lockstep.

func audioTask(mode string, quality string, codec string) *DownloadTask {
	t := newDownloadTask("https://example.com/v", "/tmp", "")
	t.mode = mode
	t.audioQuality = quality
	t.audioCodec = codec
	return t
}

func TestAudioBuildFormatString(t *testing.T) {
	cases := []struct {
		name   string
		task   *DownloadTask
		expect string
	}{
		{"transcode_uses_plain_bestaudio", audioTask("audio", "192", ""), "bestaudio"},
		{"passthrough_prefers_codec_then_falls_back", audioTask("audio", "192", "ec-3"), "bestaudio[acodec=ec-3]/bestaudio"},
	}
	for _, c := range cases {
		if got := c.task.buildFormatString(); got != c.expect {
			t.Errorf("%s: buildFormatString() = %q, want %q", c.name, got, c.expect)
		}
	}
}

func TestAudioPassthroughEveryCodec(t *testing.T) {
	for codecID := range audioPassthroughCodecs {
		task := audioTask("audio", "192", codecID)
		expect := "bestaudio[acodec=" + codecID + "]/bestaudio"
		if got := task.buildFormatString(); got != expect {
			t.Errorf("codec %s: buildFormatString() = %q, want %q", codecID, got, expect)
		}
		if args := task.buildPostprocessArgs(); len(args) != 0 {
			t.Errorf("codec %s: passthrough must not run ffmpeg, got %v", codecID, args)
		}
	}
}

func TestAudioBuildPostprocessArgs(t *testing.T) {
	cases := []struct {
		name   string
		task   *DownloadTask
		expect []string
	}{
		{"transcode_default_quality", audioTask("audio", "192", ""), []string{"-x", "--audio-format", "mp3", "--audio-quality", "192K"}},
		{"transcode_best_quality_maps_to_zero", audioTask("audio", "best", ""), []string{"-x", "--audio-format", "mp3", "--audio-quality", "0"}},
		{"transcode_custom_bitrate", audioTask("audio", "custom:320", ""), []string{"-x", "--audio-format", "mp3", "--audio-quality", "320K"}},
		{"passthrough_skips_postprocess", audioTask("audio", "custom:999", "ac-4"), []string{}},
	}
	for _, c := range cases {
		if got := c.task.buildPostprocessArgs(); !reflect.DeepEqual(got, c.expect) {
			t.Errorf("%s: buildPostprocessArgs() = %v, want %v", c.name, got, c.expect)
		}
	}
}

func TestSplitCustomAudioQuality(t *testing.T) {
	cases := []struct {
		in     string
		expect int
	}{
		{"192", 192},
		{"custom:320", 320},
		{"custom:junk", 192},
		{"custom:", 192},
		{"custom:-5", 192},
		{"custom:0", 192},
		{"", 192},
		{"abc", 192},
	}
	for _, c := range cases {
		if got := splitCustomAudioQuality(c.in); got != c.expect {
			t.Errorf("splitCustomAudioQuality(%q) = %d, want %d", c.in, got, c.expect)
		}
	}
}

func TestVideoFormatStringUnchanged(t *testing.T) {
	task := newDownloadTask("https://example.com/v", "/tmp", "")
	task.videoQuality = "720"
	task.videoFormat = "mp4"
	expect := "bestvideo[height<=720][ext=mp4]+bestaudio/bestvideo[height<=720]+bestaudio/best[height<=720]"
	if got := task.buildFormatString(); got != expect {
		t.Errorf("video buildFormatString() = %q, want %q", got, expect)
	}
}
