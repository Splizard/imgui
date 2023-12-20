package imgui

import (
	"math"
	"strconv"
)

type ImVec2 struct{ x, y float }

func NewImVec2(x, y float) *ImVec2 { return &ImVec2{x, y} }
func (v ImVec2) X() float          { return v.x }
func (v ImVec2) Y() float          { return v.y }
func (v ImVec2) Axis(axis ImGuiAxis) float {
	if axis == ImGuiAxis_X {
		return v.x
	}
	return v.y
}

func (v ImVec2) Add(b ImVec2) ImVec2  { return ImVec2{v.x + b.x, v.y + b.y} }
func (v ImVec2) Sub(b ImVec2) ImVec2  { return ImVec2{v.x - b.x, v.y - b.y} }
func (v ImVec2) Mul(b ImVec2) ImVec2  { return ImVec2{v.x * b.x, v.y * b.y} }
func (v ImVec2) Div(b ImVec2) ImVec2  { return ImVec2{v.x / b.x, v.y / b.y} }
func (v ImVec2) Scale(f float) ImVec2 { return ImVec2{v.x * f, v.y * f} }

// ImVec4 4D vector used to store clipping rectangles, colors etc. [Compile-time configurable type]
type ImVec4 struct{ x, y, z, w float }

func NewImVec4(x, y, z, w float) *ImVec4 { return &ImVec4{x, y, z, w} }
func (v *ImVec4) X() float               { return v.x }
func (v *ImVec4) Y() float               { return v.y }
func (v *ImVec4) Z() float               { return v.z }
func (v *ImVec4) W() float               { return v.w }
func ImFabs(X float) float               { return float(math.Abs(float64(X))) }
func ImSqrt(X float) float               { return float(math.Sqrt(float64(X))) }
func ImFmod(X, Y float) float            { return float(math.Mod(float64(X), float64(Y))) }
func ImCos(X float) float                { return float(math.Cos(float64(X))) }
func ImSin(X float) float                { return float(math.Sin(float64(X))) }
func ImAcos(X float) float               { return float(math.Acos(float64(X))) }
func ImAtan2(Y, X float) float           { return float(math.Atan2(float64(Y), float64(X))) }
func ImAtof(str string) float {
	f, _ := strconv.ParseFloat(str, 64)
	return float(f)
}
func ImCeil(X float) float   { return float(math.Ceil(float64(X))) }
func ImPow(X, Y float) float { return float(math.Pow(float64(X), float64(Y))) }
func ImLog(X float) float    { return float(math.Log(float64(X))) }
func ImAbs(X float) float {
	if X < 0 {
		return -X
	}
	return X
}

func ImAbsInt(X int) int {
	if X < 0 {
		return -X
	}
	return X
}

func ImSign(X float) float {
	if X < 0 {
		return -1
	}
	if X > 0 {
		return 1
	}
	return 0
}

func ImRsqrt(X float) float { return 1.0 / float(math.Sqrt(float64(X))) }

func ImMinVec2(lhs, rhs *ImVec2) (res ImVec2) {
	if lhs.x < rhs.x {
		res.x = lhs.x
	} else {
		res.x = rhs.x
	}
	if lhs.y < rhs.y {
		res.y = lhs.y
	} else {
		res.y = rhs.y
	}
	return
}

func ImMaxVec2(lhs, rhs *ImVec2) (res ImVec2) {
	if lhs.x > rhs.x {
		res.x = lhs.x
	} else {
		res.x = rhs.x
	}
	if lhs.y > rhs.y {
		res.y = lhs.y
	} else {
		res.y = rhs.y
	}
	return

}

func ImClampVec2(v, mn *ImVec2, mx ImVec2) (res ImVec2) {
	if v.x < mn.x {
		res.x = mn.x
	} else if v.x > mx.x {
		res.x = mx.x
	} else {
		res.x = v.x
	}
	if v.y < mn.y {
		res.y = mn.y
	} else if v.y > mx.y {
		res.y = mx.y
	} else {
		res.y = v.y
	}
	return
}

func ImLerpVec2(a, b *ImVec2, t float) ImVec2 {
	return ImVec2{a.x + (b.x-a.x)*t, a.y + (b.y-a.y)*t}
}

