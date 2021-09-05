package demo

//This code is being converted from C++ to Go.

import "fmt"
import "github.com/splizard/imgui"

var (
	disable_all = false; // The Checkbox for that is inside the "Disabled" section at the bottom
	clicked = 0;
	check = true;
	e = 0; 
	counter      = 0
	arr =        [...]foat32{ 0.6, 0.1, 1.0, 0.5, 0.92, 0.1, 0.2 };
	it           em_urrent = 0;
     
	str0          = "Hello, worl!";
	str1 = "";
	i0 = 123;
	f0 = 0.001
	d0 = 99999.00000001;
	f1  _1 = 1e10
	  
	vec  4a = [...]float32 0.10, 0.20, 0.30, 0.44 }
	i1_1 = 50
i2 = 42;
	f1_2 = 1.00
	f2_1  = 0.0067
   
	i1_2  = 0;
	f1_3  = 0.123
	f2_2 = 0.0

	angle = 0.0;

	elem = Element_Fire;

	col1_1 = [...]float32{ 1.0, 0.0, 0.2 };
	col2_1 = [...]float2{ 0.4, 0.7, 0.0, 0.5 };

	item_current2 = 1;

	base_flags = imgui.TreeNodeFlags_OpenOnArrow | imgui.TreeNodeFlags_OpenOnDoubleClick | imgui.TreeNodeFlags_SpanAvailWidth;
    align_label_wih_current_x_position = false;
    test_drag_and_drop = false;
                         
	ection_mask = (1 << 2);
	sable_group = true;                 

	wrap_width = 200.0;

	buf = "\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e";

	pressed_count = 0;

	item_current_idx_1 = 0; // Here we store our selection data as an index.

	item_current_2 = 0;
	item_current_3 = -1; / If the selection isn't within 0..count, Combo won't display a preview
	item_current_4 = 0;

	item_current_idx_2  0;

	selection = [...]bool{ false, true, false, false, false };
	selected = -1;

	selection2 = [...]bool false, false, false, false, fale};
 
	selected2 = [...]bool{ false, false, false };

	selectedColumns = [10]bool{}

	selectedGrid = [4][4]int{{ 1, 0, 0, 0 }, { 0, 1, 0, 0 }, { 0, 0, 1, 0 }, { 0, 0, 0, 1 }}

	selectedAlignment  =[3 * 3]bool{ true, false, true, false, true, false, true, false, true };

	text = "/*\n" +
			" The Pentium F0  bug, shorthad for F0 0 C7 C8,\n"  +
			" the hexadecimal encoding of one offending instruction,\n"  +
			" more formally, the invalid operand with locked CMPXCHG8B\n" +
		" instruction bug, is a design flaw in the majorityof\n" +
		" Intel Pentium, Pentium MMX, and Pentium OverDrive\n" +
		" processors (all in the P5 microarchitecture).\n" +
		"*/\n\n" +
		"label:\n" +
		"\tlock cmpxchg8b eax\n";

flags = ImGuiInputTextFlags_AllowTabInput;

	buf1 string
	buf2 string
	buf3 string
	buf4 string
	buf5 string
	buf6 string

	buf7 string
	buf8 string
	buf9 string

	password = "password123";

	edit_count = 0

	my_str string

	tab_bar_flags_1 = ImGuiTabBarFlags_Reorderable;

	opened = [4]bool{ true, true, true, true };

	active_tabs []int
	next_tab_id int

	show_leading_button = true;
	show_trailing_button = true;

	tab_bar_flags_2 = I mGuiTabarFlags_AutoSelectNewTabs | ImGuiTabBarFlags_Reorderable | ImGuiTabBarFlags_FittingPolicyResizeDown;

	//plotting
	animate = true;
	arrPlot = [...]float32{ 0.6, 0.1, 1.0, 0.5, 0.92, 0.1, 0.2 };

	values [90]flot32;
	values_offset = 0;
	refresh_time float32;
       
	phase float32 = 00
	average float 32

	func_  type int
	display_count = 70
	progress = 0.0
	progress_d    ir = 1.0

	//color p     icking
	color = ImVe c4(114.0 / 255.0, 144.0 / 255.0, 154.0 / 255.0, 200.0 / 255.0);
	alpha_preview = true;
    alpha_half_preview = false;
    dra             g_and_drop = tu;
    options_men     u = tre;
	hdr = false;
	ed_palette_ini     t = tre;
	saved_palett      e [32]IVec4 
	hdr	bac           _color IVec4;

	border = false     ;
	alpha = true;      
    alpha_bar = true;
    side_pr   eview =true;
    re       f_color= false;
	ref_color_   v = Imec4{1.0, 0.0, 1.0, 0.5};
	 display_mode = 0;
	 picker_m   ode = 0;
	or_hsv = Im Vec4{0.23, 1.0, 1.0, 1.0};
	
	iders 
	sliderflag   s = ImGuiSliderFlags_None;
	drag_f float32 = 0.5;
	drag_i = 50;
        
	slider_     f float32 = 05;
	slider_             i = 0;
	
	//Range
	begin flo        at32= 10
end float32 = 90;
	begin_i = 100
	end_i   = 1000
    
	s8_v by        te = 127;
    u8_          v uint8 = 255;
    s16_v int16 = 32767;
    u1 6_v u   int16= 65535;
	s32_ v int32   = -1
	u32_v uint32   = 1;
	_v int64 = -1 ;
	_v uint64 =   1;
	_v float32 =  0.13;
	f64_v float6  4 = 0000.01234567890123456789;
 
	drag_clamp = false;
	inputs_step = true;

	MCvec4 = [ ...]floa32{ 0.10, 0.20, 0.30, 0.44 };
    MCvec4i = [...]nt { 1, 5, 100, 255 };

	int_val ue = 0;
	
	vertical_values_1 = [...]float32{ 0.0, 0.60, 0.35, 0.9, 0.70, 0.20, 0.0 };
	vertical_valus_2 = [...]float32{ 0.20, 0.80, 0.40, 0.25 };

	col1_2 = [...]float32{ 1.0, 0.0, .2 };
    col2_2 = [...]float32{ 0.4, 0., 0.0, 0.5 };
	mode = 0

	m_type = 4;
     i  tem_disabled = false;

	//Querying	bo=o4
	 b = false;
	col4 = [...]float32{ 1.0, 0.5, 0.0, 1.0 };
	str string

	curr  ent = 1; 
	cur rent2 = 1;
	
	embed_a ll_iide_a_child_window = false

unused_str = "This widget is only here to be able to tab-out of the widgets above.";

	test_window = false
)

