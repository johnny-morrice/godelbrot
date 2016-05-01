# godelbrot

## Summary

A Unix-style Mandelbrot set explorer in Go.

## Features

* Designed as Unix toolkit of orthogonal command-line apps
* Provides middleware REST service to manage queuing and concurrency.
* Configuration file generation tool (configbrot)
* Subdividing regions algorithm
* Arbitrary precision mode (and extensible internals)
* Greyscale is default (for integration into an external pipeline)

## Philosophy

* Worse is better
* Few features
* Designed for extensibility (both internally, and in command-line app usage).

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
    $ configbrot | restfulbrot

    $ # Terminal B
    $ # Render through the webservice
    $ configbrot | clientbrot > mandelbrot.png

`restfulbrot` supports a range of options useful for the implementors of viewing clients.

`clientbrot` is a command line client of restfulbrot.

    $ # No configbrot needed when zooming to item stored server-side.
    $ clientbrot --cycle --getrq CJRRiGU_neADTL-GWEoC -xmin 100 -xmax 350 -ymin 100 ymax 280 > img/zoom.png

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