func ImLerpVec2WithVec2(a, b *ImVec2, t ImVec2) ImVec2 {
	return ImVec2{a.x + (b.x-a.x)*t.x, a.y + (b.y-a.y)*t.y}
}

func ImLerpVec4(a, b *ImVec4, t float) ImVec4 {
	return ImVec4{a.x + (b.x-a.x)*t, a.y + (b.y-a.y)*t, a.z + (b.z-a.z)*t, a.w + (b.w-a.w)*t}
}

func ImSaturate(f float) float {
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 1
	}
	return f
}

func ImLengthSqrVec2(a ImVec2) float { return a.x*a.x + a.y*a.y }
func ImLengthSqrVec4(a ImVec4) float { return a.x*a.x + a.y*a.y + a.z*a.z + a.w*a.w }

func ImInvLength(lhs ImVec2, fail_value float) float {
	var d = (lhs.x * lhs.x) + (lhs.y * lhs.y)
	if d > 0.0 {
		return ImRsqrt(d)
	}
	return fail_value
}

func ImFloor(f float) float { return (float)((int)(f)) }

func ImFloorSigned(f float) float {
	if f >= 0 || float((int)(f)) == f {
		return (float)((int)(f))
	}
	return (float)((int)(f)) - 1
}

func ImFloorVec(v *ImVec2) *ImVec2 { return &ImVec2{(float)((int)(v.x)), (float)((int)(v.y))} }
func ImModPositive(a, b int) int   { return (a + b) % b }
func ImDot(a, b *ImVec2) float     { return a.x*b.x + a.y*b.y }
func ImRotate(v *ImVec2, cos_a, sin_a float) *ImVec2 {
	return &ImVec2{v.x*cos_a - v.y*sin_a, v.x*sin_a + v.y*cos_a}
}

func ImLinearSweep(current, target, speed float) float {
	if current < target {
		return ImMin(current+speed, target)
	}
	if current > target {
		return ImMax(current-speed, target)
	}
	return current
}

func ImMul(lhs, rhs *ImVec2) *ImVec2 { return &ImVec2{lhs.x * rhs.x, lhs.y * rhs.y} }

func ImMin(a, b float) float {
	if a < b {
		return a
	}
	return b
}

func ImMinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ImMaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ImMax(a, b float) float {
	if a > b {
		return a
	}
	return b
}

func ImClamp(v, mn, mx float) float {
	if v < mn {
		return mn
	}
	if v > mx {
		return mx
	}
	return v
}

func ImClamp64(v, mn, mx float64) float64 {
	if v < mn {
		return mn
	}
	if v > mx {
		return mx
	}
	return v
}

func ImClampInt(v, mn, mx int) int {
	if v < mn {
		return mn
	}
	if v > mx {
		return mx
	}
	return v
}

func ImClampInt64(v, mn, mx int64) int64 {
	if v < mn {
		return mn
	}
	if v > mx {
		return mx
	}
	return v
}

func ImClampUint64(v, mn, mx uint64) uint64 {
	if v < mn {
		return mn
	}
	if v > mx {
		return mx
	}
	return v
}

func ImLerp(a, b, t float) float { return a + (b-a)*t }
func ImSwap(a, b float)          { a, b = b, a }

func ImAddClampOverflow(a, b, mn, mx float) float {
	if b < 0 && (a < mn-b) {
		return mn
	}
	if b > 0 && (a > mx-b) {
		return mx
	}
	return a + b
}

func ImSubClampOverflow(a, b, mn, mx float) float {
	if b > 0 && (a < mn+b) {
		return mn
	}
	if b < 0 && (a > mx+b) {
		return mx
	}
	return a - b
}

func ImBezierCubicCalc(p1, p2, p3, p4 *ImVec2, t float32) ImVec2 {
	var u = 1.0 - t
	var w1 = u * u * u
	var w2 = 3 * u * u * t
	var w3 = 3 * u * t * t
	var w4 = t * t * t
	return ImVec2{w1*p1.x + w2*p2.x + w3*p3.x + w4*p4.x, w1*p1.y + w2*p2.y + w3*p3.y + w4*p4.y}
}