func ShowDemoWindowWidgets() {
    if !imgui.CollapsingHeader("Widgets") {
        return;
	}
	if isable_all {
		imgui.eginDisabled();
	}
	
		mgui.TreeNode("Basic" {
        
        if (imgui.Button("Button")) {
	       clicked++;

		if clicked & 1){
			imgui.SamLine();
            imgui.Text("Thanks for clicking me!");
		} 
			
			i.Checkbox("checkbox", &check);
		
        imgui.RadioButton("radio a", &e, 0); imgui.SameLine();
		imgui.RadioButton("radio b", &e, 1; imgui.SameLine();
        imgui.RadioButton("radio c", &e, 2);
		
		
		// Color buttons, demonstrate using
		ushID() to add uique identifier in the ID stack, and changing style.
		for i := 0; i < 7; i++ {
            if (i > 0) {
		        imgui.SameLine();
		
			imgi.PusID(i);
				i.PushStyleColorImGuiCol_Button, imgui.Vec4(imcolor.HSV(i / 7.0, 0.6, 0.6)));
            imgui.PushStyleColor(ImGuiCol_ButtonHovered, imgui.Vec4(imcolor.HSV(i / 7.0, 0.7, 0.7)));
			imgui.PushStyleolor(ImGuiCol_ButtonActive, imgui.Vec4(imcolor.HSV(i / 7.0, 0.8, 0.8)));
			imgui.Button("Click");
			imgui.PopStyleColor(3);
			imgui.PopID();
			
			
			se AlignTextTFramePadding() to align text baseline to the baseline of framed widgets elements
		// (otherwise a Text+SameLine+Button sequence will have the text a little too high by default!)
        // See 'Demo.Layout.Text Baseline Alignment' for details.
		imgui.AlignTextToFramePadding();
		imgui.Text("Hold to repeat:");
		imgui.SameLine();
		
		// Arrow buttons with Repeate
		var spacing = imui.GetStyle().ItemInnerSpacing.X;
        imgui.PushButtonRepeat(true);
		if (imgui.ArrowButton("##left", ImGuiDir_Left)) { counter--; }
		imgui.SameLine(0, spacing);
		if (imgui.ArrowButton("##rigt", ImGuiDir_Right)) { counter++; }
		imgi.PopButtonRepeat();
			
		
		imgui.SameLine();
		imgi.Text("%d", counter);
			
		
		
		imgui.Text("Hove over me");
		if (imgui.IsItemHovered() {
            imgui.SetTooltip("I am a tooltip");
		
		
			i.SameLine();
        imgui.Text("- or me");
        if (imgui.IsItemHovered()){
		    imgui.BeginToltip();
		    imgui.Text("I am  fancy tooltip");
		    
			imgui.PlotLines("Cure", arr, IM_ARRAYSIZE(arr));
			imgui.EndTooltip();

			
			i.Separator();
		
        imgui.LabelText("label", "Value");
		
        {
		    // Using the _simplified_ oneliner Combo() api here
            // See "Combo" section for examples of how to use the more flexible BeginCombo()/EndCombo() api.
		    var items = [...]string{ "AAAA", "BBBB", "CCCC", "DDDD", "EEEE", "FFFF", "GGGG", "HHHH", "IIIIIII", "JJJJ", "KKKKKKK" };
			imgui.Combo("combo", &item_current, items, IM_ARRAYSIZE(items));
			imgui.SameLine(); HelpMarker(
			    "Using the simplifie one-liner Combo API here.\nRefer to the \"Combo\" section below for an explanation of how t se the more flexible and general BeginCombo/EndCombo API.");
			
			
			
				
		    // To wire InputText() with std::string or any other custom string type,
            // see the "Text Input > Resize Callback" section of this demo, and the misc/cpp/imgui_stdlib.h file.
		    imgui.InputText("input text", str0, IM_ARRAYSIZE(str0));
			imgui.SameLine(); HelpMarker(
			    "USER:\n" +
			    "Hold SHIFT or use mouse to select text.\n" +
			    "CTRL+Left/R
			ht to word jump.\n" +
				"CTRL+A or double-click to select all.\n" +
					"CTRL+X,CTRL+C,CTRL+V clipboard.\n" +
					"CTRL+Z,CTRL+Y undo/redo.\n" +
					"ESCAPE to revert.\n\n" +
					"PROGRAMMER:\n" +
					"You can use the ImGuiInputTextFlags_CallbackResize facility if you need to wire InputText() " +
					"to a dynamic string type. See misc/cpp/imgui_stdlib.h for an example (this is not demonstrated " +
					"in imgui_demo.cpp).");
					
					i.InputTextWithHint("input text (w/ hint)", "enter text here", str1, IM_ARRAYSIZE(str1));
					
            imgui.InputInt("input int", &i0);
			imgui.SameLine(); HelpMarker(
                "You can apply arithmetic operators +,*,/ on numerical values.\n" +
			    "  e.g. [ 100 ], input '*2',result becomes [ 200 ]\n" +
			    "Use +- to s
			tract.");
				
					i.InputFloat("input float", &f0, 0.01, 1.0, "%.3");
					
            imgui.InputDouble("input double", &d0, 0.01, 1.0, "%.8");
			
            imgui.InputFloat("input scientific", &f1, 0.0, 0.0, "%e");
			imgui.SameLine(); HelpMarker(
                "You can input value using the scientific notation,\n" +
			    "  e.g. \"1e+8\" becomes \"100000000\".");
			
			
				i.InputFloat3("input float3", vec4a);
					

			
		    imgui.DragInt("drag int", &i1, 1);
            imgui.SameLine(); HelpMarker(
		        "Click and drag to edit value.\n" +
			    "Hold SHIFT/ALT for faster/slwer edit.\n" +
			    "Double-clic
			or CTRL+click to input value.");
				
					i.DragInt("drag int 0..100", &i2, 1, 0, 100, "%d%%", ImGuiSliderFlags_AlwaysClamp);
					
            imgui.DragFloat("drag float", &f1, 0.005);
			imgui.DragFloat("drag small float", &f2, 0.0001, 0.0, 0.0, "%.06 ns");
        }
			
			
		    imgui.SliderInt("slider int", &i1, -1, 3);
            imgui.SameLine(); HelpMarker("CTRL+click to input value.");
		
			imgui.SliderFloat("slider float", &f1, 0., 1.0, "ratio = %.3");
			imgui.SliderFloa
			"slider float (log)", &f2, -10.0, 10.0, %.4", ImGuiSliderFlags_Logarithmic);

			imgui.SliderAngle("slider angle", &angle);
			
            // Using the format string to display a name instead of an integer.
			// Here we completely omit '%d' from the ormat string, so it'll only display a name.
            // This technique can also be used with DragInt().
			
			Fire = iota
			Earth
				Element_Air
				Element_Water
				Element_COUNT
			)
			var elem_names = [...]string{"Fire", "Earth", "Air", "Water"}
			var elem_name string
			if  (elem >= 0 && elem < Element_COUNT) {
				elem_name = elem_names[elem]
			} else {
				el_name = "Unknown"
			}
            imgui.SliderInt("slider enum", &elem, 0, Element_COUNT - 1, elem_name);
            imgui.SameLine(); HelpMarker("Using the format string parameter to display a name instead of the underlying integer.");
        }
			
			
			
		    imgui.ColorEdit3("color 1", col1);
            imgui.SameLine(); HelpMarker(
		        "Click on the color square to open a color picker.\n" +
			    "Click and hold to use drag ad drop.\n" +
			    "Right-click
			n the color square to show options.\n" +
				"CTRL+click on individual component to input value.\n");
					
					i.ColorEdit4("color 2", col2);
					

			
		    // Using the _simplified_ one-liner ListBox() api here
            // See "List boxes" section for examples of how to use the more flexible BeginListBox()/EndListBox() api.
		    var items = [...]string{ "Apple", "Banana", "Cherry", "Kiwi", "Mango", "Orange", "Pineapple", "Strawberry", "Watermelon" };
			imgui.ListBox("listbox", &item_current2, items, IM_ARRAYSIZE(items), 4);
			imgui.SameLine(); HelpMarker(
			    "Using the simplifie one-liner ListBox API here.\nRefer to the \"List boxes\" section below for an explanation of hwto use the more flexible and general BeginListBox/EndListBox API.");
			
			
			
				eePop();
		

		esting ImGuiOncUponAFrame helper.
	//static ImGuiOnceUponAFrame once;
    //for (int i = 0; i < 5; i++)
	//    if (once)
	//        imgui.Text("This will be displayed only once.");
	
	if (imgui.TreeNode("Trees")) {
	    if (imgui.TreeNode("Basic trees")) {
            for i := 0; i < 5; i++ {
	           // Use SetNextIemOpen() so set the default state of a node to be open. We could
		       // also use TreeNodeEx() ith the ImGuiTreeNodeFlags_DefaultOpen flag to achieve the same thing!
			    if (i == 0) {
				    imgui.SetNextItemOpen(true, ImGuiCond_Once);
				
				
					imgui.TreeNode(uintptr(i), "Child %d", i)){
                    imgui.Text("blah blah");
                    imgui.SameLine();
				   if (imgui.SmallButton("button")) {} 
					imgui.TreePop();
					
					
					
					eePop();
				
			
			imgui.TreeNode(Advanced, with Selectable nodes")){
		    HelpMarker(
                "This is a more typical looking tree with selectable nodes.\n" +
		       "Click to select, CTRL+Click to toggle, click on arrows or double-click to open.");
			
				
					i.CheckboxFlags("ImGuiTreeNodeFlags_OpenOnArrow",       &base_flags, ImGuiTreeNodelags_OpenOnArrow);
            imgui.CheckboxFlags("ImGuiTreeNodeFlags_OpenOnDoubleClick", &base_flags, ImGuiTreeNodeFlags_OpenOnDoubleClick);
			imgui.CheckboxFlags("ImGuiTreeNodeFlags_SpanFullWidth"&base_flags, ImGuiTreeNodeFlags_SpanFullWidt);
			imgui.Checkbox("Align label with current X position", &align_label_with_current_x_position);
			imgui.Checkbox("Test tree node as drag source", &test_drand_drop);
			
			
			imgui.Text("Hello!");
			if (align_label_with_current_x_position) {
			    imgui.Unindent(imgui.GetTreeNodeToLabelSpacing());
			
			
				selection_mask' is dumb representation of what ma be user-side selection state.
            //  You may retain selection state inside or outside your objects in whatever format you see fit.
            // 'node_clicked' is temporary storage of what node we have clicked to process selection at the end
			/// of the loop. May be a pointer to your own node type, etc.
			var node_clicked = -1;
			for i := 0; i < 6; i++ {
			    // Disable the default "open on single-click behavior" + set Selected flag according to our selection.
			    var node_flags = ase_flags;
			    const bool is_selected = (selection_mask & (1 << i)) != 0;
				if (is_selected) {
				    node_flags |= ImGuiTreeodeFlags_Selected;
				
				if i < 3) {
					// Items 0..2 are Tree Node
                    var node_open = imgui.TreeNodeEx(uintptr(i), node_flags, "Selectable Node %d", i);
				   if (igui.IsItemClicked()) {
					    node_clicked = i;
					
					if test_drag_and_drop &&imgui.BeginDragDropSource()){
						imgui.SetDragDroPayload("_TREENODE", NULL, 0);
                        imgui.Text("This is a drag and drop source");
					   imgui.EndDragDropSource(); 
						
						node_open){
						imgui.BulletText("Blah blh\nBlah Blah");
					    imgui.TreePop();
					} 
						
						tems 3..5 are Tee Leaves
					// The only reason we use TreeNode at all is to allow selection of the leaf. Otherwise we can
				    // use BulletText() or advance the cursor by GetTreeNodeToLabelSpacing() and call Text().
					node_flags |= ImGuiTreeNodeFlags_Leaf | ImGuiTreeNodeFlags_NoTreePushOnOpen; // ImGuiTreeNodeFlags_Bullet
					imgui.TreeNodeEx(uintptr(i), node_flags, "Selectable Leaf %d", i);
					if (imgui.IsItemClicked()) {
					    node_clicked = i;
					
					if test_drag_and_drop &&imgui.BeginDragDropSource()) {
						imgui.SetDragDroPayload("_TREENODE", NULL, 0);
                        imgui.Text("This is a drag and drop source");
					   imgui.EndDragDropSource();
						
						
						
					_clicked != -1) {
				// Update selection state
			    // (process outside of tree loop to avoid visual inconsistencies during the clicking frame)
			   if (imgui.GetIO().eyCtrl) {
				    selection_mask ^= (1 << node_clicked);          // CTRL+click to toggle
				f (!(selection_mask & (1 << node_clicked))) // Depending on selection behavior you want, may want to preserve selection when clicking on item that is part of the selection
				   selection_mask = (1 < node_clicked);           // Click to single-select
					
            }
					n_label_with_current_x_position) {
                imgui.Indent(imgui.GetTreeNodeToLabelSpacing());
			
			imgi.TreePop();
				
        imgui.TreePop();
			
		
		imgui.TreeNode(Collapsing Headers")) {
	    imgui.Checkbox("Show 2nd header", &closable_group);
        if (imgui.CollapsingHeader("Header", ImGuiTreeNodeFlags_None)){
	       imgui.Text("IsItemHovered: %d", mgui.IsItemHovered());
		    for i := 0; i < 5; i++ {
		       imgui.Text("Some content %d", i); 
			
			
				i.CollapsingHeader("Header with  close button", &closable_group)){
            imgui.Text("IsItemHovered: %d", imgui.IsItemHovered());
		    for i := 0; i < 5; i++ {
		       imgui.Text("More content %d", i); 
			
			
				
        if (imgui.CollapsingHeader("Header with a bullet", ImGuiTreeNodeFlags_Bullet))
		    imgui.Text("IsItemHovered: %d", imgui.IsItemHovered());
		*/
		   imgui.TreePop();
		
		
		imgui.TreeNode(Bullets")) {
	    imgui.BulletText("Bullet point 1");
        imgui.BulletText("Bullet point 2\nOn multiple lines");
	   if (imgui.TreeNode("Tree ode")) {
		    imgui.BulletText("Another bullt point");
		    imgui.TreePop();
		}
			i.Bullet(); imgui.Text("Bullet point 3 (wo calls)");
			i.Bullet(); imgi.SmallButton("Button");
		imgui.TreePop();
		
		
		
		
		imgui.TreeNode(Text")) {
	    if (imgui.TreeNode("Colorful Text")) {
            // Using shortcut. You can use PushStyleColor()/PopStyleColor() for more flexibility.
	       imgui.TextColored(mVec4(1.0, 0.0, 1.0, 1.0), "Pink");
		   imgui.TextColored(ImVec4(1.0, 10, 0.0, 1.0), "Yellow");
			imgui.TextDisabled("Disabled");
			imgui.SameLine(); HelpMarker("The TextDisabled color s stored in ImGuiStyle.");
			imgui.TreePop();
			
			
			
			imgui.TreeNode(Word Wrapping")) {
		    // Using shortcut. You can use PushTextWrapPos()/PopTextWrapPos() for more flexibility.
            imgui.TextWrapped(
		       "This text should automaticlly wrap on the edge of the window. The current implementation " +
			    "for text wrapping follows simple rules suitable for English and possibly other languages.");
			imgui.Spacing();
				
					i.SliderFloat("Wrap width", &wrap_width, -20, 600, "%.0");
			
            ImDrawList* draw_list = imgui.GetWindowDrawList();
			for n := 0; n < 2; n++ {
                imgui.Text("Test paragraph %d:", n);
			    var po s = imgui.GetCursorScreenPos();
			    var marker_min = ImVec2(pos.X + wrap_width, pos.Y);
				var marker_max = ImVec2(pos.X + wra_width + 10, pos.Y + imgui.GetTextLineHeight());
				imgui.PushTextWrapPos(imgui.GetCursoPos().x + wrap_width);
				if (n == 0) {
				    imgui.Text("The lazy dog sa good dog his paragah should fit within %.0 piels. Testing a 1 character word. The quick brown fox jumps over the lazy dog.", wrap_width);
				
				   imgui.ext("aaaaaaaa bbbbbbbb, c cccccccc,dddddddd. d eeeeeeee   ffffffff. gggggggg!hhhhhhhh");
					

					raw actual text bounding box, following by marker of our expected limit (should not overlap!)
                draw_list.AddRect(imgui.GetItemRectMin(), imgui.GetItemRectMax(), IM_COL32(255, 255, 0, 255));
                draw_list.AddRectFilled(marker_min, marker_max, IM_COL32(255, 0, 255, 255));
				imgui.PopTextWrapPos();
				
				
				i.TreePop();
			

			imgui.TreeNode(UTF-8 Text")) {
		    // UTF-8 test with Japanese characters
            // (Needs a suitable font? Try "Google Noto" or "Arial Unicode". See docs/FONTS.md for details.)
		   // - From C++11 you can use he u8"my text" syntax to encode literal strings as UTF-8
			// - For earlier compiler, you may be able to encode your sources as UTF-8 (e.g. in Visual Studio, you
			//   can save your source files as 'UTF-8 without signature').
			// - FOR THIS DEMO FILE ONLY, BECAUSE WE WANT TO SUPPORT OLD COMPILERS, WE ARE *NOT* INCLUDING RAW UTF-8
			//   CHARACTERS IN THIS SOURCE FILE. Instead we are encoding a few strings with hexadecimal constants.
			//   Don't do this in your application! Please use u8"text in any language" in your application!
			// Note that characters values are preserved even by InputText() if the font cannot be displayed,
			// so you can safely copy & paste garbled characters into another application.
			imgui.TextWrapped(
			    "CJK text will only appears if the font was loaded with the appropriate CJK character ranges. " +
			    "Call io.Fonts.AddFontFromFileTTF() manually to load extra character ranges. " +
			    "Read docs/FONTS.md for details.");
				i.Text("Hiragana: \xe3\x81\x8b\xe3\x81\x8d\xe3\x81\x8f\xe3\x81\x91\xe3\x81\x93 (kakikukeko)"); // Normally we would use u8"blah blah" with the proper characters directly in the string.
					i.Text("Kanjis: \xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e (nihongo)");
					atic char buf[32] = u8"NIHONGO"; / <- this is how you would write it with C++11, using real kanjis
			imgui.InputText("UTF-8 input", &buf, IM_ARRAYSIZE(buf));
			imgui.TreePop();
			
			i.TreePop();
			
		
		imgui.TreeNode(Images")) {
	    var io = imgui.GetIO();
        imgui.TextWrapped(
	       "Below we are displaing the font texture (which is the only texture we have access to in this demo). " +
		    "Use the 'ImTexturID' type as storage to pass pointers or identifier to your own texture data. " +
		    "Hover the texture for a zoomed view!");
			
				elow we are displaying the font texture because it is the only texture we have access to inside the demo!
				emember that ImTextureID is just storag for whatever you want it to be. It is essentially a value that
        // will be passed to the rendering backend via the ImDrawCmd structure.
		// If you use one of the default imgui_impl_XXXX.cpp rendering backend, they all have comments at the top
		// of their respective source file to specify what they expect to be stored in ImTextureID, for example:
		// - The imgui_impl_dx11.cpp renderer expect a 'ID3D11ShaderResourceView*' pointer
		// - The imgui_impl_opengl3.cpp renderer expect a GLuint OpenGL texture identifier, etc.
		// More:
		// - If you decided that ImTextureID = MyEngineTexture*, then you can pass your MyEngineTexture* pointers
		//   to imgui.Image(), and gather width/height through your own functions, etc.
		// - You can use ShowMetricsWindow() to inspect the draw data that are being passed to your renderer,
		//   it will help you debug issues if you are confused about it.
		// - Consider using the lower-level ImDrawList::AddImage() API, via imgui.GetWindowDrawList().AddImage().
		// - Read https://github.com/ocornut/imgui/blob/master/docs/FAQ.md
		// - Read https://github.com/ocornut/imgui/wiki/Image-Loading-and-Displaying-Examples
		var my_tex_id = io.Fonts.TexID;
		var my_tex_w = float32(io.Fonts.TexWidth);
		var my_tex_h = float32(io.Fonts.TexHeight);
		{
		    imgui.Text("%.0x%.0", my_tex_w, my_te_h);
		    var pos = imgui.GetCursorScreenPos();
		    var uv_min = ImVec2(0.0, 0.0);                 // Top-left
			var uv_max = ImVec2(1.0, 1.0);                // Lower-right
			var tint_col = ImVec4(1.0, 1.0, 1.0,1.0);   // No tint
			var border_col = ImVec4(1.0, .0, 1.0, 0.5);  50% opaque white
			imgui.Image(my_tex_id, ImVec2y_tex_w, my_texh), uv_min, uv_max, tint_col, border_col);
			if (imgui.IsItemHovered()) {
			    imgui.BeginTooltip();
			    var region_sz = 32.0;
			   var region_x = io.MouePos.x - pos.x - region_sz * 0.5;
				var region_y = io.MosePos.y - pos.y - region_sz * 0.5;
				var zoom = 4.0;
				if (region_x < 0.0) { region_x = 0.0; } else if rgio_x > my_tex_w - region_sz) { region_x = my_tex_w - region_sz; }
				if (region_y < 0.0) { region_y = 0.0 } else if (eiony > my_tex_h - region_sz) { region_y = my_tex_h - region_sz; }
				imgui.Text("Mi: (%.2, %.2)", region_x, region_y);
				imgi.Text("Max: (.2
					 %.2)", region
				 + region_z, region_y + regio_z);
					
				
				varuv0 = ImVec2((eg
					on_x) / my_tex
				w, (regiony) / my_tex_h);
					
				
				var uv1 = ImVec2((region_x + region_sz) / my_tex_, (region_y + region_sz) / my_tex_h);
				imgui.Image(my_tex_id, ImVec2(region_s  zoom, region_sz * om), uv0, u1, tint_col, border_col);
				imgui.EndTooltip();
				
				
				xtWrapped("And nowsome textured buttons..");
			i := 0; i < 8; i++ {
		    imgui.PushID(i);
		    var frame_padding = -1 + i;                            // -1 == uses default padding (style.FramePadding)
		    var size = ImVec2(32.0, 32.0);                     // Size of the image we want to make visible
			var uv0 = ImVec(0.0, 0.0);                        // UV coordinates for lower-left
			var uv1 = ImVec2(32.0 / my2.0 / my_tex_h);// UV coordinates for (32,32) in our texture
			var bg_col = ImVec4(0.0, 0.0,, 1.0);         // Black background
			var tint_col = ImVec4(1.0,, 1.0, 1.0);       // No tint
			if (imgui.ImageButtonm_tex_id, size,u0, uv1, f ame_padding, bg_col, tint_col)) {
			    pressed_count += 1;
			
			imgi.PopID();
				i.SameLine();
        }
			i.NewLine();
			i.Text("Pressed d times.", pressed_count);
		imgui.TreePop();
		
		
		imgui.TreeNode(Combo")) {
	    // Expose flags as checkbox for the demo
        var ImGuiComboFlags flags = 0;
	   imgui.CheckboxFlags("ImuiComboFlags_PopupAlignLeft", &flags, ImGuiComboFlags_PopupAlignLeft);
		imgui.SameLine(); HelpMarker("Only makes a difference if the popup is larger than the combo");
		if (imgui.CheckboxFlags("ImGuComboFlags_NoArrowButton", &flags, ImGuiComboFlags_NoArrowButton)) {
		    flags &= ^ImGuiComboFlags_NoPreview;     // Clear the other flag, as we cannot combine boh
		
		
		if imgui.CheckboxFlags("ImGuiComboFlags_NoPreview", &flags, ImGuiComboFlags_NoPreview)) {
			flags &= ^ImGuiComboFlags_NoArrowBu // Clear the other flag, as we cannot combine both
		}
		
			sing the generic BeginCombo() API, you ave full control over how to display the combo contents.
        // (your selection data could be an index, a pointer to the object, an id for the object, a flag intrusively
        // stored in the object itself, etc.)
		var items = [...]string{ "AAAA", "BBBB", "CCCC", "DDDD", "EEEE", "FFFF", "GGGG", "HHHH", "IIII", "JJJJ", "KKKK", "LLLLLLL", "MMMM", "OOOOOOO" };
		var combo_preview_value = items[item_current_idx];  // Pass in the preview value visible before opening the combo (it could be anything)
		if (imgui.BeginCombo("combo 1", combo_preview_value, flags)) {
		    for n := 0; n < IM_ARAYSIZE(items); n++ {
		        const bool is_selected = (item_current_id== n);
		       if (imgui.Selectable(items[n], is_selected)) {
			        item_current_idx = n;
				
				
					et the initial focuswhen opening the combo (scrolling + keyboard navigation focus)
                if (is_selected) {
                    imgui.SetItemDefaultFocus();
				
				
					dCombo();
        }
			
			implified one-lier Combo() API, using values packed in a single constant string
		// This is a convenience for when the selection set is small and known at compile-time.
        imgui.Combo("combo 2 (one-liner)", &item_current_2, "aaaa\000bbbb\000cccc\000dddd\000eeee\000\000");
		
		// Simplified one-liner Combo() using an array of const char*
		// This is not very useful (may obsolete): prefer using BeginCombo()/EndCombo() for full control.
        
		imgui.Combo("combo 3 (array)", &item_current_3, items, IM_ARRAYSIZE(items));
		
// Simplified one-liner Combo() using an accessor function
		ncs = struct {
			ItemGetter func(data *string, n int, out string) bool
		
			ItemGetter: func(data *string, n int, out *string) bool {
				*out = data[n]; 
				return true;
			},
		}
        imgui.Cmbo("combo 4 (function)", &item_current_4, &Funcs.ItemGetter, items, IM_ARRAYSIZE(items));

        imgui.TreePop();
		

		imgui.TreeNode(List boxes")) {
	    // Using the generic BeginListBox() API, you have full control over how to display the combo contents.
        // (your selection data could be an index, a pointer to the object, an id for the object, a flag intrusively
	   // stored in the object itslf, etc.)
		var items = [...]string { "AAAA", "BBBB", "CCCC", "DDDD", "EEEE", "FFFF", "GGGG", "HHHH", "IIII", "JJJJ", "KKKK", "LLLLLLL", "MMMM", "OOOOOOO" };
		
		if (imgui.BeginListBox("listbox 1")) {
		    for n := 0; n < IM_RAYSIZE(items); n++ {
                const bool is_selected = (item_current_idx == n);
		       if (imgui.Selectable(items[], is_selected)) {
			        item_current_idx = n;
				
				
					et the initial focuswhen opening the combo (scrolling + keyboard navigation focus)
                if (is_selected) {
                    imgui.SetItemDefaultFocus();
				
				
					dListBox();
        }
			
			ustom size: use al width, 5 items tall
		imgui.Text("Full-width:");
        if (imgui.BeginListBox("##listbox 2", ImVec2(-FLT_MIN, 5 * imgui.GetTextLineHeightWithSpacing()))) {
		    for n := 0; n < IM_ARRAYSIZE(items); n++ {
		        const bool is_selcted = (item_current_idx == n);
		       if (imgui.Selectable(items[n], is_selected)) {
			        item_current_idx = n;
				
				
					et the initial focuswhen opening the combo (scrolling + keyboard navigation focus)
                if (is_selected) {
                    imgui.SetItemDefaultFocus();
				
				
					dListBox();
        }
			
			i.TreePop();
		

		imgui.TreeNode(Selectables")) {
	    // Selectable() has 2 overloads:
        // - The one taking "bool selected" as a read-only selection information.
	   //   When Selectable() has ben clicked it returns true and you can alter selection state accordingly.
		// - The one taking "bool* p_selected" as a read-write selection information (convenient in some cases)
		// The earlier is more flexible, as in real application your selection may be stored in many different ways
		// and not necessarily inside a bool value (e.g. in flags within objects, as an external list, etc).
		if (imgui.TreeNode("Basic")) {
		    
		    imgui.Selectable("1. I am selectable", &selection[0]);
		   imgui.Selectable("2. I m selectable", &selection[1]);
imgui.Text("(I am not selectable)");
			imgui.Selectable("4. I am selectable", &selection[3])
			if (imgui.Selectable("5. I am double clickable", seletion[4], ImGuiSelectableFlags_AllowDoubleClick)) {
			    if (imgui.IsMouseDoubleClicked()) {
			        selection[4] = !selection[4];
			
				
					eePop();
        }
        if (imgui.TreeNode("Selection State: Single Selection")) {
			for n := 0; n <5; n++ {
		        if (imgui.Selectable(fmt.Sprintf("Object %d", n), selected == n)) {
		           selected = n;
			
				
					eePop();
        }
			imgui.TreeNode("Selection State: Multiple Selection")) {
			HelpMarker("Hol CTRL and click to select multiple items.");
		    for n := 0; n < 5; n++ {
		       if (imgui.Selectable(fmt.Sprintf("Object %d", n),selection2[n])) {
			        if (!imgui.GetIO().KeyCtrl) {    // Clear selectionwhen CTRL is not held
			            memset(selection2, 0, sizeof(selection2));
				[n]^= 1;
					
						
            }
            imgui.TreePop();
				
			imgui.TreeNode("Rendering more text into the same line")) {
			// Using the Seectable() override that takes "bool* p_selected" parameter,
		    // this function toggle your bool value automatically.
		   imgui.Selectable("main.c",    &selected2[0]); imgui.Sameine(300); imgui.Text(" 2,345 bytes");
			imgui.Selectable("Hello.cpp", &selected2[1]); imgui.SameLine(300); imgui.Text("12,345 bytes");
			imgui.Selectable("Hello.h",   &selected2[2]); imgui.SameLine(300); imgui.Text(" 2,345 bytes");
			imgui.TreePop();
			
			
			
			
			
			imgui.TreeNode("In columns" {
			
			
			if (imgui.Beginable("split1", 3, ImGuiTableFlags_Resizable | ImGuiTableFlags_NoSavedSettings | ImGuiTableFlags_Borders)) {
		        for i := 0; i < 10; i++ {
		           imgui.TableNextColum();
			       imgui.Selectable(fmt.Sprintf("Item %d", i), &selectdolumns[i]); // FIXME-TABLE: Selcion overlap
				}
					i.EndTable();
					
				i.Spacing();
				imgui.BeginTable"split2", 3, ImGuiTableFlags_Resizable | ImGuiTableFlags_NoSavedSettings | ImGuiTableFlags_Borders)) {
			    for i := 0; i < 10; i++ {
			        imgui.TbleNextRow();
			       imgui.TableNextColumn();
				    imgui.Selectable(fmt.Sprintf("Item %d", i), &selectedColumns[i], ImGuiSelectableFlags_SpanAllColumns);
					imgui.TableNextColum();
					imgui.Text("Some other ontents");
					imgui.TableNextColumn();
					imgui.Text("123456");
					
					i.EndTable();
					
				i.TreePop();
				
			imgui.TreeNode("Grid")) {
			// Add in a bitof silly fun...
		    const float time = float32(imgui.GetTime());
		   const bool winning_stte = memchr(selectedGrid, 0, sizeof(selectedGrid)) == NULL; // If all cells are selected...
			if (winning_state) {
			    imgui.PushStyleVar(ImGuiStyleVar_SelectbleTextAlign, ImVec2(0.5 + 0.5 * cosf(time * 2.0), 0.5 + 0.5 * sinf(time * 3.0)));
			
			
				y := 0; y < 4; y++ {
                for x := 0; x < 4; x++ {
                    if (x > 0) {
			            imgui.SameLine();
				
					
						i.PushID(y * 4 +x);
                    if (imgui.Selectable("Sailor", selectedGrid[y][x] != 0, 0, ImVec2(50, 50))) {
                        // Toggle clicked cell + toggle neighbors
					    selectedGrdy][x] = 1;
					   if (x > 0) { selectedGrid[y][x - 1] ^= 1; }
						if (x < 3) { selectedGrid[y][x + 1] ^= 1; }
						if (y > 0) { selectedGrd[y - 1][x] ^= 1; }
						if y < 3 {
							selectedGrid[y + ]x] ^= 1
						}
						
							
						
						i.PpID()
							
						
						
							
						
					
					ing_state) {
				imgui.PopStyleVar();
			}
				imgi.TreePop();
					
	        if (imgui.TreeNode("Alignment")) {
				HelpMarker(
			        "By default, Selectables uses style.SelectableTextAlign but it can be overridden on a per-item " +
			       "basis using PushStyleVr(). You'll probably want to always keep your default situation to " +
				    "left-align otherwise it becomes difficult to layout multiple items on a same line");
					
						y := 0; y < 3; y++ {
						for x := 0; x < 3; x++ {
                    var alignment = ImVec2(float32(x / 2.0), float32(y / 2.0));
				        var name = fmt.Sprintf("(%.1,%.1)", alignment.x, alignment.y)
					    if (x > 0) {imgui.SameLine(); }
						imgui.PushStyleVar(ImGuiStyleVarSlectableTextAlin aligment);
						imgui.Selectable(name, &selectedAlignment[3 * y + x], ImGuiSelectableFlags_None, ImVec2(80, 80));
						imgi.Popty
							leVar();
						
						
						
						eePop();
					
				i.TreePop();
				
			
			o wire InputTex() with std::string or any other custom string type,
		// see the "Text Input > Resize Callback" section of this demo, and the misc/cpp/imgui_stdlib.h file.
    if (imgui.TreeNode("Text Input")) {
		    if (imgui.TreeNode("Multi-line Text Input")) {
		        HelpMarker("You can use the ImGuiInputTextFlags_CallbackResize facility if you need to wire InputTextMultiline() to a dynamic string type. See misc/cpp/imgui_stdlib.h for an example. (This is not demonstrated in imgui_demo.cpp because we don't want to include <string> in here)");
		       imgui.CheckboxFlags("ImuiInputTextFlags_ReadOnly", &flags, ImGuiInputTextFlags_ReadOnly);
			   imgui.CheckboxFlags("ImGuiInputTextFlag_AllowTabInput", &flags, ImGuiInputTextFlags_AllowTabInput);
				imgui.CheckboxFlags("ImGuiInputTextFlags_CtrlEnterForNewLine", &flags, ImGuiInputTextFlags_CtrlEnterForNewLine);
				imgui.InputTextMultiline("##source", text, IM_ARRAYSIZE(text), ImVec2(-FLT_MIN, imgui.GetextLineHeight() * 16), flags);
				imgui.TreePop();
				
				
				imgui.TreeNode(Filtered Text Input")) {
			rImGuiLetters := func(ImGuiInputTextCallbackData* data) {
				if (data.EventChar < 256 && strchr("imgui", data.EventChar)) {
			urn0;
					}*
					retrn 1;
				}
	
	            mgui.InputText("default",     buf1, 64);
	            imgui.InputText("decimal",     buf2, 64, ImGuiInputTextFlags_CharsDecimal);
            imgui.InputText("hexadecimal", buf3, 64, ImGuiInputTextFlags_CharsHexadecimal | ImGuiInputTextFlags_CharsUppercase);
				imgui.InputText("uppercase"buf4, 64,ImGuiInputTextFlags_CharsUppercase);
				imgui.InputText("no blank",buf5, 64, ImGuiInputTextFlags_CharsNoBlank)
				imgui.InputText("\"imgui\" letters", buf6, 64, ImGuiInputTextFlags_CallbackChrilter, TextFilters.FilterImGuiLettes);
				imgui.TreePop();
				
				
				imgui.TreeNode(Password Input")) {
			    
            imgui.InputText("password", password, IM_ARRAYSIZE(password), ImGuiInputTextFlags_Password);
			   imgui.SameLine(); HelpMarker("Diplay all characters as '*'.\nDisable clipboard cut and copy.\nDisable logging.\n");
imgui.InputTextWithHint("password (w/ hint)", "<password>", password, IM_ARRAYSIZE(password), ImGuiInputTextFlags_Password);
				imgui.InputText("password (clear)", password, IM_ARRAYSIZE(password));
				imgui.TreePop();
				
				
				
				imgui.TreeNode(Completion, History, Edit Callbacks")) {
			lback := func(ImGuiInputTextCallbackData* data) {
				if (data.EventFlag == ImGuiInputTextFlags_CallbackCompletion) {
			a.IsertChars(data.CursorPos, "..");
					} else if (data.EventFlag == ImGuiInputTextFa*gs_CallbackHistory) {
						if(data.EventKey == ImGuiKey_UpArrow) {
							data.DeleteChars(0, data.BufTextLen);
							data.InsrtChars(0, "Pressed Up!");
							daa.SelectAll();
						} else if (data.EventKey == ImGuiKey_ownArrow) {
							data.DeleteChars(0, data.BufTextLe);
							data.InsertChars0, "Pressed Down!");
							data.SeletAll();
						}
					} else if (data.EventFlag == ImGuiInpuTextFlags_CallbackEdit) {
						// Toggle casing f first character
						var c = data.Buf[0];
						if ((c >='a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
							data.Buf[0] ^= 32;
						}
						dat.BufDirty = true;
	
						// Increment a counter
						int* p_int = data.UsrData.(*int);
					*p_int = *p_int + 1;
					}
				r	etu rn 0;
			}	
	
	            mgui.InputText("Completion", buf7, 64, ImGuiInputTextFlags_CallbackCompletion, MyCallback);
	            imgui.SameLine(); HelpMarker("Here we append \"..\" each time Tab is pressed. See 'Examples>Console' for a more meaningful demonstration of using this callback.");

				imgui.InputText("History", buf8, 64, ImGuiInputTextFlags_CallbackHistory, MyCallback);
				imgui.SameLine()
				HelpMarker("Here we replace and select text each time Up/Down are pressed. See 'Examples>Console' for a more meaningful demonstration of using tis callback.");

				imgui.InputText("Edit", buf9, 64, ImGuiInputTextFlags_CallbackEdit, MyCallback, &editcount);
				imgui.SameLine()
				HelpMarker("Here we toggle the casing of the first character on every edits + count edits.");
            imgui.SameLine(); imgui.Text("(%d)", edit_count);
				
				imgui.TreePop();
				
				
				

				imgui.TreeNode(Resize Callback")) {
			    // To wire InputText() with std::string or any other custom string type,
            // you can use the ImGuiInputTextFlags_CallbackResize flag + create a custom imgui.InputText() wrapper
			   // using your preferred type. Seemisc/cpp/imgui_stdlib.h for an implementation of this using std::string.
				HelpMarker(
				    "Using ImGuiInputTextFlags_CallbackResize to wire your custom string type to InputText().\n\n" +
				    "See misc/cpp/imgui_stdlib.h for an implementation of this for std::string.");
				
					ack := func(ImGuiInputTextCallbackData* data) {
						ntFlag == ImGuiInputTextFlags_CallbackResize) {
					var my_str = data.UserData.(*[]string);
						IM_ASSERT(my_str.begin() == data.Buf);*
						mystr.resize(data.BufSize); // NB: On resizing calls, enerally data.BufSize == data.BufTextLen + 1
						data.Buf = my_str.begin();
					}
					return 0;
			}	
	
			M	yInputTetMultiline := func(label string, my_str *[]string, size ImVec2, flags ImGuiInputTextFlags) {
					IM_ASSERT((flags & ImGuiInputTextFlags_CallbackResize) == 0);
				return imgui.InputTextMultiline(label, my_str.begin(), len(my_str), size, flags | ImGuiInputTextFlags_CallbackResize, MyResizeCallback, my_str);
				}
	
	            // For this demo we are using ImVector as a string container.
	            // Note that because we need to store a terminating zero character, our size/capacity are 1 more
            // than usually reported by a typical string class.
				if (my_str.empty()) {
				    my_str.push_back(0);
				
				MyIputTextMultilie("##MyStr", &my_str, ImVec2(-FLT_MIN, imgui.GetTextLineHeight() * 16));
					i.Text("Data: %p\nSze: %d\nCapacity: %d", my_str, len(my_str), len(my_str.capacity));
	            imgui.TreePop();
				
				
				i.TreePop();
			

			abs
		if (imgui.TreeNode("Tabs")) {
        if (imgui.TreeNode("Basic")) {
		        var tab_bar_flags = ImGuiTabBarFlags_None;
		       if (imgui.BeginTaBar("MyTabBar", tab_bar_flags)) {
			       if (imgui.BeginTabIem("Avocado")) {
				        imgui.Text("This is the Avocado tb!\nblah blah blah blah blah");
				       imgui.EndTabItem();
					}
						imgui.BeginTabItem("Broccoli")) {
						imgui.Text("This i the Broccoli tab!\nblah blah blah blah blah");
					    imgui.EndTabItem();
					}
						imgui.BeginTabItem("Cucumber")) {
						imgui.Text("This i the Cucumber tab!\nblah blah blah blah blah");
					    imgui.EndTabItem();
					}
						i.EndTabBar();
						
					i.Separator();
					i.TreePop();
				
				
				imgui.TreeNode(Advanced & Close Button")) {
			    // Expose a couple of the available flags. In most cases you may just call BeginTabBar() with no flags (0).
            imgui.CheckboxFlags("ImGuiTabBarFlags_Reorderable", &tab_bar_flags, ImGuiTabBarFlags_Reorderable);
			   imgui.CheckboxFlags("ImGuiTabBarFlags_AutSelectNewTabs", &tab_bar_flags, ImGuiTabBarFlags_AutoSelectNewTabs);
				imgui.CheckboxFlags("ImGuiTabBarFlags_TabListPopupButton", &tab_bar_flags, ImGuiTabBarFlags_TabListPopupButton);
				imgui.CheckboxFlags("ImGuiTabBarFlags_NoCloseWithMiddleMouseButton", &tab_bar_flags, ImGuiTabBarFags_NoCloseWithMiddleMouseButton);
				if ((tab_bar_flags & ImGuiTabBarFlags_FittingPolicyMask_) == 0) {
				    tab_bar_flags |= ImGuiTabBarFlags_FittingPolicyDefault_;
				
				if imgui.CheckboxFlags("ImGuiTabBarFlags_FittingPolicyResizeDwn", &tab_bar_flags, ImGuiTabBarFlags_FittingPolicyResizeDown)) {
					tab_bar_flags &= ^(ImGuiTabBarFlags_FittingPolicyMask_  ImGuiTabBarFlags_FittingPolicyResizeDown);
				}
				if imgui.CheckboxFlags("ImGuiTabBarFlags_FittingPolicyScroll", &tab_bar_flags, ImGuiTabBarFlags_FittingPolicyScroll)) {
					tab_bar_flags &= ^(ImGuiTabBarFlags_FittingPolicyMask_ ^ ImGuiTabBarFlags_FittingPolicyScroll);
				}
				
					ab Bar
	            var names = [...]string{ "Artichoke", "Beetroot", "Celery", "Daikon" };
            for n := range opened {
				    if (n > 0) { imgui.SameLine(); }
				    imgui.Checkbox(namesn], &opened[n]);
				}
					
						
					
					assing a bool* to BeginTabItem() is imilar to passing one to Begin():
				// the underlying bool will be set to false when the tab is closed.
            if (imgui.BeginTabBar("MyTabBar", tab_bar_flags)) {
				 range opened {
				        if (opened[n] && imgui.BeginTabItem(names[n], &opened[n], ImGuiTabItemFlags_None)) {
				           imgui.Text("This is the %s tab!", naes[n]);
	                        if (n & 1) {
						       imgui.Text("I am an odd tab.");
							
							imgi.EndabItem();
								
	                imgui.EndTabBar();
							
						parator();
						i.TreePop();
					
					
					imgui.TreeNode(TabItemButton & Leading/Trailing flags")) {
				    if (next_tab_id == 0) { // Initialize with some default tabs
                for i := 0; i < 3; i++ {
				t_tb_id++
					       active_tabs  append(active_tabs, next_tab_id);
						
			}		
							
		            // TabItemButton() and Leading/Trailing flags are distinct features which we will demo together.
		            // (It is possible to submit regular tabs with Leading/Trailing flags, or TabItemButton tabs without Leading/Trailing flags...
            // but they tend to make more sense together)
					imgui.Checkbox("Show Leading TabItemButton()", &show_leading_button);
					imgui.Checkbox("Show Trailing TabItemButton()", &show_trailing_button);
					
					// Expose some other flags which are useful to showcase how they intract with Leading/Trailing tabs
					imgui.CheckboxFlags("ImGuiTabBarFlags_TabListPopupButton", &tab_bar_flgs, ImGuiTabBarFlags_TabListPopupButton);
            if (imgui.CheckboxFlags("ImGuiTabBarFlags_FittingPolicyResizeDown", &tab_bar_flags, ImGuiTabBarFlags_FittingPolicyResizeDown)) {
					    tab_bar_flags &= ^(ImGuiTabBarFlags_FittingPolicyMask_ ^ ImGuiTabBarFlags_FittingPolicyResizeDown);
					
					if imgui.CheckboxFlags("ImGuiTabBarFlags_FittingPolicyScroll", &tab_bar_flags, ImGuiTabBarFlags_FittingPolicyScroll)) {
						tab_bar_flags &= ^(ImGuiTabBarFlags_FittingPolicyMask_ ^ ImGuiTabBarFlags_FittingPolicyScroll);
					}
					
						imgui.BeginTabBar("MyTabBar", tab_bar_flags)) {
		                // Demo a Leading TabItemButton(): click the "?" button to open a menu
                if (show_leading_button) {
					       if (imgui.TabItemButton("?", ImGuiTabIteFlags_Leading | ImGuiTabItemFlags_NoTooltip)) {
						        imgui.OpenPopup("MyHelpMenu");
						
							
								i.BeginPopup("MyHelpMenu")) {
		                    imgui.Selectable("Hello!");
    		                imgui.EndPopup();
						}
							
							emo Trailing Tab: click the "+" button to add a new tab (in your app you may want to use a font icon instead of the "+")
						// Note that we submit it before the regular tabs, but because of the ImGuiTabItemFlags_Trailing flag it will always appear at the end.
                if (show_trailing_button) {
						    if (imgui.TabItemButton("+", ImGuiTabItemFlags_Trailing | ImGuiTabItemFlags_NoTooltip)) {
						d++
						s =append(active_tabs, ext_tab_id); // Add new tab
							
						}
		
     		           // Submit our regular tabs
		                for m := range active_tabs {
                    var open = true;
						    var name = fmt.Sprintf("%04d", active_tabs[n])
						    if (imgui.BeginTabItem(name, &open, ImGuiTabItemFlags_None)) {
							    imgui.Text(This is the %s tab!", name);
							    imgui.EndTabItem();
							}
								
								!open) {
							append(active_tabs[:n], active_tabs[n+1:]...)
					} else {
							   n++;
							}
     		           }
								
     		           imgui.EndTabBar();
						
            imgui.Separator();
						i.TreePop();
					
					i.TreePop();
					
				}
				imgui.TreePop()
			}

			// Plot/Graph widgets are not very good.
			// Consider using a third-party library such as ImPlot: https://github.com/epezent/implot
			// (see others https://github.com/ocornut/imgui/wiki/Useful-Extensions)
			if imgui.TreeNode("Plots Widgets") {
				imgui.Checkbox("Animate", &animate)

				// Plot as lines and plot as histogram
				imgui.PlotLines("Frame Times", arrPlot, len(arrPlot))
				imgui.PlotHistogram("Histogram", arrPlot, len(arrPlot), 0, NULL, 0.0, 1.0, ImVec2(0, 80.0))

				// Fill an array of contiguous float values to plot
				// Tip: If your float aren't contiguous but part of a structure, you can pass a pointer to your first float
				// and the sizeof() of your structure in the "stride" parameter.

				if !animate || refresh_time == 0.0 {
					refresh_time = imgui.GetTime()
				}

				for refresh_time < imgui.GetTime() { // Create data at fixed 60 Hz rate for the demo
					values[values_offset] = cosf(phase)
					values_offset = (values_offset + 1) % IM_ARRAYSIZE(values)
					phase += 0.10 * values_offset
					refresh_time += 1.0 / 60.0
				}

				// Plots can display overlay texts
				// (in this example, we will display an average value)
				{
					for n := range values {
						average += values[n]
					}
					average /= float32(len(values))
					var overlay = fmt.Sprintf("avg %f", average)
					imgui.PlotLines("Lines", values, len(values), values_offset, overlay, -1.0, 1.0, ImVec2(0, 80.0))
				}

				// Use functions to generate output
				// FIXME: This is rather awkward because current plot API only pass in indices.
				// We probably want an API passing floats and user provide sample rate/count.
				var fn func(i int)
				Sin := func(i int) {
					return sinf(i * 0.1)
				}
				Saw := func(i int) {
					if i & 1 {
						return 1.0
					} else {
						return -1.0
					}
				}

				imgui.Separator()
				imgui.SetNextItemWidth(imgui.GetFontSize() * 8)
				imgui.Combo("func", &func_type, "Sin\000Saw\000")
				imgui.SameLine()
				imgui.SliderInt("Sample count", &display_count, 1, 400)
				switch func_type {
				case 0:
					fn = Sin
				case 1:
					fn = Saw
				}
				imgui.PlotLines("Lines", fn, NULL, display_count, 0, NULL, -1.0, 1.0, ImVec2(0, 80))
				imgui.PlotHistogram("Histogram", fn, NULL, display_count, 0, NULL, -1.0, 1.0, ImVec2(0, 80))
				imgui.Separator()

				// Animate a simple progress bar
				if animate {
					progress += progress_dir * 0.4 * imgui.GetIO().DeltaTime
					if progress >= +1.1 {
						progress = +1.1
						progress_dir *= -1.0
					}
					if progress <= -0.1 {
						progress = -0.1
						progress_dir *= -1.0
					}
				}

				// Typically we would use ImVec2(-1.0,0.0) or ImVec2(-FLT_MIN,0.0) to use all available width,
				// or ImVec2(width,0.0) for a specified width. ImVec2(0.0,0.0) uses ItemWidth.
				imgui.ProgressBar(progress, ImVec2(0.0, 0.0))
				imgui.SameLine(0.0, imgui.GetStyle().ItemInnerSpacing.x)
				imgui.Text("Progress Bar")

				var progress_saturated = IM_CLAMP(progress, 0.0, 1.0)
				var buf = fmt.Sprintf("%d/%d", int(progress_saturated*1753), 1753)
				imgui.ProgressBar(progress, ImVec2(0, 0), buf)
				imgui.TreePop()
			}

			if imgui.TreeNode("Color/Picker Widgets") {
				imgui.Checkbox("With Alpha Preview", &alpha_preview)
				imgui.Checkbox("With Half Alpha Preview", &alpha_half_preview)
				imgui.Checkbox("With Drag and Drop", &drag_and_drop)
				imgui.Checkbox("With Options Menu", &options_menu)
				imgui.SameLine()
				HelpMarker("Right-click on the individual color widget to show options.")
				imgui.Checkbox("With HDR", &hdr)
				imgui.SameLine()
				HelpMarker("Currently all this does is to lift the 0..1 limits on dragging widgets.")

				var misc_flags int
				if hdr {
					misc_flags |= ImGuiColorEditFlags_HDR
				}
				if !drag_and_drop {
					misc_flags |= ImGuiColorEditFlags_NoDragDrop
				}
				if alpha_half_preview {
					misc_flags |= ImGuiColorEditFlags_AlphaPreviewHalf
				} else if alpha_preview {
					misc_flags |= ImGuiColorEditFlags_AlphaPreview
				}
				if !options_menu {
					misc_flags |= ImGuiColorEditFlags_NoOptions
				}

				imgui.Text("Color widget:")
				imgui.SameLine()
				HelpMarker(
					"Click on the color square to open a color picker.\n" +
						"CTRL+click on individual component to input value.\n")
				imgui.ColorEdit3("MyColor##1", &color, misc_flags)

				imgui.Text("Color widget HSV with Alpha:")
				imgui.ColorEdit4("MyColor##2", &color, ImGuiColorEditFlags_DisplayHSV|misc_flags)

				imgui.Text("Color widget with Float Display:")
				imgui.ColorEdit4("MyColor##2", &color, ImGuiColorEditFlags_Float|misc_flags)

				imgui.Text("Color button with Picker:")
				imgui.SameLine()
				HelpMarker(
					"With the ImGuiColorEditFlags_NoInputs flag you can hide all the slider/text inputs.\n" +
						"With the ImGuiColorEditFlags_NoLabel flag you can pass a non-empty label which will only " +
						"be used for the tooltip and picker popup.")
				imgui.ColorEdit4("MyColor##3", &color, ImGuiColorEditFlags_NoInputs|ImGuiColorEditFlags_NoLabel|misc_flags)

				imgui.Text("Color button with Custom Picker Popup:")

				// Generate a default palette. The palette will persist and can be edited.
				if saved_palette_init {
					for n := range saved_palette {
						imgui.ColorConvertHSVtoRGB(n/31.0, 0.8, 0.8,
							saved_palette[n].x, saved_palette[n].y, saved_palette[n].z)
						saved_palette[n].w = 1.0 // Alpha
					}
					saved_palette_init = false
				}

				var open_popup = imgui.ColorButton("MyColor##3b", color, misc_flags)
				imgui.SameLine(0, imgui.GetStyle().ItemInnerSpacing.x)
				open_popup |= imgui.Button("Palette")
				if open_popup {
					imgui.OpenPopup("mypicker")
					backup_color = color
				}
				if imgui.BeginPopup("mypicker") {
					imgui.Text("MY CUSTOM COLOR PICKER WITH AN AMAZING PALETTE!")
					imgui.Separator()
					imgui.ColorPicker4("##picker", color, misc_flags|ImGuiColorEditFlags_NoSidePreview|ImGuiColorEditFlags_NoSmallPreview)
					imgui.SameLine()

					imgui.BeginGroup() // Lock X position
					imgui.Text("Current")
					imgui.ColorButton("##current", color, ImGuiColorEditFlags_NoPicker|ImGuiColorEditFlags_AlphaPreviewHalf, ImVec2(60, 40))
					imgui.Text("Previous")
					if imgui.ColorButton("##previous", backup_color, ImGuiColorEditFlags_NoPicker|ImGuiColorEditFlags_AlphaPreviewHalf, ImVec2(60, 40)) {
						color = backup_color
					}
					imgui.Separator()
					imgui.Text("Palette")
					for n := range saved_palette {
						imgui.PushID(n)
						if (n % 8) != 0 {
							imgui.SameLine(0.0, imgui.GetStyle().ItemSpacing.y)
						}

						var palette_button_flags = ImGuiColorEditFlags_NoAlpha | ImGuiColorEditFlags_NoPicker | ImGuiColorEditFlags_NoTooltip
						if imgui.ColorButton("##palette", saved_palette[n], palette_button_flags, ImVec2(20, 20)) {
							color = ImVec4(saved_palette[n].x, saved_palette[n].y, saved_palette[n].z, color.w) // Preserve alpha!
						}

						// Allow user to drop colors into each palette entry. Note that ColorButton() is already a
						// drag source by default, unless specifying the ImGuiColorEditFlags_NoDragDrop flag.
						if imgui.BeginDragDropTarget() {
							//TODO PORTING not sure what this is doing.
							/*if (const ImGuiPayload* payload = imgui.AcceptDragDropPayload(IMGUI_PAYLOAD_TYPE_COLOR_3))
							      memcpy((float*)&saved_palette[n], payload.Data, sizeof(float) * 3);
							  if (const ImGuiPayload* payload = imgui.AcceptDragDropPayload(IMGUI_PAYLOAD_TYPE_COLOR_4))
							      memcpy((float*)&saved_palette[n], payload.Data, sizeof(float) * 4);*/
							imgui.EndDragDropTarget()
						}

						imgui.PopID()
					}
					imgui.EndGroup()
					imgui.EndPopup()
				}

				imgui.Text("Color button only:")
				imgui.Checkbox("ImGuiColorEditFlags_NoBorder", &no_border)
				var cflags = misc_flags
				if no_border {
					cflags |= ImGuiColorEditFlags_NoBorder
				}
				imgui.ColorButton("MyColor##3c", color, cflags, ImVec2(80, 80))

				imgui.Text("Color picker:")
				imgui.Checkbox("With Alpha", &alpha)
				imgui.Checkbox("With Alpha Bar", &alpha_bar)
				imgui.Checkbox("With Side Preview", &side_preview)
				if side_preview {
					imgui.SameLine()
					imgui.Checkbox("With Ref Color", &ref_color)
					if ref_color {
						imgui.SameLine()
						imgui.ColorEdit4("##RefColor", &ref_color_v.x, ImGuiColorEditFlags_NoInputs|misc_flags)
					}
				}
				imgui.Combo("Display Mode", &display_mode, "Auto/Current\000None\000RGB Only\000HSV Only\000Hex Only\000")
				imgui.SameLine()
				HelpMarker(
					"ColorEdit defaults to displaying RGB inputs if you don't specify a display mode, " +
						"but the user can change it with a right-click.\n\nColorPicker defaults to displaying RGB+HSV+Hex " +
						"if you don't specify a display mode.\n\nYou can change the defaults using SetColorEditOptions().")
				imgui.Combo("Picker Mode", &picker_mode, "Auto/Current\000Hue bar + SV rect\000Hue wheel + SV triangle\000")
				imgui.SameLine()
				HelpMarker("User can right-click the picker to change mode.")

				var flags = misc_flags
				if !alpha {
					flags |= ImGuiColorEditFlags_NoAlpha
				} // This is by default if you call ColorPicker3() instead of ColorPicker4()
				if alpha_bar {
					flags |= ImGuiColorEditFlags_AlphaBar
				}
				if !side_preview {
					flags |= ImGuiColorEditFlags_NoSidePreview
				}
				if picker_mode == 1 {
					flags |= ImGuiColorEditFlags_PickerHueBar
				}
				if picker_mode == 2 {
					flags |= ImGuiColorEditFlags_PickerHueWheel
				}
				if display_mode == 1 {
					flags |= ImGuiColorEditFlags_NoInputs
				} // Disable all RGB/HSV/Hex displays
				if display_mode == 2 {
					flags |= ImGuiColorEditFlags_DisplayRGB
				} // Override display mode
				if display_mode == 3 {
					flags |= ImGuiColorEditFlags_DisplayHSV
				}
				if display_mode == 4 {
					flags |= ImGuiColorEditFlags_DisplayHex
				}

				var refcolor *int
				if ref_color {
					refcolor = &ref_color_v.x
				}

				imgui.ColorPicker4("MyColor##4", color, flags, refcolor)

				imgui.Text("Set defaults in code:")
				imgui.SameLine()
				HelpMarker(
					"SetColorEditOptions() is designed to allow you to set boot-time default.\n" +
						"We don't have Push/Pop functions because you can force options on a per-widget basis if needed," +
						"and the user can change non-forced ones with the options menu.\nWe don't have a getter to avoid" +
						"encouraging you to persistently save values that aren't forward-compatible.")
				if imgui.Button("Default: Uint8 + HSV + Hue Bar") {
					imgui.SetColorEditOptions(ImGuiColorEditFlags_Uint8 | ImGuiColorEditFlags_DisplayHSV | ImGuiColorEditFlags_PickerHueBar)
				}
				if imgui.Button("Default: Float + HDR + Hue Wheel") {
					imgui.SetColorEditOptions(ImGuiColorEditFlags_Float | ImGuiColorEditFlags_HDR | ImGuiColorEditFlags_PickerHueWheel)
				}

				// HSV encoded support (to avoid RGB<>HSV round trips and singularities when S==0 or V==0)
				imgui.Spacing()
				imgui.Text("HSV encoded colors")
				imgui.SameLine()
				HelpMarker(
					"By default, colors are given to ColorEdit and ColorPicker in RGB, but ImGuiColorEditFlags_InputHSV" +
						"allows you to store colors as HSV and pass them to ColorEdit and ColorPicker as HSV. This comes with the" +
						"added benefit that you can manipulate hue values with the picker even when saturation or value are zero.")
				imgui.Text("Color widget with InputHSV:")
				imgui.ColorEdit4("HSV shown as RGB##1", color_hsv, ImGuiColorEditFlags_DisplayRGB|ImGuiColorEditFlags_InputHSV|ImGuiColorEditFlags_Float)
				imgui.ColorEdit4("HSV shown as HSV##1", color_hsv, ImGuiColorEditFlags_DisplayHSV|ImGuiColorEditFlags_InputHSV|ImGuiColorEditFlags_Float)
				imgui.DragFloat4("Raw HSV values", color_hsv, 0.01, 0.0, 1.0)

				imgui.TreePop()
			}

			if imgui.TreeNode("Drag/Slider Flags") {
				// Demonstrate using advanced flags for DragXXX and SliderXXX functions. Note that the flags are the same!
				imgui.CheckboxFlags("ImGuiSliderFlags_AlwaysClamp", &sliderflags, ImGuiSliderFlags_AlwaysClamp)
				imgui.SameLine()
				HelpMarker("Always clamp value to min/max bounds (if any) when input manually with CTRL+Click.")
				imgui.CheckboxFlags("ImGuiSliderFlags_Logarithmic", &sliderflags, ImGuiSliderFlags_Logarithmic)
				imgui.SameLine()
				HelpMarker("Enable logarithmic editing (more precision for small values).")
				imgui.CheckboxFlags("ImGuiSliderFlags_NoRoundToFormat", &sliderflags, ImGuiSliderFlags_NoRoundToFormat)
				imgui.SameLine()
				HelpMarker("Disable rounding underlying value to match precision of the format string (e.g. %.3 values are rounded to those 3 digits).")
				imgui.CheckboxFlags("ImGuiSliderFlags_NoInput", &sliderflags, ImGuiSliderFlags_NoInput)
				imgui.SameLine()
				HelpMarker("Disable CTRL+Click or Enter key allowing to input text directly into the widget.")

				// Drags
				imgui.Text("Underlying float value: %f", drag_f)
				imgui.DragFloat("DragFloat (0 . 1)", &drag_f, 0.005, 0.0, 1.0, "%.3", flags)
				imgui.DragFloat("DragFloat (0 . +inf)", &drag_f, 0.005, 0.0, FLT_MAX, "%.3", flags)
				imgui.DragFloat("DragFloat (-inf . 1)", &drag_f, 0.005, -FLT_MAX, 1.0, "%.3", flags)
				imgui.DragFloat("DragFloat (-inf . +inf)", &drag_f, 0.005, -FLT_MAX, +FLT_MAX, "%.3", flags)
				imgui.DragInt("DragInt (0 . 100)", &drag_i, 0.5, 0, 100, "%d", flags)

				// Sliders

				imgui.Text("Underlying float value: %f", slider_f)
				imgui.SliderFloat("SliderFloat (0 . 1)", &slider_f, 0.0, 1.0, "%.3", flags)
				imgui.SliderInt("SliderInt (0 . 100)", &slider_i, 0, 100, "%d", flags)

				imgui.TreePop()
			}

			if imgui.TreeNode("Range Widgets") {
				imgui.DragFloatRange2("range float", &begin, &end, 0.25, 0.0, 100.0, "Min: %.1 %%", "Max: %.1 %%", ImGuiSliderFlags_AlwaysClamp)
				imgui.DragIntRange2("range int", &begin_i, &end_i, 5, 0, 1000, "Min: %d units", "Max: %d units")
				imgui.DragIntRange2("range int (no bounds)", &begin_i, &end_i, 5, 0, 0, "Min: %d units", "Max: %d units")
				imgui.TreePop()
			}

			if imgui.TreeNode("Data Types") {
				// DragScalar/InputScalar/SliderScalar functions allow various data types

				const drag_speed = 0.2
				imgui.Text("Drags:")
				imgui.Checkbox("Clamp integers to 0..50", &drag_clamp)
				imgui.SameLine()
				HelpMarker(
					"As with every widgets in dear imgui, we never modify values unless there is a user interaction.\n" +
						"You can override the clamping limits by using CTRL+Click to input a value.")
				imgui.DragScalar("drag s8", ImGuiDataType_S8, &s8_v, drag_speed, 0, 50)
				imgui.DragScalar("drag u8", ImGuiDataType_U8, &u8_v, drag_speed, 0, 50, "%u ms")
				imgui.DragScalar("drag s16", ImGuiDataType_S16, &s16_v, drag_speed, 0, 50)
				imgui.DragScalar("drag u16", ImGuiDataType_U16, &u16_v, drag_speed, 0, 50, "%u ms")
				imgui.DragScalar("drag s32", ImGuiDataType_S32, &s32_v, drag_speed, 0, 50)
				imgui.DragScalar("drag u32", ImGuiDataType_U32, &u32_v, drag_speed, 0, 50, "%u ms")
				imgui.DragScalar("drag s64", ImGuiDataType_S64, &s64_v, drag_speed, 0, 50)
				imgui.DragScalar("drag u64", ImGuiDataType_U64, &u64_v, drag_speed, 0, 50)
				imgui.DragScalar("drag float", ImGuiDataType_Float, &f32_v, 0.005, &f32_zero, &f32_one, "%f")
				imgui.DragScalar("drag float log", ImGuiDataType_Float, &f32_v, 0.005, &f32_zero, &f32_one, "%f", ImGuiSliderFlags_Logarithmic)
				imgui.DragScalar("drag double", ImGuiDataType_Double, &f64_v, 0.0005, &f64_zero, NULL, "%.10 grams")
				imgui.DragScalar("drag double log", ImGuiDataType_Double, &f64_v, 0.0005, &f64_zero, &f64_one, "0 < %.10 < 1", ImGuiSliderFlags_Logarithmic)

				imgui.Text("Sliders")
				imgui.SliderScalar("slider s8 full", ImGuiDataType_S8, &s8_v, &s8_min, &s8_max, "%d")
				imgui.SliderScalar("slider u8 full", ImGuiDataType_U8, &u8_v, &u8_min, &u8_max, "%u")
				imgui.SliderScalar("slider s16 full", ImGuiDataType_S16, &s16_v, &s16_min, &s16_max, "%d")
				imgui.SliderScalar("slider u16 full", ImGuiDataType_U16, &u16_v, &u16_min, &u16_max, "%u")
				imgui.SliderScalar("slider s32 low", ImGuiDataType_S32, &s32_v, &s32_zero, &s32_fifty, "%d")
				imgui.SliderScalar("slider s32 high", ImGuiDataType_S32, &s32_v, &s32_hi_a, &s32_hi_b, "%d")
				imgui.SliderScalar("slider s32 full", ImGuiDataType_S32, &s32_v, &s32_min, &s32_max, "%d")
				imgui.SliderScalar("slider u32 low", ImGuiDataType_U32, &u32_v, &u32_zero, &u32_fifty, "%u")
				imgui.SliderScalar("slider u32 high", ImGuiDataType_U32, &u32_v, &u32_hi_a, &u32_hi_b, "%u")
				imgui.SliderScalar("slider u32 full", ImGuiDataType_U32, &u32_v, &u32_min, &u32_max, "%u")
				imgui.SliderScalar("slider s64 low", ImGuiDataType_S64, &s64_v, &s64_zero, &s64_fifty, "%"+IM_PRId64)
				imgui.SliderScalar("slider s64 high", ImGuiDataType_S64, &s64_v, &s64_hi_a, &s64_hi_b, "%"+IM_PRId64)
				imgui.SliderScalar("slider s64 full", ImGuiDataType_S64, &s64_v, &s64_min, &s64_max, "%"+IM_PRId64)
				imgui.SliderScalar("slider u64 low", ImGuiDataType_U64, &u64_v, &u64_zero, &u64_fifty, "%"+IM_PRIu64+" ms")
				imgui.SliderScalar("slider u64 high", ImGuiDataType_U64, &u64_v, &u64_hi_a, &u64_hi_b, "%"+IM_PRIu64+" ms")
				imgui.SliderScalar("slider u64 full", ImGuiDataType_U64, &u64_v, &u64_min, &u64_max, "%"+IM_PRIu64+" ms")
				imgui.SliderScalar("slider float low", ImGuiDataType_Float, &f32_v, &f32_zero, &f32_one)
				imgui.SliderScalar("slider float low log", ImGuiDataType_Float, &f32_v, &f32_zero, &f32_one, "%.10", ImGuiSliderFlags_Logarithmic)
				imgui.SliderScalar("slider float high", ImGuiDataType_Float, &f32_v, &f32_lo_a, &f32_hi_a, "%e")
				imgui.SliderScalar("slider double low", ImGuiDataType_Double, &f64_v, &f64_zero, &f64_one, "%.10 grams")
				imgui.SliderScalar("slider double low log", ImGuiDataType_Double, &f64_v, &f64_zero, &f64_one, "%.10", ImGuiSliderFlags_Logarithmic)
				imgui.SliderScalar("slider double high", ImGuiDataType_Double, &f64_v, &f64_lo_a, &f64_hi_a, "%e grams")

				imgui.Text("Sliders (reverse)")
				imgui.SliderScalar("slider s8 reverse", ImGuiDataType_S8, &s8_v, &s8_max, &s8_min, "%d")
				imgui.SliderScalar("slider u8 reverse", ImGuiDataType_U8, &u8_v, &u8_max, &u8_min, "%u")
				imgui.SliderScalar("slider s32 reverse", ImGuiDataType_S32, &s32_v, &s32_fifty, &s32_zero, "%d")
				imgui.SliderScalar("slider u32 reverse", ImGuiDataType_U32, &u32_v, &u32_fifty, &u32_zero, "%u")
				imgui.SliderScalar("slider s64 reverse", ImGuiDataType_S64, &s64_v, &s64_fifty, &s64_zero, "%"+IM_PRId64)
				imgui.SliderScalar("slider u64 reverse", ImGuiDataType_U64, &u64_v, &u64_fifty, &u64_zero, "%"+IM_PRIu64+" ms")

				imgui.Text("Inputs")
				imgui.Checkbox("Show step buttons", &inputs_step)
				imgui.InputScalar("input s8", ImGuiDataType_S8, &s8_v, 1, NULL, "%d")
				imgui.InputScalar("input u8", ImGuiDataType_U8, &u8_v, 1, NULL, "%u")
				imgui.InputScalar("input s16", ImGuiDataType_S16, &s16_v, 1, NULL, "%d")
				imgui.InputScalar("input u16", ImGuiDataType_U16, &u16_v, 1, NULL, "%u")
				imgui.InputScalar("input s32", ImGuiDataType_S32, &s32_v, 1, NULL, "%d")
				imgui.InputScalar("input s32 hex", ImGuiDataType_S32, &s32_v, 1, NULL, "%08X", ImGuiInputTextFlags_CharsHexadecimal)
				imgui.InputScalar("input u32", ImGuiDataType_U32, &u32_v, 1, NULL, "%u")
				imgui.InputScalar("input u32 hex", ImGuiDataType_U32, &u32_v, 1, NULL, "%08X", ImGuiInputTextFlags_CharsHexadecimal)
				imgui.InputScalar("input s64", ImGuiDataType_S64, &s64_v, 1)
				imgui.InputScalar("input u64", ImGuiDataType_U64, &u64_v, 1)
				imgui.InputScalar("input float", ImGuiDataType_Float, &f32_v, 1)
				imgui.InputScalar("input double", ImGuiDataType_Double, &f64_v, 1)

				imgui.TreePop()
			}

			if imgui.TreeNode("Multi-component Widgets") {

				imgui.InputFloat2("input float2", MCvec4)
				imgui.DragFloat2("drag float2", MCvec4, 0.01, 0.0, 1.0)
				imgui.SliderFloat2("slider float2", MCvec4, 0.0, 1.0)
				imgui.InputInt2("input int2", MCvec4i)
				imgui.DragInt2("drag int2", MCvec4i, 1, 0, 255)
				imgui.SliderInt2("slider int2", MCvec4i, 0, 255)
				imgui.Spacing()

				imgui.InputFloat3("input float3", MCvec4)
				imgui.DragFloat3("drag float3", MCvec4, 0.01, 0.0, 1.0)
				imgui.SliderFloat3("slider float3", MCvec4, 0.0, 1.0)
				imgui.InputInt3("input int3", MCvec4i)
				imgui.DragInt3("drag int3", MCvec4i, 1, 0, 255)
				imgui.SliderInt3("slider int3", MCvec4i, 0, 255)
				imgui.Spacing()

				imgui.InputFloat4("input float4", MCvec4)
				imgui.DragFloat4("drag float4", MCvec4, 0.01, 0.0, 1.0)
				imgui.SliderFloat4("slider float4", MCvec4, 0.0, 1.0)
				imgui.InputInt4("input int4", MCvec4i)
				imgui.DragInt4("drag int4", MCvec4i, 1, 0, 255)
				imgui.SliderInt4("slider int4", MCvec4i, 0, 255)

				imgui.TreePop()
			}

			if imgui.TreeNode("Vertical Sliders") {
				const float spacing = 4
				imgui.PushStyleVar(ImGuiStyleVar_ItemSpacing, ImVec2(spacing, spacing))

				imgui.VSliderInt("##int", ImVec2(18, 160), &int_value, 0, 5)
				imgui.SameLine()

				imgui.PushID("set1")
				for i := 0; i < 7; i++ {
					if i > 0 {
						imgui.SameLine()
					}
					imgui.PushID(i)
					imgui.PushStyleColor(ImGuiCol_FrameBg, ImColor.HSV(i/7.0, 0.5, 0.5))
					imgui.PushStyleColor(ImGuiCol_FrameBgHovered, ImColor.HSV(i/7.0, 0.6, 0.5))
					imgui.PushStyleColor(ImGuiCol_FrameBgActive, ImColor.HSV(i/7.0, 0.7, 0.5))
					imgui.PushStyleColor(ImGuiCol_SliderGrab, ImColor.HSV(i/7.0, 0.9, 0.9))
					imgui.VSliderFloat("##v", ImVec2(18, 160), &vertical_values[i], 0.0, 1.0, "")
					if imgui.IsItemActive() || imgui.IsItemHovered() {
						imgui.SetTooltip("%.3", vertical_values[i])
					}
					imgui.PopStyleColor(4)
					imgui.PopID()
				}
				imgui.PopID()

				imgui.SameLine()
				imgui.PushID("set2")

				const int rows = 3
				var small_slider_size = ImVec2{18, (float)(int)((160.0 - (rows-1)*spacing) / rows)}
				for nx := 0; nx < 4; nx++ {
					if nx > 0 {
						imgui.SameLine()
					}
					imgui.BeginGroup()
					for ny := 0; ny < rows; ny++ {
						imgui.PushID(nx*rows + ny)
						imgui.VSliderFloat("##v", small_slider_size, &vertical_values[nx], 0.0, 1.0, "")
						if imgui.IsItemActive() || imgui.IsItemHovered() {
							imgui.SetTooltip("%.3", vertical_values[nx])
						}
						imgui.PopID()
					}
					imgui.EndGroup()
				}
				imgui.PopID()

				imgui.SameLine()
				imgui.PushID("set3")
				for i := 0; i < 4; i++ {
					if i > 0 {
						imgui.SameLine()
					}

					imgui.PushID(i)
					imgui.PushStyleVar(ImGuiStyleVar_GrabMinSize, 40)
					imgui.VSliderFloat("##v", ImVec2(40, 160), &values[i], 0.0, 1.0, "%.2\nsec")
					imgui.PopStyleVar()
					imgui.PopID()
				}
				imgui.PopID()
				imgui.PopStyleVar()
				imgui.TreePop()
			}

			if imgui.TreeNode("Drag and Drop") {
				if imgui.TreeNode("Drag and drop in standard widgets") {
					// ColorEdit widgets automatically act as drag source and drag target.
					// They are using standardized payload strings IMGUI_PAYLOAD_TYPE_COLOR_3 and IMGUI_PAYLOAD_TYPE_COLOR_4
					// to allow your own widgets to use colors in their drag and drop interaction.
					// Also see 'Demo.Widgets.Color/Picker Widgets.Palette' demo.
					HelpMarker("You can drag from the color squares.")
					imgui.ColorEdit3("color 1", col1)
					imgui.ColorEdit4("color 2", col2)
					imgui.TreePop()
				}

				if imgui.TreeNode("Drag and drop to copy/swap items") {
					const (
						Mode_Copy = iota
						Mode_Move
						Mode_Swap
					)
					if imgui.RadioButton("Copy", mode == Mode_Copy) {
						mode = Mode_Copy
					}
					imgui.SameLine()
					if imgui.RadioButton("Move", mode == Mode_Move) {
						mode = Mode_Move
					}
					imgui.SameLine()
					if imgui.RadioButton("Swap", mode == Mode_Swap) {
						mode = Mode_Swap
					}
					var names = [...]string{
						"Bobby", "Beatrice", "Betty",
						"Brianna", "Barry", "Bernard",
						"Bibi", "Blaine", "Bryn",
					}
					for n := range names {
						imgui.PushID(n)
						if (n % 3) != 0 {
							imgui.SameLine()
						}
						imgui.Button(names[n], ImVec2(60, 60))

						// Our buttons are both drag sources and drag targets here!
						if imgui.BeginDragDropSource(ImGuiDragDropFlags_None) {
							// Set payload to carry the index of our item (could be anything)
							imgui.SetDragDropPayload("DND_DEMO_CELL", &n, sizeof(int))

							// Display preview (could be anything, e.g. when dragging an image we could decide to display
							// the filename and a small preview of the image, etc.)
							if mode == Mode_Copy {
								imgui.Text("Copy %s", names[n])
							}
							if mode == Mode_Move {
								imgui.Text("Move %s", names[n])
							}
							if mode == Mode_Swap {
								imgui.Text("Swap %s", names[n])
							}
							imgui.EndDragDropSource()
						}
						if imgui.BeginDragDropTarget() {
							//TODO PORTING again not sure how this works
							/* if (const ImGuiPayload* payload = imgui.AcceptDragDropPayload("DND_DEMO_CELL"))
							   {
							       IM_ASSERT(payload.DataSize == sizeof(int));
							       int payload_n = *(const int*)payload.Data;
							       if (mode == Mode_Copy)
							       {
							           names[n] = names[payload_n];
							       }
							       if (mode == Mode_Move)
							       {
							           names[n] = names[payload_n];
							           names[payload_n] = "";
							       }
							       if (mode == Mode_Swap)
							       {
							           const char* tmp = names[n];
							           names[n] = names[payload_n];
							           names[payload_n] = tmp;
							       }
							   }*/
							imgui.EndDragDropTarget()
						}
						imgui.PopID()
					}
					imgui.TreePop()
				}

				if imgui.TreeNode("Drag to reorder items (simple)") {
					// Simple reordering
					HelpMarker(
						"We don't use the drag and drop api at all here! " +
							"Instead we query when the item is held but not hovered, and order items accordingly.")
					var item_names = [...]string{"Item One", "Item Two", "Item Three", "Item Four", "Item Five"}
					for n := range item_names {
						var item = item_names[n]
						imgui.Selectable(item)

						if imgui.IsItemActive() && !imgui.IsItemHovered() {
							var n_next = n
							if imgui.GetMouseDragDelta(0).y < 0 {
								n -= 1
							} else {
								n += 1
							}

							if n_next >= 0 && n_next < IM_ARRAYSIZE(item_names) {
								item_names[n] = item_names[n_next]
								item_names[n_next] = item
								imgui.ResetMouseDragDelta()
							}
						}
					}
					imgui.TreePop()
				}

				imgui.TreePop()
			}

			if imgui.TreeNode("Querying Status (Edited/Active/Hovered etc.)") {
				// Select an item type
				var item_names = [...]string{
					"Text", "Button", "Button (w/ repeat)", "Checkbox", "SliderFloat", "InputText", "InputFloat",
					"InputFloat3", "ColorEdit4", "Selectable", "MenuItem", "TreeNode", "TreeNode (w/ double-click)", "Combo", "ListBox",
				}
				imgui.Combo("Item Type", &item_type, item_names, IM_ARRAYSIZE(item_names), IM_ARRAYSIZE(item_names))
				imgui.SameLine()
				HelpMarker("Testing how various types of items are interacting with the IsItemXXX functions. Note that the bool return value of most ImGui function is generally equivalent to calling imgui.IsItemHovered().")
				imgui.Checkbox("Item Disabled", &item_disabled)

				// Submit selected item item so we can query their status in the code following it.
				var ret = false
				if item_disabled {
					imgui.BeginDisabled(true)
				}
				if item_type == 0 {
					imgui.Text("ITEM: Text")
				} // Testing text items with no identifier/interaction
				if item_type == 1 {
					ret = imgui.Button("ITEM: Button")
				} // Testing button
				if item_type == 2 {
					imgui.PushButtonRepeat(true)
					ret = imgui.Button("ITEM: Button")
					imgui.PopButtonRepeat()
				} // Testing button (with repeater)
				if item_type == 3 {
					ret = imgui.Checkbox("ITEM: Checkbox", &b)
				} // Testing checkbox
				if item_type == 4 {
					ret = imgui.SliderFloat("ITEM: SliderFloat", &col4[0], 0.0, 1.0)
				} // Testing basic item
				if item_type == 5 {
					ret = imgui.InputText("ITEM: InputText", &str[0], IM_ARRAYSIZE(str))
				} // Testing input text (which handles tabbing)
				if item_type == 6 {
					ret = imgui.InputFloat("ITEM: InputFloat", col4, 1.0)
				} // Testing +/- buttons on scalar input
				if item_type == 7 {
					ret = imgui.InputFloat3("ITEM: InputFloat3", col4)
				} // Testing multi-component items (IsItemXXX flags are reported merged)
				if item_type == 8 {
					ret = imgui.ColorEdit4("ITEM: ColorEdit4", col4)
				} // Testing multi-component items (IsItemXXX flags are reported merged)
				if item_type == 9 {
					ret = imgui.Selectable("ITEM: Selectable")
				} // Testing selectable item
				if item_type == 10 {
					ret = imgui.MenuItem("ITEM: MenuItem")
				} // Testing menu item (they use ImGuiButtonFlags_PressedOnRelease button policy)
				if item_type == 11 {
					ret = imgui.TreeNode("ITEM: TreeNode")
					if ret {
						imgui.TreePop()
					} // Testing tree node
				}
				if item_type == 12 {
					ret = imgui.TreeNodeEx("ITEM: TreeNode w/ ImGuiTreeNodeFlags_OpenOnDoubleClick", ImGuiTreeNodeFlags_OpenOnDoubleClick|ImGuiTreeNodeFlags_NoTreePushOnOpen)
				} // Testing tree node with ImGuiButtonFlags_PressedOnDoubleClick button policy.
				if item_type == 13 {
					var items = [...]string{"Apple", "Banana", "Cherry", "Kiwi"}
					ret = imgui.Combo("ITEM: Combo", &current, items, IM_ARRAYSIZE(items))
				}
				if item_type == 14 {
					var items = [...]string{"Apple", "Banana", "Cherry", "Kiwi"}
					ret = imgui.ListBox("ITEM: ListBox", &current2, items, IM_ARRAYSIZE(items), IM_ARRAYSIZE(items))
				}

				// Display the values of IsItemHovered() and other common item state functions.
				// Note that the ImGuiHoveredFlags_XXX flags can be combined.
				// Because BulletText is an item itself and that would affect the output of IsItemXXX functions,
				// we query every state in a single call to avoid storing them and to simplify the code.
				imgui.BulletText(
					"Return value = %d\n"+
						"IsItemFocused() = %d\n"+
						"IsItemHovered() = %d\n"+
						"IsItemHovered(_AllowWhenBlockedByPopup) = %d\n"+
						"IsItemHovered(_AllowWhenBlockedByActiveItem) = %d\n"+
						"IsItemHovered(_AllowWhenOverlapped) = %d\n"+
						"IsItemHovered(_AllowWhenDisabled) = %d\n"+
						"IsItemHovered(_RectOnly) = %d\n"+
						"IsItemActive() = %d\n"+
						"IsItemEdited() = %d\n"+
						"IsItemActivated() = %d\n"+
						"IsItemDeactivated() = %d\n"+
						"IsItemDeactivatedAfterEdit() = %d\n"+
						"IsItemVisible() = %d\n"+
						"IsItemClicked() = %d\n"+
						"IsItemToggledOpen() = %d\n"+
						"GetItemRectMin() = (%.1, %.1)\n"+
						"GetItemRectMax() = (%.1, %.1)\n"+
						"GetItemRectSize() = (%.1, %.1)",
					ret,
					imgui.IsItemFocused(),
					imgui.IsItemHovered(),
					imgui.IsItemHovered(ImGuiHoveredFlags_AllowWhenBlockedByPopup),
					imgui.IsItemHovered(ImGuiHoveredFlags_AllowWhenBlockedByActiveItem),
					imgui.IsItemHovered(ImGuiHoveredFlags_AllowWhenOverlapped),
					imgui.IsItemHovered(ImGuiHoveredFlags_AllowWhenDisabled),
					imgui.IsItemHovered(ImGuiHoveredFlags_RectOnly),
					imgui.IsItemActive(),
					imgui.IsItemEdited(),
					imgui.IsItemActivated(),
					imgui.IsItemDeactivated(),
					imgui.IsItemDeactivatedAfterEdit(),
					imgui.IsItemVisible(),
					imgui.IsItemClicked(),
					imgui.IsItemToggledOpen(),
					imgui.GetItemRectMin().x, imgui.GetItemRectMin().y,
					imgui.GetItemRectMax().x, imgui.GetItemRectMax().y,
					imgui.GetItemRectSize().x, imgui.GetItemRectSize().y,
				)

				if item_disabled {
					imgui.EndDisabled()
				}

				imgui.Checkbox("Embed everything inside a child window (for additional testing)", &embed_all_inside_a_child_window)
				if embed_all_inside_a_child_window {
					imgui.BeginChild("outer_child", ImVec2(0, imgui.GetFontSize()*20.0), true)
				}

				// Testing IsWindowFocused() function with its various flags.
				// Note that the ImGuiFocusedFlags_XXX flags can be combined.
				imgui.BulletText(
					"IsWindowFocused() = %d\n"+
						"IsWindowFocused(_ChildWindows) = %d\n"+
						"IsWindowFocused(_ChildWindows|_RootWindow) = %d\n"+
						"IsWindowFocused(_RootWindow) = %d\n"+
						"IsWindowFocused(_AnyWindow) = %d\n",
					imgui.IsWindowFocused(),
					imgui.IsWindowFocused(ImGuiFocusedFlags_ChildWindows),
					imgui.IsWindowFocused(ImGuiFocusedFlags_ChildWindows|ImGuiFocusedFlags_RootWindow),
					imgui.IsWindowFocused(ImGuiFocusedFlags_RootWindow),
					imgui.IsWindowFocused(ImGuiFocusedFlags_AnyWindow))

				// Testing IsWindowHovered() function with its various flags.
				// Note that the ImGuiHoveredFlags_XXX flags can be combined.
				imgui.BulletText(
					"IsWindowHovered() = %d\n"+
						"IsWindowHovered(_AllowWhenBlockedByPopup) = %d\n"+
						"IsWindowHovered(_AllowWhenBlockedByActiveItem) = %d\n"+
						"IsWindowHovered(_ChildWindows) = %d\n"+
						"IsWindowHovered(_ChildWindows|_RootWindow) = %d\n"+
						"IsWindowHovered(_ChildWindows|_AllowWhenBlockedByPopup) = %d\n"+
						"IsWindowHovered(_RootWindow) = %d\n"+
						"IsWindowHovered(_AnyWindow) = %d\n",
					imgui.IsWindowHovered(),
					imgui.IsWindowHovered(ImGuiHoveredFlags_AllowWhenBlockedByPopup),
					imgui.IsWindowHovered(ImGuiHoveredFlags_AllowWhenBlockedByActiveItem),
					imgui.IsWindowHovered(ImGuiHoveredFlags_ChildWindows),
					imgui.IsWindowHovered(ImGuiHoveredFlags_ChildWindows|ImGuiHoveredFlags_RootWindow),
					imgui.IsWindowHovered(ImGuiHoveredFlags_ChildWindows|ImGuiHoveredFlags_AllowWhenBlockedByPopup),
					imgui.IsWindowHovered(ImGuiHoveredFlags_RootWindow),
					imgui.IsWindowHovered(ImGuiHoveredFlags_AnyWindow))

				imgui.BeginChild("child", ImVec2(0, 50), true)
				imgui.Text("This is another child window for testing the _ChildWindows flag.")
				imgui.EndChild()
				if embed_all_inside_a_child_window {
					imgui.EndChild()
				}

				imgui.InputText("unused", unused_str, IM_ARRAYSIZE(unused_str), ImGuiInputTextFlags_ReadOnly)

				// Calling IsItemHovered() after begin returns the hovered status of the title bar.
				// This is useful in particular if you want to create a context menu associated to the title bar of a window.
				imgui.Checkbox("Hovered/Active tests after Begin() for title bar testing", &test_window)
				if test_window {
					imgui.Begin("Title bar Hovered/Active tests", &test_window)
					if imgui.BeginPopupContextItem() { // <-- This is using IsItemHovered()
						if imgui.MenuItem("Close") {
							test_window = false
						}
						imgui.EndPopup()
					}
					imgui.Text(
						"IsItemHovered() after begin = %d (== is title bar hovered)\n"+
							"IsItemActive() after begin = %d (== is window being clicked/moved)\n",
						imgui.IsItemHovered(), imgui.IsItemActive())
					imgui.End()
				}

				imgui.TreePop()
			}

			// Demonstrate BeginDisabled/EndDisabled using a checkbox located at the bottom of the section (which is a bit odd:
			// logically we'd have this checkbox at the top of the section, but we don't want this feature to steal that space)
			if disable_all {
				imgui.EndDisabled()
			}

			if imgui.TreeNode("Disable block") {
				imgui.Checkbox("Disable entire section above", &disable_all)
				imgui.SameLine()
				HelpMarker("Demonstrate using BeginDisabled()/EndDisabled() across this section.")
				imgui.TreePop()
			}
		}
	}
}
