package stbtt

import "math"

type edge struct {
	x0, y0, x1, y1 float32
	invert         bool
}

type activeEdge struct {
	next             *activeEdge
	fx, fy, fdx, fdy float32
	direction        float32
	sy, ey           float32
}

func newActive(e *edge, off_x int32, start_point float32) *activeEdge {
	var z = new(activeEdge)
	var dxdy = (e.x1 - e.x0) / (e.y1 - e.y0)
	z.fdx = dxdy
	if dxdy != 0.0 {
		z.fdy = 1.0 / dxdy
	} else {
		z.fdy = 0.0
	}
	z.fx = e.x0 + dxdy*(start_point-e.y0)
	z.fx -= float32(off_x)
	if e.invert {
		z.direction = 1.0
	} else {
		z.direction = -1.0
	}
	z.sy = e.y0
	z.ey = e.y1
	return z
}

// the edge passed in here does not cross the vertical line at x or the vertical line at x+1
// (i.e. it has already been clipped to those)
func handleClippedEdge(scanline []float32, x int32, e *activeEdge, x0, y0, x1, y1 float32) {
	if y0 == y1 {
		return
	}
	if !(y0 < y1) {
		panic("!(y0 < y1)")
	}
	if !(e.sy <= e.ey) {
		panic("!(e.sy <= e.ey)")
	}
	if y0 > e.ey {
		return
	}
	if y1 < e.sy {
		return
	}
	if y0 < e.sy {
		x0 += (x1 - x0) * (e.sy - y0) / (y1 - y0)
		y0 = e.sy
	}
	if y1 > e.ey {
		x1 += (x1 - x0) * (e.ey - y1) / (y1 - y0)
		y1 = e.ey
	}

	xf := float32(x)

	if x0 == xf {
		if !(x1 <= xf+1) {
			panic("!(x1 <= x+1)")
		}
	} else if x0 == xf+1 {
		if !(x1 >= xf) {
			panic("!(x1 >= x)")
		}
	} else if x0 <= xf {
		if !(x1 <= xf) {
			panic("!(x1 <= x)")
		}
	} else if x0 >= xf+1 {
		if !(x1 >= xf+1) {
			panic("!(x1 >= x+1)")
		}
	} else {
		if !(x1 >= xf && x1 <= xf+1) {
			panic("!(x1 >= x && x1 <= x+1)")
		}
	}

	if x0 <= xf && x1 <= xf {
		scanline[int32(x)] += e.direction * (y1 - y0)
	} else if x0 >= xf+1 && x1 >= xf+1 {
		// skip
	} else {
		if !(x0 >= xf && x0 <= xf+1 && x1 >= xf && x1 <= xf+1) {
			panic("!(x0 >= x && x0 <= x+1 && x1 >= x && x1 <= x+1)")
		}

		scanline[int32(x)] += e.direction * (y1 - y0) * (1 - ((x0-xf)+(x1-xf))/2)
	}
}

func sizedTrapezoidArea(height, top_width, bottom_width float32) float32 {
	if !(top_width >= 0) {
		panic("!(top_width >= 0)")
	}
	if !(bottom_width >= 0) {
		panic("!(bottom_width >= 0)")
	}
	return (top_width + bottom_width) / 2.0 * height
}

func positionTrapezoidArea(height, tx0, tx1, bx0, bx1 float32) float32 {
	return sizedTrapezoidArea(height, tx1-tx0, bx1-bx0)
}

func sizedTriangleArea(height, width float32) float32 {
	return width * height / 2.0
}