func ImBezierCubicClosestPoint(p1, p2, p3, p4 *ImVec2, p *ImVec2, num_segments int) ImVec2 {
	IM_ASSERT(num_segments > 0) // Use ImBezierCubicClosestPointCasteljau()
	var p_last = *p1
	var p_closest ImVec2
	var p_closest_dist2 float = FLT_MAX
	var t_step = 1.0 / (float)(num_segments)
	for i_step := 1; int(i_step) <= num_segments; i_step++ {
		var p_current = ImBezierCubicCalc(p1, p2, p3, p4, t_step*float(i_step))
		var p_line = ImLineClosestPoint(&p_last, &p_current, p)
		var dist2 = ImLengthSqrVec2(p.Sub(p_line))
		if dist2 < p_closest_dist2 {
			p_closest = p_line
			p_closest_dist2 = dist2
		}
		p_last = p_current
	}
	return p_closest
}

// ImBezierCubicClosestPointCasteljau tess_tol is generally the same value you would find in ImGui::GetStyle().CurveTessellationTol
// Because those ImXXX functions are lower-level than ImGui:: we cannot access this value automatically.
func ImBezierCubicClosestPointCasteljau(p1, p2, p3, p4 *ImVec2, p *ImVec2, tess_tol float32) ImVec2 {
	IM_ASSERT(tess_tol > 0.0)
	var p_last = *p1
	var p_closest ImVec2
	var p_closest_dist2 float = FLT_MAX
	ImBezierCubicClosestPointCasteljauStep(p, &p_closest, &p_last, p_closest_dist2, p1.x, p1.y, p2.x, p2.y, p3.x, p3.y, p4.x, p4.y, tess_tol, 0)
	return p_closest
}

// ImBezierCubicClosestPointCasteljauStep Closely mimics PathBezierToCasteljau
func ImBezierCubicClosestPointCasteljauStep(p, p_closest, p_last *ImVec2, p_closest_dist2, x1, y1, x2, y2, x3, y3, x4, y4, tess_tol float, level int) {
	var dx = x4 - x1
	var dy = y4 - y1
	var d2 = (x2-x4)*dy - (y2-y4)*dx
	var d3 = (x3-x4)*dy - (y3-y4)*dx
	if d2 < 0 {
		d2 = -d2
	}
	if d3 < 0 {
		d3 = -d3
	}
	if (d2+d3)*(d2+d3) < tess_tol*(dx*dx+dy*dy) {
		var p_current = ImVec2{x4, y4}
		var p_line = ImLineClosestPoint(p_last, &p_current, p)
		var dist2 = ImLengthSqrVec2(p.Sub(p_line))
		if dist2 < p_closest_dist2 {
			*p_closest = p_line
			p_closest_dist2 = dist2
		}
		*p_last = p_current
	} else if level < 10 {
		var x12 = (x1 + x2) * 0.5
		var y12 = (y1 + y2) * 0.5
		var x23 = (x2 + x3) * 0.5
		var y23 = (y2 + y3) * 0.5
		var x34 = (x3 + x4) * 0.5
		var y34 = (y3 + y4) * 0.5
		var x123 = (x12 + x23) * 0.5
		var y123 = (y12 + y23) * 0.5
		var x234 = (x23 + x34) * 0.5
		var y234 = (y23 + y34) * 0.5
		var x1234 = (x123 + x234) * 0.5
		var y1234 = (y123 + y234) * 0.5
		ImBezierCubicClosestPointCasteljauStep(p, p_closest, p_last, p_closest_dist2, x1, y1, x12, y12, x123, y123, x1234, y1234, tess_tol, level+1)
		ImBezierCubicClosestPointCasteljauStep(p, p_closest, p_last, p_closest_dist2, x1234, y1234, x234, y234, x34, y34, x4, y4, tess_tol, level+1)
	}
}

func ImBezierQuadraticCalc(p1, p2, p3 *ImVec2, t float32) ImVec2 {
	u := 1.0 - t
	w1 := u * u
	w2 := 2 * u * t
	w3 := t * t
	return ImVec2{w1*p1.x + w2*p2.x + w3*p3.x, w1*p1.y + w2*p2.y + w3*p3.y}
}

