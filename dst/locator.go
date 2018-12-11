package dst

import (
	"fmt"
	"io"

	"github.com/recentralized/structure/data"
	"github.com/recentralized/structure/meta"
	"github.com/recentralized/structure/uri"
)

// Locator is the interface for defining the locations of data in a
// destination. Locators always return relative URLs, which can be resolved
// with the destination's base URIs when needed.
type Locator interface {

	// NewHash generates the hash for data.
	NewHash(io.Reader) (data.Hash, error)

	// IndexURI returns the document that stores the index.
	IndexURI() uri.URI

	// RefsURI returns the document that should store this ref. This is
	// generally the same as IndexURI, but returning a different value
	// would allow you to shard the refs.
	RefsURI(data.Hash) uri.URI

	// DataURI returns the location that this data should be stored.
	DataURI(data.Hash, *meta.Meta) uri.URI

	// MetaURI returns the location that this meta should be stored.
	MetaURI(data.Hash, *meta.Meta) uri.URI
}

// NewFilesystemLocator initializes the standard locator for use on filesystems
// and filesystem-like storage media such as AWS S3.
func NewFilesystemLocator() Locator {
	return fsLocator{
		indexFile: "index.json",
		classToCategory: map[data.Class]string{
			data.Image: "media",
		},
		unknownCategory: "unknown",
		zeroDateDir:     "Undated",
	}
}

type fsLocator struct {
	indexFile       string
	classToCategory map[data.Class]string
	unknownCategory string
	zeroDateDir     string
}

func (l fsLocator) NewHash(r io.Reader) (data.Hash, error) {
	return data.NewHash(r)
}

func (l fsLocator) IndexURI() uri.URI {
	return uri.TrustedNew(l.indexFile)
}

func (l fsLocator) RefsURI(hash data.Hash) uri.URI {
	return l.IndexURI()
}

// media/2006/2006-01-02/<hash>.<ext>
// media/Undated/hash(<hash>)/<hash>.<ext>
// <category>/hash(<hash>)/<hash>.<ext>
func (l fsLocator) DataURI(hash data.Hash, meta *meta.Meta) uri.URI {
	var (
		key string
		ext = meta.Type.Ext()
		cls = meta.Type.Class()
	)

	// Categorize by the class of data.
	category := l.classToCategory[cls]
	if category == "" {
		category = l.unknownCategory
	}

	// Customize the path location for each category.
	switch category {

	// "media" category names files by hash and organized by date.
	case "media":
		t := meta.DateCreated()
		if t.IsZero() {
			key = fmt.Sprintf("%s/%s/%s.%s", category, l.zeroDateDir, l.dirs(hash), ext)
			return uri.TrustedNew(key)
		}
		year := t.Format("2006")
		date := t.Format("2006-01-02")
		key = fmt.Sprintf("%s/%s/%s/%s.%s", category, year, date, hash.String(), ext)
		return uri.TrustedNew(key)

	// Otherwise organize by breaking down the hash.
	default:
		if ext == "" {
			key = fmt.Sprintf("%s/%s", category, l.dirs(hash))
		} else {
			key = fmt.Sprintf("%s/%s.%s", category, l.dirs(hash), ext)
		}
		return uri.TrustedNew(key)
	}
}

// meta/hash(<hash>)/<hash>.json
func (l fsLocator) MetaURI(hash data.Hash, meta *meta.Meta) uri.URI {
	key := fmt.Sprintf("meta/%s.%s", l.dirs(hash), "json")
	return uri.TrustedNew(key)
}

func (l fsLocator) dirs(hash data.Hash) string {
	s := hash.String()
	if len(s) > 4 {
		return fmt.Sprintf("%s/%s/%s", s[0:2], s[2:4], s[4:])
	}
	return s
}