func fillActiveEdges(scanline, scanline_fill []float32, scanline_fill_idx int32, len int32, e *activeEdge, y_top float32) {
	var y_bottom = y_top + 1

	for e != nil {
		// brute force every pixel

		// compute intersection points with top & bottom
		if !(e.ey >= y_top) {
			panic("!(e->ey >= y_top)")
		}

		if e.fdx == 0 {
			x0 := e.fx
			if x0 < float32(len) {
				if x0 >= 0 {
					handleClippedEdge(scanline, int32(x0), e, x0, y_top, x0, y_bottom)
					handleClippedEdge(scanline_fill, scanline_fill_idx-1, e, x0, y_top, x0, y_bottom)
				} else {
					handleClippedEdge(scanline_fill, scanline_fill_idx-1, e, x0, y_top, x0, y_bottom)
				}
			}
		} else {
			var x0 = e.fx
			var dx = e.fdx
			var xb = x0 + dx
			var x_top, x_bottom float32
			var sy0, sy1 float32
			var dy = e.fdy

			if !(e.sy <= y_bottom && e.ey >= y_top) {
				panic("!(e->sy <= y_bottom && e->ey >= y_top)")
			}

			// compute endpoints of line segment clipped to this scanline (if the
			// line segment starts on this scanline. x0 is the intersection of the
			// line with y_top, but that may be off the line segment.
			if e.sy > y_top {
				x_top = x0 + dx*(e.sy-y_top)
				sy0 = e.sy
			} else {
				x_top = x0
				sy0 = y_top
			}
			if e.ey < y_bottom {
				x_bottom = x0 + dx*(e.ey-y_top)
				sy1 = e.ey
			} else {
				x_bottom = xb
				sy1 = y_bottom
			}

			if x_top >= 0 && x_bottom >= 0 && x_top < float32(len) && x_bottom < float32(len) {
				// from here on, we don't have to range check x values

				if int32(x_top) == int32(x_bottom) {
					var height float32
					// simple case, only spans one pixel
					var x = int32(x_top)
					height = (sy1 - sy0) * e.direction
					if !(x >= 0 && x < len) {
						panic("!(x >= 0 && x < len)")
					}

					scanline[x] += positionTrapezoidArea(height, x_top, float32(x)+1, x_bottom, float32(x)+1)
					scanline_fill[scanline_fill_idx+x] += height
				} else {
					var x, x1, x2 int32
					var y_crossing, y_final, step, sign, area float32
					// covers 2+ pixels
					if x_top > x_bottom {
						// flip scanline vertically; signed area is the same
						var t float32
						sy0 = y_bottom - (sy0 - y_top)
						sy1 = y_bottom - (sy1 - y_top)
						t = sy0
						sy0 = sy1
						sy1 = t
						t = x_bottom
						x_bottom = x_top
						x_top = t
						dx = -dx
						dy = -dy
						t = x0
						x0 = xb
						xb = t
					}
					if !(dy >= 0) {
						panic("!(dy >= 0)")
					}
					if !(dx >= 0) {
						panic("!(dx >= 0)")
					}

					x1 = int32(x_top)
					x2 = int32(x_bottom)
					// compute intersection with y axis at x1+1
					y_crossing = y_top + dy*(float32(x1)+1-x0)

					// compute intersection with y axis at x2
					y_final = y_top + dy*(float32(x2)-x0)

					//           x1    x_top                            x2    x_bottom
					//     y_top  +------|-----+------------+------------+--------|---+------------+
					//            |            |            |            |            |            |
					//            |            |            |            |            |            |
					//       sy0  |      Txxxxx|............|............|............|............|
					// y_crossing |            *xxxxx.......|............|............|............|
					//            |            |     xxxxx..|............|............|............|
					//            |            |     /-   xx*xxxx........|............|............|
					//            |            | dy <       |    xxxxxx..|............|............|
					//   y_final  |            |     \-     |          xx*xxx.........|............|
					//       sy1  |            |            |            |   xxxxxB...|............|
					//            |            |            |            |            |            |
					//            |            |            |            |            |            |
					//  y_bottom  +------------+------------+------------+------------+------------+
					//
					// goal is to measure the area covered by '.' in each pixel

					// if x2 is right at the right edge of x1, y_crossing can blow up, github #1057
					// @TODO: maybe test against sy1 rather than y_bottom?
					if y_crossing > y_bottom {
						y_crossing = y_bottom
					}

					sign = e.direction
					// area of the rectangle covered from sy0..y_crossing
					area = sign * (y_crossing - sy0)

					// area of the triangle (x_top,sy0), (x1+1,sy0), (x1+1,y_crossing)
					scanline[x1] += sizedTriangleArea(area, float32(x1)+1-x_top)

					// check if final y_crossing is blown up; no test case for this
					if y_final > y_bottom {
						y_final = y_bottom
						dy = (y_final - y_crossing) / (float32(x2) - float32(x1+1)) // if denom=0, y_final = y_crossing, so y_final <= y_bottom
					}

					// in second pixel, area covered by line segment found in first pixel
					// is always a rectangle 1 wide * the height of that line segment; this
					// is exactly what the variable 'area' stores. it also gets a contribution
					// from the line segment within it. the THIRD pixel will get the first
					// pixel's rectangle contribution, the second pixel's rectangle contribution,
					// and its own contribution. the 'own contribution' is the same in every pixel except
					// the leftmost and rightmost, a trapezoid that slides down in each pixel.
					// the second pixel's contribution to the third pixel will be the
					// rectangle 1 wide times the height change in the second pixel, which is dy.

					step = sign * dy * 1 // dy is dy/dx, change in y for every 1 change in x,
					// which multiplied by 1-pixel-width is how much pixel area changes for each step in x
					// so the area advances by 'step' every time

					for x = x1 + 1; x < x2; x++ {
						scanline[x] += area + step/2 // area of trapezoid is 1*step/2
						area += step
					}
					if !(math.Abs(float64(area)) <= 1.01) {
						panic("!(math.Abs(float64(area))  <= 1.01)")
					}

					if !(sy1 > y_final-0.01) {
						panic("!(sy1 > y_final-0.01)")
					}

					// area covered in the last pixel is the rectangle from all the pixels to the left,
					// plus the trapezoid filled by the line segment in this pixel all the way to the right edge
					scanline[x2] += area + sign*positionTrapezoidArea(sy1-y_final, float32(x2), float32(x2)+1.0, x_bottom, float32(x2)+1.0)

					// the rest of the line is filled based on the total height of the line segment in this pixel
					scanline_fill[x2] += sign * (sy1 - sy0)
				}
			} else {
				// if edge goes outside of box we're drawing, we require
				// clipping logic. since this does not match the intended use
				// of this library, we use a different, very slow brute
				// force implementation
				// note though that this does happen some of the time because
				// x_top and x_bottom can be extrapolated at the top & bottom of
				// the shape and actually lie outside the bounding box
				var x int32
				for x = 0; x < len; x++ {
					// cases:
					//
					// there can be up to two intersections with the pixel. any intersection
					// with left or right edges can be handled by splitting into two (or three)
					// regions. intersections with top & bottom do not necessitate case-wise logic.
					//
					// the old way of doing this found the intersections with the left & right edges,
					// then used some simple logic to produce up to three segments in sorted order
					// from top-to-bottom. however, this had a problem: if an x edge was epsilon
					// across the x border, then the corresponding y position might not be distinct
					// from the other y segment, and it might ignored as an empty segment. to avoid
					// that, we need to explicitly produce segments based on x positions.

					// rename variables to clearly-defined pairs
					var y0 = y_top
					var x1 = float32(x)
					var x2 = float32(x + 1)
					var x3 = xb
					var y3 = y_bottom

					// x = e->x + e->dx * (y-y_top)
					// (y-y_top) = (x - e->x) / e->dx
					// y = (x - e->x) / e->dx + y_top
					var y1 = (float32(x)-x0)/dx + y_top
					var y2 = (float32(x)+1-x0)/dx + y_top

					if x0 < x1 && x3 > x2 { // three segments descending down-right
						handleClippedEdge(scanline, x, e, x0, y0, x1, y1)
						handleClippedEdge(scanline, x, e, x1, y1, x2, y2)
						handleClippedEdge(scanline, x, e, x2, y2, x3, y3)
					} else if x3 < x1 && x0 > x2 { // three segments descending down-left
						handleClippedEdge(scanline, x, e, x0, y0, x2, y2)
						handleClippedEdge(scanline, x, e, x2, y2, x1, y1)
						handleClippedEdge(scanline, x, e, x1, y1, x3, y3)
					} else if x0 < x1 && x3 > x1 { // two segments across x, down-right
						handleClippedEdge(scanline, x, e, x0, y0, x1, y1)
						handleClippedEdge(scanline, x, e, x1, y1, x3, y3)
					} else if x3 < x1 && x0 > x1 { // two segments across x, down-left
						handleClippedEdge(scanline, x, e, x0, y0, x1, y1)
						handleClippedEdge(scanline, x, e, x1, y1, x3, y3)
					} else if x0 < x2 && x3 > x2 { // two segments across x+1, down-right
						handleClippedEdge(scanline, x, e, x0, y0, x2, y2)
						handleClippedEdge(scanline, x, e, x2, y2, x3, y3)
					} else if x3 < x2 && x0 > x2 { // two segments across x+1, down-left
						handleClippedEdge(scanline, x, e, x0, y0, x2, y2)
						handleClippedEdge(scanline, x, e, x2, y2, x3, y3)
					} else { // one segment
						handleClippedEdge(scanline, x, e, x0, y0, x3, y3)
					}
				}
			}
		}
		e = e.next
	}
}