func ImLineClosestPoint(a, b, p *ImVec2) ImVec2 {
	var ap = p.Sub(*a)
	var ab_dir = b.Sub(*a)
	var dot = ap.x*ab_dir.x + ap.y*ab_dir.y
	if dot < 0.0 {
		return *a
	}
	var ab_len_sqr = ab_dir.x*ab_dir.x + ab_dir.y*ab_dir.y
	if dot > ab_len_sqr {
		return *b
	}
	return a.Add(ab_dir.Scale(dot / ab_len_sqr))
}

func ImTriangleContainsPoint(a, b, c, p *ImVec2) bool {
	var b1 = ((p.x-b.x)*(a.y-b.y) - (p.y-b.y)*(a.x-b.x)) < 0.0
	var b2 = ((p.x-c.x)*(b.y-c.y) - (p.y-c.y)*(b.x-c.x)) < 0.0
	var b3 = ((p.x-a.x)*(c.y-a.y) - (p.y-a.y)*(c.x-a.x)) < 0.0
	return (b1 == b2) && (b2 == b3)
}

func ImTriangleBarycentricCoords(a, b, c, p *ImVec2, out_u, out_v, out_w *float32) {
	var v0 = b.Sub(*a)
	var v1 = c.Sub(*a)
	var v2 = p.Sub(*a)
	var denom = v0.x*v1.y - v1.x*v0.y
	*out_v = (v2.x*v1.y - v1.x*v2.y) / denom
	*out_w = (v0.x*v2.y - v2.x*v0.y) / denom
	*out_u = 1.0 - *out_v - *out_w
}

func ImTriangleArea(a, b, c *ImVec2) float32 {
	return ImFabs((a.x*(b.y-c.y))+(b.x*(c.y-a.y))+(c.x*(a.y-b.y))) * 0.5
}

func ImTriangleClosestPoint(a, b, c, p *ImVec2) ImVec2 {
	var proj_ab = ImLineClosestPoint(a, b, p)
	var proj_bc = ImLineClosestPoint(b, c, p)
	var proj_ca = ImLineClosestPoint(c, a, p)
	var dist2_ab = ImLengthSqrVec2(p.Sub(proj_ab))
	var dist2_bc = ImLengthSqrVec2(p.Sub(proj_bc))
	var dist2_ca = ImLengthSqrVec2(p.Sub(proj_ca))
	var m = ImMin(dist2_ab, ImMin(dist2_bc, dist2_ca))
	if m == dist2_ab {
		return proj_ab
	}
	if m == dist2_bc {
		return proj_bc
	}
	return proj_ca
}

func ImGetDirQuadrantFromDelta(dx, dy float32) ImGuiDir {
	if ImFabs(dx) > ImFabs(dy) {
		if dx > 0 {
			return ImGuiDir_Right
		}
		return ImGuiDir_Left
	}
	if dy > 0 {
		return ImGuiDir_Down
	}
	return ImGuiDir_Up
}

type ImVec1 struct{ x float }

type ImVec2ih struct {
	x int16
	y int16
}

type ImRect struct {
	Min ImVec2
	Max ImVec2
}

func ImRectFromVec4(v *ImVec4) ImRect { return ImRect{ImVec2{v.x, v.y}, ImVec2{v.z, v.w}} }
func (this *ImRect) GetCenter() ImVec2 {
	return ImVec2{(this.Min.x + this.Max.x) * 0.5, (this.Min.y + this.Max.y) * 0.5}
}
func (this *ImRect) GetSize() ImVec2  { return ImVec2{this.Max.x - this.Min.x, this.Max.y - this.Min.y} }
func (this *ImRect) GetWidth() float  { return this.Max.x - this.Min.x }
func (this *ImRect) GetHeight() float { return this.Max.y - this.Min.y }
func (this *ImRect) GetArea() float   { return (this.Max.x - this.Min.x) * (this.Max.y - this.Min.y) }
func (this *ImRect) GetTL() ImVec2    { return this.Min }
func (this *ImRect) GetTR() ImVec2    { return ImVec2{this.Max.x, this.Min.y} }
func (this *ImRect) GetBL() ImVec2    { return ImVec2{this.Min.x, this.Max.y} }
func (this *ImRect) GetBR() ImVec2    { return this.Max }
func (this *ImRect) ContainsVec(p ImVec2) bool {
	return p.x >= this.Min.x && p.y >= this.Min.y && p.x < this.Max.x && p.y < this.Max.y
}
func (this *ImRect) ContainsRect(r ImRect) bool {
	return r.Min.x >= this.Min.x && r.Min.y >= this.Min.y && r.Max.x <= this.Max.x && r.Max.y <= this.Max.y
}
func (this *ImRect) Overlaps(r ImRect) bool {
	return r.Min.y < this.Max.y && r.Max.y > this.Min.y && r.Min.x < this.Max.x && r.Max.x > this.Min.x
}
func (this *ImRect) AddVec(p ImVec2) {
	if this.Min.x > p.x {
		this.Min.x = p.x
	}
	if this.Min.y > p.y {
		this.Min.y = p.y
	}
	if this.Max.x < p.x {
		this.Max.x = p.x
	}
	if this.Max.y < p.y {
		this.Max.y = p.y
	}
}

