# godelbrot

## Summary

A Fancy Unix-style Mandelbrot set explorer in Go

## Dependencies

godelbrot has no dependencies outside of the go compiler, and its standard
library.

The godelbrot developers are using go 1.4, mostly because godebug does not
support 1.5 at the date of publication.

## Get it

    $ go get github.com/johnny-morrice/godelbrot/renderbrot
    $ go get github.com/johnny-morrice/godelbrot/configbrot
    $ go get github.com/johnny-morrice/godelbrot/colorbrot

## Use it

    $ configbrot | renderbrot > mandelbrot.png

configbrot generates configuration files and supports a range of options.  Try

    $ configbrot -help

colorbrot is provided as a convenience for those who may like to recolour the output.

## Credits

**John Morrice**

http://functorama.com

https://github.com/johnny-morrice

**Gavin Leech**

https://github.com/technicalities