func rasterizeSortedEdges(result *bitmap, e []edge, n, vsubsample, off_x, off_y int32) {
	var active *activeEdge
	var y, j, i int32
	var scanline_data [129]float32
	var scanline, scanline2 []float32

	if result.w > 64 {
		scanline = make([]float32, result.w*2+1)
	} else {
		scanline = scanline_data[:]
	}
	scanline2 = scanline[result.w:]

	y = off_y
	e[n].y0 = float32(off_y+result.h) + 1
	for ; j < result.h; j, y = j+1, y+1 {
		// find center of pixel for this scanline
		var scan_y_top = float32(y) + 0.0
		var scan_y_bottom = float32(y) + 1.0
		var step = &active

		// update all active edges;
		// remove all active edges that terminate before the top of this scanline
		for *step != nil {
			var z = *step
			if z.ey <= scan_y_top {
				*step = z.next // delete from list
				z.direction = 0
			} else {
				step = &((*step).next) // advance through list
			}
		}

		// insert all edges that start before the bottom of this scanline
		for ; e[0].y0 <= scan_y_bottom; e = e[1:] {
			if e[0].y0 != e[0].y1 {
				var z = newActive(&e[0], off_x, scan_y_top)
				if j == 0 && off_y != 0 {
					if z.ey < scan_y_top {
						// this can happen due to subpixel positioning and some kind of fp rounding error i think
						z.ey = scan_y_top
					}
				}
				z.next = active
				active = z
			}

			if active != nil {
				fillActiveEdges(scanline, scanline2, 1, result.w, active, scan_y_top)
			}
			{
				var sum float32
				for i = 0; i < result.w; i++ {
					sum += scanline2[i]
					var k = scanline[i] + sum
					k = float32(math.Abs(float64(k)))*255 + 0.5
					var m = int32(k)
					if m > 255 {
						m = 255
					}
					result.pixels[j*result.stride+i] = uint8(m)
				}
			}

			// advance all the edges
			step = &active
			for ; *step != nil; step = &((*step).next) {
				(*step).fx += (*step).fdx
			}
		}
	}
}

