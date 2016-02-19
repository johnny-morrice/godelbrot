# godelbrot

## Summary

A Unix-style Mandelbrot set explorer in Go.

## Features

* Designed as Unix toolkit of orthogonal command-line apps
* Configuration file generation tool (configbrot)
* Subdividing regions algorithm
* Arbitrary precision mode (and extensible internals)

## Philosophy

* Worse is better
* Greyscale is default
* Few features
* Designed for extensibility (both internally, and in command-line app usage).

## Get it

    $ go get github.com/johnny-morrice/godelbrot/renderbrot
    $ go get github.com/johnny-morrice/godelbrot/configbrot
    $ go get github.com/johnny-morrice/godelbrot/colorbrot

## Use it

    $ configbrot | renderbrot > mandelbrot.png

configbrot generates configuration files and supports a range of options.  Try

    $ configbrot -help

colorbrot is provided as a convenience for those who may like to recolour the output.

## You might also like

Webdelbrot is a web front-end for Godelbrot.  

http://github.com/johnny-morrice/webdelbrot

## Credits

**John Morrice**

http://functorama.com

https://github.com/johnny-morrice

**Gavin Leech**

https://github.com/technicalities

## License

We use an MIT style license.  See LICENSE.txt for terms of use and distribution.
