# godelbrot

## Summary

A Unix-style Mandelbrot set explorer in Go.

## Demo - CHROME ONLY

Better browser support en route: see [webdelbrot issue #9](https://github.com/johnny-morrice/webdelbrot/issues/9).

[Webdelbrot client](http://godelbrot.functorama.com)

Webdelbrot is a gopherjs web client for Godelbrot.  Godelbrot is a fractal
render backend designed to allow easy implementation of GUI clients.

## Features

* Designed as Unix toolkit of orthogonal command-line apps
* Provides middleware REST service to manage queuing and concurrency
* Zoom into image given pixel boundaries (no math for you! :)
* Configuration file generation tool (`configbrot`)
* Subdividing regions algorithm
* Arbitrary precision mode (and extensible internals)
* Greyscale is default (for integration into an external pipeline)

## Philosophy

* Worse is better
* Few features
* Designed for extensibility (both internally, and in command-line app usage)
* Ease of client implementation

## Get it

    $ go get github.com/johnny-morrice/godelbrot
    $ go install -tags production github.com/johnny-morrice/godelbrot...

As Godelbrot is a multiple-binary toolkit so the latter command installs all other binaries.

## Use it

    $ godelbrot > mandelbrot.png

Which is equivalent to:

    $ configbrot | renderbrot > mandelbrot.png

`configbrot` generates configuration files and supports a range of options.  Try

    $ configbrot -help

For a persisent process try `restfulbrot`

    $ # Terminal A
    $ # Run fractal webservice
    $ # Options given to configbrot here become defaults for all render clients.
    $ configbrot -palette pretty | restfulbrot

    $ # Terminal B
    $ # Render through the webservice
    $ # Customize defaults with configbrot options
    $ configbrot -width 1920 -height 1080 | clientbrot > mandelbrot.png

`restfulbrot` supports a range of options useful for the implementors of viewing clients.

`clientbrot` is a command line client of restfulbrot.

    $ # No configbrot needed when zooming to item stored server-side.
    $ clientbrot --cycle --getrq CJRRiGU_neADTL-GWEoC -xmin 100 -xmax 350 -ymin 100 -ymax 280 > img/zoom.png

`colorbrot` is provided as a convenience for those who may like to recolour the output.

## You might also like

Webdelbrot is a web front-end for Godelbrot.

http://github.com/johnny-morrice/webdelbrot

## Credits

**John Morrice**

http://functorama.com

https://github.com/johnny-morrice

See HISTORY.md for contributors

## License

We use an MIT style license.  See LICENSE.txt for terms of use and distribution.