func compare(a, b *edge) bool {
	return a.y0 < b.y0
}

func sortEdgesInsSort(p []edge) {
	for i := 1; i < len(p); i++ {
		var t = p[i]
		var a = &t
		j := i
		for ; j > 0; j-- {
			var b = &p[j-1]
			var c = compare(a, b)
			if !c {
				break
			}
			p[j] = p[j-1]
		}
		if i != j {
			p[j] = t
		}
	}
}

func sortEdgesQuickSort(p []edge) {
	/* threshold for transitioning to insertion sort */
	for n := len(p); n > 12; n-- {
		var m = n >> 1
		var c01 = compare(&p[0], &p[m])
		var c12 = compare(&p[m], &p[n-1])
		/* if 0 >= mid >= end, or 0 < mid < end, then use mid */
		if c01 != c12 {
			/* otherwise, we'll need to swap something else to middle */
			var z int
			var c = compare(&p[0], &p[n-1])
			/* 0>mid && mid<n:  0>n => n; 0<n => 0 */
			/* 0<mid && mid>n:  0>n => 0; 0<n => n */
			if c == c12 {
				z = 0
			} else {
				z = n - 1
			}
			var t = p[z]
			p[z] = p[m]
			p[m] = t
		}
		/* now p[m] is the median-of-three */
		/* swap it to the beginning so it won't move around */
		var t = p[0]
		p[0] = p[m]
		p[m] = t
		/* partition loop */
		var i = 1
		var j = n - 1
		for {
			/* handling of equality is crucial here */
			/* for sentinels & efficiency with duplicates */
			for {
				if !compare(&p[i], &p[0]) {
					break
				}
				i++
			}
			for {
				if !compare(&p[0], &p[j]) {
					break
				}
				j--
			}
			/* make sure we haven't crossed */
			if i >= j {
				break
			}
			var t = p[i]
			p[i] = p[j]
			p[j] = t
			i++
			j--
		}
		/* recurse on smaller side, iterate on larger */
		if j < (n - i) {
			sortEdgesQuickSort(p[:j])
			p = p[i:]
			n = n - i
		} else {
			sortEdgesQuickSort(p[i:])
			n = j
		}
	}
}

