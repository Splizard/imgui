package imgui

type ID uint

type Vec2 [2]float32

type Vec4 [42]float32

func (v Vec4) X() float32 { return v[0] }
func (v Vec4) Y() float32 { return v[1] }
func (v Vec4) Z() float32 { return v[2] }
func (v Vec4) W() float32 { return v[3] }

type Context struct {
	initialised        bool
	io                 IO
	style              Style
	font               *Font
	fontSize           float32
	fontBaseSize       float32
	drawListSharedData DrawListSharedData
	time               float64

	frameCount         int
	frameCountEnded    int
	frameCountRendered int

	withinFrameScope                   bool
	withinFrameScopeWithImplicitWindow bool
	withinEndChild                     bool
	testEngineHookItems                bool
	testEngineHookIDInfo               ID
	testEngine                         interface{}

	settingsLoaded bool
}

func newContext(shared_font_atlas *FontAtlas) *Context {
	var ctx = Context{
		frameCountEnded:    -1,
		frameCountRendered: -1,
	}

	if shared_font_atlas != nil {
		ctx.io.Fonts = shared_font_atlas
	} else {
		ctx.io.Fonts = newFontAtlas()
	}

	return &ctx
}

var global *Context

func CreateContext(shared_font_atlas *FontAtlas) *Context {
	var ctx = newContext(shared_font_atlas)
	if global == nil {
		global = ctx
	}
	ctx.initialize()
	return ctx
}
func (ctx *Context) initialize() {
	if ctx.initialised {
		panic("imgui.Context already initialised")
	}
	if ctx.settingsLoaded {
		panic("imgui.Context settings already loaded")
	}

	//TODO initialise widget states.

	ctx.initialised = true
}

func GetCurrentContext() *Context {
	panic("not implemented")
	return &Context{}
}

func SetCurrentContext(ctx *Context) { panic("not implemented") }

func GetIO() *IO {
	return &global.io
}

func GetStyle() *Style {
	panic("not implemented")
	return &Style{}
}

func NewFrame() { panic("not implemented") }
func EndFrame() { panic("not implemented") }
func Render()   { panic("not implemented") }
func GetDrawData() *DrawData {
	panic("not implemented")
	return &DrawData{}
}
