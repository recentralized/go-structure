package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/recentralized/structure/cid"
	"github.com/recentralized/structure/content"
	"github.com/recentralized/structure/dst"
	"github.com/recentralized/structure/index"
	"github.com/recentralized/structure/meta"
	"github.com/recentralized/structure/uri"
)

func main() {
	index, err := buildIndex()
	if err != nil {
		fmt.Printf("Failed to build index: %s", err)
		os.Exit(1)
	}

	locator := dst.NewFilesystemLocator()

	err = addRefs(locator, index)
	if err != nil {
		fmt.Printf("Failed to add refs: %s", err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		fmt.Printf("Failed to create json: %s", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}

func buildIndex() (*index.Index, error) {
	srcPath, err := uri.NewDirPath("/tmp/src")
	if err != nil {
		return nil, fmt.Errorf("Could not create src path: %s", err)
	}

	dstPath, err := uri.NewDirPath("/tmp/dst")
	if err != nil {
		return nil, fmt.Errorf("Could not create dst path: %s", err)
	}

	src := index.NewSrc(srcPath.URI)
	dst := index.NewDstAllAt(dstPath.URI)

	idx := index.New()
	idx.Srcs = []index.Src{src}
	idx.Dsts = []index.Dst{dst}

	return idx, nil
}

func addRefs(loc dst.Locator, idx *index.Index) error {
	src := idx.Srcs[0]
	dst := idx.Dsts[0]

	data := []byte("fictional image data")
	cid, err := loc.NewHash(bytes.NewReader(data))
	if err != nil {
		return err
	}

	srcItem, meta, err := buildSrcItem(src)
	if err != nil {
		return err
	}

	dstItem, err := buildDstItem(loc, dst, cid, meta)
	if err != nil {
		return err
	}

	idx.AddRef(index.Ref{
		Hash: cid,
		Src:  srcItem,
		Dst:  dstItem,
	})

	return nil
}

func buildSrcItem(src index.Src) (index.SrcItem, *meta.Meta, error) {
	var item index.SrcItem

	dataPath := uri.TrustedNew("fictional/image.jpg")
	dataURI, err := src.SrcURI.ResolveReference(dataPath)
	if err != nil {
		return item, nil, err
	}

	metaPath := uri.TrustedNew("fictional/image.xmp")
	metaURI, err := src.SrcURI.ResolveReference(metaPath)
	if err != nil {
		return item, nil, err
	}

	item = index.SrcItem{
		SrcID:      src.SrcID,
		DataURI:    dataURI,
		MetaURI:    metaURI,
		ModifiedAt: time.Date(2018, 11, 12, 0, 0, 0, 0, time.UTC),
	}

	doc := meta.New()
	doc.ContentType = content.JPG
	doc.Inherent = meta.Content{
		Created: time.Date(2018, 11, 10, 0, 0, 0, 0, time.UTC),
	}

	return item, doc, nil
}

func buildDstItem(loc dst.Locator, dst index.Dst, cid cid.ContentID, meta *meta.Meta) (index.DstItem, error) {
	var item index.DstItem

	dataURI := loc.DataURI(cid, meta)
	metaURI := loc.MetaURI(cid, meta)

	item = index.DstItem{
		DstID:    dst.DstID,
		DataURI:  dataURI,
		MetaURI:  metaURI,
		StoredAt: time.Date(2018, 11, 13, 0, 0, 0, 0, time.UTC),
	}
	return item, nil
}
