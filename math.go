package imgui

import (
	"math"
	"strconv"
)

type ImVec2 struct {
	x, y float
}

func NewImVec2(x, y float) *ImVec2 {
	return &ImVec2{x, y}
}

// ImVec4: 4D vector used to store clipping rectangles, colors etc. [Compile-time configurable type]
type ImVec4 struct {
	x, y, z, w float
}

func (v *ImVec4) X() float {
	return v.x
}

func (v *ImVec4) Y() float {
	return v.y
}

func (v *ImVec4) Z() float {
	return v.z
}

func (v *ImVec4) W() float {
	return v.w
}

func ImFabs(X float) float {
	return float(math.Abs(float64(X)))
}

func ImSqrt(X float) float {
	return float(math.Sqrt(float64(X)))
}

func ImFmod(X, Y float) float {
	return float(math.Mod(float64(X), float64(Y)))
}

func ImCos(X float) float {
	return float(math.Cos(float64(X)))
}

func ImSin(X float) float {
	return float(math.Sin(float64(X)))
}

func ImAcos(X float) float {
	return float(math.Acos(float64(X)))
}

func ImAtan2(Y, X float) float {
	return float(math.Atan2(float64(Y), float64(X)))
}

func ImAtof(str string) float {
	f, _ := strconv.ParseFloat(str, 64)
	return float(f)
}

func ImCeil(X float) float {
	return float(math.Ceil(float64(X)))
}

func ImPow(X, Y float) float {
	return float(math.Pow(float64(X), float64(Y)))
}

func ImLog(X float) float {
	return float(math.Log(float64(X)))
}

func ImAbs(X float) float {
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

func ImRsqrt(X float) float {
	return 1.0 / float(math.Sqrt(float64(X)))
}

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

func ImLengthSqrVec2(a ImVec2) float {
	return a.x*a.x + a.y*a.y
}

func ImLengthSqrVec4(a ImVec4) float {
	return a.x*a.x + a.y*a.y + a.z*a.z + a.w*a.w
}

func ImInvLength(lhs ImVec2, fail_value float) float {
	var d float = (lhs.x * lhs.x) + (lhs.y * lhs.y)
	if d > 0.0 {
		return ImRsqrt(d)
	}
	return fail_value
}

func ImFloor(f float) float {
	return (float)((int)(f))
}

func ImFloorSigned(f float) float {
	if f >= 0 || float((int)(f)) == f {
		return (float)((int)(f))
	}
	return (float)((int)(f)) - 1
}

func ImFloorVec(v *ImVec2) *ImVec2 {
	return &ImVec2{(float)((int)(v.x)), (float)((int)(v.y))}
}

func ImModPositive(a, b int) int {
	return (a + b) % b
}

func ImDot(a, b *ImVec2) float {
	return a.x*b.x + a.y*b.y
}

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

func ImMul(lhs, rhs *ImVec2) *ImVec2 {
	return &ImVec2{lhs.x * rhs.x, lhs.y * rhs.y}
}

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

func ImLerp(a, b, t float) float {
	return a + (b-a)*t
}

func ImSwap(a, b float) {
	a, b = b, a
}

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
	panic("not implemented")
}

func ImBezierCubicClosestPoint(p1, p2, p3, p4 *ImVec2, p *ImVec2, num_segments int) ImVec2 {
	panic("not implemented")
}

func ImBezierCubicClosestPointCasteljau(p1, p2, p3, p4 *ImVec2, p *ImVec2, tess_tol float32) ImVec2 {
	panic("not implemented")
}

func ImBezierQuadraticCalc(p1, p2, p3 *ImVec2, t float32) ImVec2 {
	panic("not implemented")
}

func ImLineClosestPoint(a, b, p *ImVec2) ImVec2 {
	panic("not implemented")
}

func ImTriangleContainsPoint(a, b, c, p *ImVec2) bool {
	panic("not implemented")
}

func ImTriangleBarycentricCoords(a, b, c, p *ImVec2, out_u, out_v, out_w *float32) {
	panic("not implemented")
}

func ImTriangleArea(a, b, c *ImVec2) float32 {
	return ImFabs((a.x*(b.y-c.y))+(b.x*(c.y-a.y))+(c.x*(a.y-b.y))) * 0.5
}

func ImGetDirQuadrantFromDelta(dx, dy float32) ImGuiDir {
	panic("not implemented")
}

type ImVec1 struct {
	X float
}

type ImVec2ih struct {
	x int16
	y int16
}

type ImRect struct {
	Min ImVec2
	Max ImVec2
}

func (this *ImRect) GetCenter() ImVec2 {
	return ImVec2{(this.Min.x + this.Max.x) * 0.5, (this.Min.y + this.Max.y) * 0.5}
}

func (this *ImRect) GetSize() ImVec2 {
	return ImVec2{this.Max.x - this.Min.x, this.Max.y - this.Min.y}
}

func (this *ImRect) GetWidth() float {
	return this.Max.x - this.Min.x
}

func (this *ImRect) GetHeight() float {
	return this.Max.y - this.Min.y
}

func (this *ImRect) GetArea() float {
	return (this.Max.x - this.Min.x) * (this.Max.y - this.Min.y)
}

func (this *ImRect) GetTL() ImVec2 {
	return this.Min
}

func (this *ImRect) GetTR() ImVec2 {
	return ImVec2{this.Max.x, this.Min.y}
}

func (this *ImRect) GetBL() ImVec2 {
	return ImVec2{this.Min.x, this.Max.y}
}

func (this *ImRect) GetBR() ImVec2 {
	return this.Max
}

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

// ImDrawList: Helper function to calculate a circle's segment count given its radius and a "maximum error" value.
// Estimation of number of circle segment based on error is derived using method described in https://stackoverflow.com/a/2244088/15194693
// Number of segments (N) is calculated using equation:
//   N = ceil ( pi / acos(1 - error / r) )     where r > 0, error <= r
// Our equation is significantly simpler that one in the post thanks for choosing segment that is
// perpendicular to X axis. Follow steps in the article from this starting condition and you will
// will get this result.
//
// Rendering circles with an odd number of segments, while mathematically correct will produce
// asymmetrical results on the raster grid. Therefore we're rounding N to next even number (7->8, 8->8, 9->10 etc.)
//
func IM_ROUNDUP_TO_EVEN(V float) float {
	return ((((V) + 1) / 2) * 2)
}

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