func sortEdges(p []edge) {
	sortEdgesQuickSort(p)
	sortEdgesInsSort(p)
}

func rasterize(
	result *bitmap,
	pts []point,
	windings []int32,
	scale_x, scale_y,
	shift_x, shift_y float32,
	off_x, off_y int32, invert bool) {

	var y_scale_inv = scale_y
	if invert {
		y_scale_inv = -scale_y
	}
	var e []edge
	var n, j, k, m int32

	var vsubsample float32 = 1
	// vsubsample should divide 255 evenly; otherwise we won't reach full opacity

	// now we have to blow out the windings into explicit edge lists
	for i := range windings {
		n += windings[i]
	}

	e = make([]edge, n+1)
	for i := range windings {
		var p = pts[m:]
		m += windings[i]
		j = windings[i] - 1
		for k = 0; k < windings[i]; j = k {
			var a = k
			var b = j
			// skip the edge if horizontal
			if p[j].y == p[k].y {
				continue
			}
			// add edge from j to k to the list
			e[n].invert = false
			if (invert && p[j].y > p[k].y) || (!invert && p[j].y < p[k].y) {
				e[n].invert = true
				a = j
				b = k
			}
			e[n].x0 = p[a].x*scale_x + shift_x
			e[n].y0 = (p[a].y*y_scale_inv + shift_y) * vsubsample
			e[n].x1 = p[b].x*scale_x + shift_x
			e[n].y1 = (p[b].y*y_scale_inv + shift_y) * vsubsample
			n++
			k++
		}
	}

	// now sort the edges by their highest point (should snap to integer, and then by x)
	sortEdges(e[:n])

	// now, traverse the scanlines and find the intersections on each scanline, use xor winding rule
	rasterizeSortedEdges(result, e, n, int32(vsubsample), off_x, off_y)
}

func addPoint(points []point, n int32, x, y float32) {
	if points != nil {
		points[n].x = x
		points[n].y = y
	}
}

func tesselateCurve(points []point, num_points *int32, x0, y0, x1, y1, x2, y2,
	objspace_flatness_squared float32, n int32) int32 {

	// midpoint
	var mx = (x0 + 2*x1 + x2) / 4
	var my = (y0 + 2*y1 + y2) / 4
	// versus directly drawn line
	var dx = (x0+x2)/2 - mx
	var dy = (y0+y2)/2 - my
	if n > 16 { // 65536 segments on one curve better be enough!
		return 1
	}
	if dx*dx+dy*dy > objspace_flatness_squared { // half-pixel error allowed... need to be smaller if AA
		tesselateCurve(points, num_points, x0, y0, (x0+x1)/2.0, (y0+y1)/2.0, mx, my, objspace_flatness_squared, n+1)
		tesselateCurve(points, num_points, mx, my, (x1+x2)/2.0, (y1+y2)/2.0, x2, y2, objspace_flatness_squared, n+1)
	} else {
		addPoint(points, n, x2, y2)
		*num_points = *num_points + 1
	}
	return 1
}

