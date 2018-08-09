Gocam is a tool to produce Gcode from DXF files in order to pilot a CNC router. It is inspired by [bCNC](https://github.com/vlachoudis/bCNC/) and aims to work with simple 2.5D models. **This is a very long-term goal**. Right now all it does is converting a DXF file to Gcode.

# Resources

- [offset algorithm for polyline curves](https://seant23.files.wordpress.com/2010/11/anoffsetalgorithm.pdf)
- [divide and conquer NURBS into polylines](https://www.bfft.de/en/techblog-eng/bfft-techblog-mai-divide-and-conquer-nurbs-into-polylines/)
- [time efficient NURBS curve evaluation](https://www.researchgate.net/publication/228411721/download)
- [FindKnotSpan](https://github.com/mfem/mfem/blob/master/mesh/nurbs.cpp#L214)

# TODO

- splines
- piecewise linear approximation
- better UI with options on the commandline

Greville abscissae: mean of k-1 knots. Can be used as starting points to approximate the shape of the spline. Then subdivide the parts until the error is bellow a threshold.

# 2018-08-09

Merged early support for splines. Right now all I have is a function to evaluate a spline, ie I pass a value representing a point between the start and the end of the spline and I get the corresponding coordinates.

As splines are not supported by grbl, my goal is to approximate them with polylines (hopefully lines and arcs) and work with that.

# 2018-08-07

Rewrote a good part of the program. It should behave the same, but maybe faster.

# 2018-02-27

It is possible to use gocam to convert (a subset of) dxf to gcode. Right now 3 entities are supported : line, arc and circle.

    gocam myfile.dxf > myfile.ngc

Keep in mind that the method used to construct paths from the imported entities is still rather crude. I still managed to convert a 5.5Mb DXF file in 20s or so - it's not that bad, I think.
