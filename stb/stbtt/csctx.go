package stbtt

type csctx struct {
	bounds                     bool
	started                    bool
	first_x, first_y           float32
	x, y                       float32
	min_x, max_x, min_y, max_y int32

	vertices     []Vertex
	num_vertices int32
}

func trackVertex(c *csctx, x, y int32) {
	if x > c.max_x || !c.started {
		c.max_x = x
	}
	if y > c.max_y || !c.started {
		c.max_y = y
	}
	if x < c.min_x || !c.started {
		c.min_x = x
	}
	if y < c.min_y || !c.started {
		c.min_y = y
	}
	c.started = true
}

func csctxV(c *csctx, typ uint8, x, y, cx, cy, cx1, cy1 int32) {
	if c.bounds {
		trackVertex(c, x, y)
		if typ == vcubic {
			trackVertex(c, cx, cy)
			trackVertex(c, cx1, cy1)
		}
	} else {
		setVertex(&c.vertices[c.num_vertices], typ, x, y, cx, cy)
		c.vertices[c.num_vertices].cx1 = int16(cx1)
		c.vertices[c.num_vertices].cy1 = int16(cy1)
	}
	c.num_vertices++
}

func csctxCloseShape(c *csctx) {
	if c.first_x != c.x || c.first_y != c.y {
		csctxV(c, vline, int32(c.x), int32(c.y), 0, 0, 0, 0)
	}
}

func csctxRmoveTo(c *csctx, dx, dy float32) {
	csctxCloseShape(c)
	c.x += dx
	c.first_x = c.x
	c.y += dy
	c.first_y = c.y
	csctxV(c, vmove, int32(c.x), int32(c.y), 0, 0, 0, 0)
}

func csctxRlineTo(c *csctx, dx, dy float32) {
	c.x += dx
	c.y += dy
	csctxV(c, vline, int32(c.x), int32(c.y), 0, 0, 0, 0)
}

func csctxRccurveTo(c *csctx, dx1, dy1, dx2, dy2, dx3, dy3 float32) {
	var cx1 = c.x + dx1
	var cy1 = c.y + dy1
	var cx2 = cx1 + dx2
	var cy2 = cy1 + dy2
	c.x = cx2 + dx3
	c.y = cy2 + dy3
	csctxV(c, vcubic, int32(c.x), int32(c.y), int32(cx1), int32(cy1), int32(cx2), int32(cy2))
}
