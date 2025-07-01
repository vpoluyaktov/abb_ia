package utils

import "testing"

func TestSanitizeFilePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"empty string", "", ""},
		{"no special chars", "normal/path/file.txt", "normal/path/file.txt"},
		{"single quotes", "user's/file.txt", "users/file.txt"},
		{"multiple dots", "file...txt", "file.txt"},
		{"double dots", "file..txt", "file.txt"},
		{"dollar sign", "file$name.txt", "file.name.txt"},
		{"exclamation mark", "important!file.txt", "important.file.txt"},
		{"question mark", "what?file.txt", "what.file.txt"},
		{"ellipsis", "file…name.txt", "filename.txt"},
		{"hash", "file#1.txt", "fileN1.txt"},
		{"brackets", "[file].txt", "file.txt"},
		{"colon", "file:name.txt", "file.name.txt"},
		{"multiple replacements", "[file]!name#1...txt", "file.nameN1.txt"},
		{"sequential replacements", "file!!name", "file.name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeFilePath(tt.path)
			if got != tt.want {
				t.Errorf("SanitizeFilePath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestSanitizeMp3FileName(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{"empty string", "", ""},
		{"no bitrate", "song.mp3", "song.mp3"},
		{"64kb", "song_64kb.mp3", "song.mp3"},
		{"24kb", "song_24kb.mp3", "song.mp3"},
		{"32kb", "song_32kb.mp3", "song.mp3"},
		{"40kb", "song_40kb.mp3", "song.mp3"},
		{"48kb", "song_48kb.mp3", "song.mp3"},
		{"56kb", "song_56kb.mp3", "song.mp3"},
		{"80kb", "song_80kb.mp3", "song.mp3"},
		{"96kb", "song_96kb.mp3", "song.mp3"},
		{"112kb", "song_112kb.mp3", "song.mp3"},
		{"128kb", "song_128kb.mp3", "song.mp3"},
		{"144kb", "song_144kb.mp3", "song.mp3"},
		{"160kb", "song_160kb.mp3", "song.mp3"},
		{"224kb", "song_224kb.mp3", "song.mp3"},
		{"256kb", "song_256kb.mp3", "song.mp3"},
		{"320kb", "song_320kb.mp3", "song.mp3"},
		{"vbr", "song_vbr.mp3", "song.mp3"},
		{"multiple bitrates", "song_128kb_vbr.mp3", "song.mp3"},
		{"bitrate in middle", "my_128kb_song.mp3", "my_song.mp3"},
		{"multiple occurrences", "song_128kb_128kb.mp3", "song.mp3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeMp3FileName(tt.fileName)
			if got != tt.want {
				t.Errorf("SanitizeMp3FileName(%q) = %q, want %q", tt.fileName, got, tt.want)
			}
		})
	}
}
