Gocam aims to produce Gcode from DXF files in order to pilot a CNC router. It is inspired by [bCNC](https://github.com/vlachoudis/bCNC/) and aims to work with simple 2.5D models. **This is a very long-term goal**. Right now all it does is converting a DXF file to Gcode.

# TODO

- splines
- piecewise linear approximation
- better UI with options on the commandline

# 2018-02-27

It is possible to use gocam to convert (a subset of) dxf to gcode. Right now 3 entities are supported : line, arc and circle.

    gocam myfile.dxf > myfile.ngc

Keep in mind that the method used to construct paths from the imported entities is still rather crude. I still managed to convert a 5.5Mb DXF file in 20s or so - it's not that bad, I think.