func (this *ImRect) AddRect(r ImRect) {
	if this.Min.x > r.Min.x {
		this.Min.x = r.Min.x
	}
	if this.Min.y > r.Min.y {
		this.Min.y = r.Min.y
	}
	if this.Max.x < r.Max.x {
		this.Max.x = r.Max.x
	}
	if this.Max.y < r.Max.y {
		this.Max.y = r.Max.y
	}
}

func (this *ImRect) Expand(amount float) {
	this.Min.x -= amount
	this.Min.y -= amount
	this.Max.x += amount
	this.Max.y += amount
}

func (this *ImRect) ExpandVec(amount ImVec2) {
	this.Min.x -= amount.x
	this.Min.y -= amount.y
	this.Max.x += amount.x
	this.Max.y += amount.y
}

func (this *ImRect) Translate(d ImVec2) {
	this.Min.x += d.x
	this.Min.y += d.y
	this.Max.x += d.x
	this.Max.y += d.y
}

func (this *ImRect) TranslateX(dx float) {
	this.Min.x += dx
	this.Max.x += dx
}

func (this *ImRect) TranslateY(dy float) {
	this.Min.y += dy
	this.Max.y += dy
}

func (this *ImRect) ClipWith(r ImRect) {
	this.Min = ImMaxVec2(&this.Min, &r.Min)
	this.Max = ImMinVec2(&this.Max, &r.Max)
}

func (this *ImRect) ClipWithFull(r ImRect) {
	this.Min = ImClampVec2(&this.Min, &r.Min, r.Max)
	this.Max = ImClampVec2(&this.Max, &r.Min, r.Max)
}

func (this *ImRect) Floor() {
	this.Min.x = IM_FLOOR(this.Min.x)
	this.Min.y = IM_FLOOR(this.Min.y)
	this.Max.x = IM_FLOOR(this.Max.x)
	this.Max.y = IM_FLOOR(this.Max.y)
}

func (this *ImRect) IsInverted() bool {
	return this.Min.x > this.Max.x || this.Min.y > this.Max.y
}

func (this *ImRect) ToVec4() ImVec4 {
	return ImVec4{this.Min.x, this.Min.y, this.Max.x, this.Max.y}
}

// IM_ROUNDUP_TO_EVEN ImDrawList: Helper function to calculate a circle's segment count given its radius and a "maximum error" value.
// Estimation of number of circle segment based on error is derived using method described in https://stackoverflow.com/a/2244088/15194693
// Number of segments (N) is calculated using equation:
//
//	N = ceil ( pi / acos(1 - error / r) )     where r > 0, error <= r
//
// Our equation is significantly simpler that one in the post thanks for choosing segment that is
// perpendicular to X axis. Follow steps in the article from this starting condition and you will
// will get this result.
//
// Rendering circles with an odd number of segments, while mathematically correct will produce
// asymmetrical results on the raster grid. Therefore we're rounding N to next even number (7->8, 8->8, 9->10 etc.)
func IM_ROUNDUP_TO_EVEN(V float) float { return (((V) + 1) / 2) * 2 }

const IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MIN = 4
const IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MAX = 512

func IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC(RAD, MAXERROR float) float {
	return ImClamp(IM_ROUNDUP_TO_EVEN(ImCeil(IM_PI/ImAcos(1-ImMin(MAXERROR, RAD)/RAD))), IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MIN, IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MAX)
}

func IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC_R(N, MAXERROR float) float {
	return MAXERROR / (1 - ImCos(IM_PI/ImMax(N, IM_PI)))
}
func IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC_ERROR(N, RAD float) float {
	return (1 - ImCos(IM_PI/ImMax(N, IM_PI))) / RAD
}

const IM_DRAWLIST_ARCFAST_TABLE_SIZE = 48
const IM_DRAWLIST_ARCFAST_SAMPLE_MAX = IM_DRAWLIST_ARCFAST_TABLE_SIZE

func IM_NORMALIZE2F_OVER_ZERO(VX, VY *float) {
	var d2 = *VX**VX + *VY**VY
	if d2 > 0.0 {
		var inv_len = ImRsqrt(d2)
		*VX *= inv_len
		*VY *= inv_len
	}
}

const IM_FIXNORMAL2F_MAX_INVLEN2 float = 100

func IM_FIXNORMAL2F(VX, VY *float) {
	var d2 = *VX**VX + *VY**VY
	if d2 > 0.000001 {
		var inv_len2 = 1.0 / d2
		if inv_len2 > IM_FIXNORMAL2F_MAX_INVLEN2 {
			inv_len2 = IM_FIXNORMAL2F_MAX_INVLEN2
		}
		*VX *= inv_len2
		*VY *= inv_len2
	}
}

func ImAcos01(x float) float {
	if x <= 0.0 {
		return IM_PI * 0.5
	}
	if x >= 1.0 {
		return 0.0
	}
	return ImAcos(x)
}

// PathBezierCubicCurveToCasteljau Closely mimics ImBezierCubicClosestPointCasteljau() in imgui.cpp
func PathBezierCubicCurveToCasteljau(path *[]ImVec2, x1, y1, x2, y2, x3, y3, x4, y4, tess_tol float, level int) {
	var dx = x4 - x1
	var dy = y4 - y1
	var d2 = (x2-x4)*dy - (y2-y4)*dx
	var d3 = (x3-x4)*dy - (y3-y4)*dx
	if d2 < 0.0 {
		d2 = -d2
	}
	if d3 < 0.0 {
		d3 = -d3
	}
	if (d2+d3)*(d2+d3) < tess_tol*(dx*dx+dy*dy) {
		*path = append(*path, ImVec2{x4, y4})
	} else if level < 10 {
		var x12, y12 = (x1 + x2) * 0.5, (y1 + y2) * 0.5
		var x23, y23 = (x2 + x3) * 0.5, (y2 + y3) * 0.5
		var x34, y34 = (x3 + x4) * 0.5, (y3 + y4) * 0.5
		var x123, y123 = (x12 + x23) * 0.5, (y12 + y23) * 0.5
		var x234, y234 = (x23 + x34) * 0.5, (y23 + y34) * 0.5
		var x1234, y1234 = (x123 + x234) * 0., (y123 + y234) * 0.5
		PathBezierCubicCurveToCasteljau(path, x1, y1, x12, y12, x123, y123, x1234, y1234, tess_tol, level+1)
		PathBezierCubicCurveToCasteljau(path, x1234, y1234, x234, y234, x34, y34, x4, y4, tess_tol, level+1)
	}
}

func PathBezierQuadraticCurveToCasteljau(path *[]ImVec2, x1, y1, x2, y2, x3, y3, tess_tol float, level int) {
	var dx, dy = x3 - x1, y3 - y1
	var det = (x2-x3)*dy - (y2-y3)*dx
	if det*det*4.0 < tess_tol*(dx*dx+dy*dy) {
		*path = append(*path, ImVec2{x3, y3})
	} else if level < 10 {
		var x12, y12 = (x1 + x2) * 0.5, (y1 + y2) * 0.5
		var x23, y23 = (x2 + x3) * 0.5, (y2 + y3) * 0.5
		var x123, y123 = (x12 + x23) * 0.5, (y12 + y23) * 0.5
		PathBezierQuadraticCurveToCasteljau(path, x1, y1, x12, y12, x123, y123, tess_tol, level+1)
		PathBezierQuadraticCurveToCasteljau(path, x123, y123, x23, y23, x3, y3, tess_tol, level+1)
	}
}
