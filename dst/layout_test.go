package dst

import (
	"testing"
	"time"

	"github.com/recentralized/structure/data"
	"github.com/recentralized/structure/dst/files"
	"github.com/recentralized/structure/meta"
	"github.com/recentralized/structure/uri"
)

func TestFilesystemLayout(t *testing.T) {
	tests := []struct {
		desc        string
		hash        data.Hash
		meta        *meta.Meta
		wantDataURI string
		wantMetaURI string
	}{
		{
			desc: "dated media",
			hash: data.LiteralHash("abcdefg"),
			meta: &meta.Meta{
				Type: data.JPG,
				Inherent: meta.Content{
					Created: time.Date(2015, 1, 2, 9, 9, 9, 9, time.UTC),
				},
				Sidecar: meta.Content{
					Created: time.Date(2011, 1, 2, 9, 9, 9, 9, time.UTC),
				},
			},
			wantDataURI: "media/2015/2015-01-02/abcdefg.jpg",
			wantMetaURI: "meta/ab/cd/efg.json",
		},
		{
			desc: "undated media",
			hash: data.LiteralHash("abcdefg"),
			meta: &meta.Meta{
				Type: data.JPG,
			},
			wantDataURI: "media/Undated/ab/cd/efg.jpg",
			wantMetaURI: "meta/ab/cd/efg.json",
		},
		{
			desc: "unknown class",
			hash: data.LiteralHash("abcdefg"),
			meta: &meta.Meta{
				Type: data.UnknownType,
			},
			wantDataURI: "unknown/ab/cd/efg",
			wantMetaURI: "meta/ab/cd/efg.json",
		},
	}
	for _, tt := range tests {
		layout := NewFilesystemLayout()
		got := layout.DataURI(tt.hash, tt.meta)
		if got, want := got.String(), tt.wantDataURI; got != want {
			t.Errorf("%q DataURI()\ngot  %s\nwant %s", tt.desc, got, want)
		}
		got = layout.MetaURI(tt.hash, tt.meta)
		if got, want := got.String(), tt.wantMetaURI; got != want {
			t.Errorf("%q MetaURI()\ngot  %s\nwant %s", tt.desc, got, want)
		}
	}
}

func TestFilesytemLayoutFiles(t *testing.T) {
	tests := []struct {
		desc     string
		layout   Layout
		readFunc func(string) ([]byte, error)
		data     string
	}{
		{
			desc:   "default",
			layout: NewFilesystemLayout(),
		},
		{
			desc:   "with replaced files.Read",
			layout: NewFilesystemLayout(),
			readFunc: func(name string) ([]byte, error) {
				return []byte(name + "-data"), nil
			},
			data: "fslayout_readme.txt-data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.readFunc != nil {
				var oldRead = files.Read
				defer func() {
					files.Read = oldRead
				}()
				files.Read = tt.readFunc
			}
			files := tt.layout.Files()
			if len(files) == 0 {
				t.Fatalf("expect files")
			}
			file := files[0]
			if got, want := file.URI, uri.TrustedNew("README.txt"); !got.Equal(want) {
				t.Errorf("URI got %s want %s", got, want)
			}
			if len(file.Data) == 0 {
				t.Errorf("File has no data")
			}
			if len(tt.data) > 0 && string(file.Data) != tt.data {
				t.Errorf("File data is not expected, got %s", file.Data)
			}
		})
	}
}
