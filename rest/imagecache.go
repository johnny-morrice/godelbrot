package rest

import (
	"fmt"
)

type imagecache struct {
	ca cache
}

type pic struct {
	png    []byte
	pkthsh hashcode
}

func (p pic) hash() hashcode {
	return p.pkthsh
}

func makeImageCache() imagecache {
	ic := imagecache{}
	ic.ca = makeCache(pic{})
	return ic
}

func (ic imagecache) put(hsh hashcode, png []byte) {
	sz := 0
	if __DEBUG {
		sz = len(png)
	}
	debugf("Storing image for %v (length %v)", hsh, sz)
	p := pic{}
	p.png = png
	p.pkthsh = hsh
	ic.ca.put(p)
}

func (ic imagecache) get(hsh hashcode) ([]byte, bool) {
	debugf("Fetching image for %v", hsh)
	any, present := ic.ca.get(hsh)
	if !present {
		return nil, false
	}
	p, ok := any.(pic)
	if !ok {
		panic(fmt.Sprintf("Expected type pic received: %v", any))
	}
	return p.png, true
}