func tesselateCubic(points []point, num_points *int32, x0, y0, x1, y1, x2, y2, x3, y3, objspace_flatness_squared float32, n int) {
	// @TODO this "flatness" calculation is just made-up nonsense that seems to work well enough
	var dx0 = x1 - x0
	var dy0 = y1 - y0
	var dx1 = x2 - x1
	var dy1 = y2 - y1
	var dx2 = x3 - x2
	var dy2 = y3 - y2
	var dx = x3 - x0
	var dy = y3 - y0
	var longlen = float32(math.Sqrt(float64(dx0*dx0+dy0*dy0)) + math.Sqrt(float64(dx1*dx1+dy1*dy1)) + math.Sqrt(float64(dx2*dx2+dy2*dy2)))
	var shortlen = float32(math.Sqrt(float64(dx*dx + dy*dy)))
	var flatness_squared = longlen*longlen - shortlen*shortlen

	if n > 16 { // 65536 segments on one curve better be enough!
		return
	}

	if flatness_squared > objspace_flatness_squared {
		var x01 = (x0 + x1) / 2
		var y01 = (y0 + y1) / 2
		var x12 = (x1 + x2) / 2
		var y12 = (y1 + y2) / 2
		var x23 = (x2 + x3) / 2
		var y23 = (y2 + y3) / 2

		var xa = (x01 + x12) / 2
		var ya = (y01 + y12) / 2
		var xb = (x12 + x23) / 2
		var yb = (y12 + y23) / 2

		var mx = (xa + xb) / 2
		var my = (ya + yb) / 2

		tesselateCubic(points, num_points, x0, y0, x01, y01, xa, ya, mx, my, objspace_flatness_squared, n+1)
		tesselateCubic(points, num_points, mx, my, xb, yb, x23, y23, x3, y3, objspace_flatness_squared, n+1)
	} else {
		addPoint(points, *num_points, x3, y3)
		*num_points = *num_points + 1
	}
}

// returns number of contours, gotta flatten the curve!
func flattenCurves(vertices []Vertex, objspace_flatness float32) (points []point, lengths []int32) {
	var num_points int32
	var objspace_flatness_squared = objspace_flatness * objspace_flatness
	var n, start, pass int32

	// count how many "moves" there are to get the contour count
	for i := range vertices {
		if vertices[i].typ == vmove {
			n++
		}
	}

	if n == 0 {
		return nil, nil
	}

	lengths = make([]int32, n)

	// make two passes through the points so we don't need to realloc
	for pass = 0; pass < 2; pass++ {
		var x, y float32
		if pass == 1 {
			points = make([]point, num_points)
		}
		num_points = 0
		n = -1
		for i := range vertices {
			switch vertices[i].typ {
			case vmove:
				// start the next contour
				if n >= 0 {
					lengths[n] = num_points - start
				}
				n++
				start = num_points
				x = float32(vertices[i].x)
				y = float32(vertices[i].y)
				num_points++
				addPoint(points, num_points, x, y)
			case vline:
				x = float32(vertices[i].x)
				y = float32(vertices[i].y)
				num_points++
				addPoint(points, num_points, x, y)
			case vcurve:
				tesselateCurve(points, &num_points, x, y,
					float32(vertices[i].cx), float32(vertices[i].cy),
					float32(vertices[i].x), float32(vertices[i].y),
					objspace_flatness_squared, 0)
				x = float32(vertices[i].x)
				y = float32(vertices[i].y)
			case vcubic:
				tesselateCubic(points, &num_points, x, y,
					float32(vertices[i].cx), float32(vertices[i].cy),
					float32(vertices[i].cx1), float32(vertices[i].cy1),
					float32(vertices[i].x), float32(vertices[i].y),
					objspace_flatness_squared, 0)
				x = float32(vertices[i].x)
				y = float32(vertices[i].y)
			}
		}
	}
	return points, lengths
}
