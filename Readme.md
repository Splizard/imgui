# GOLANG IMGUI PORTING EFFORTS
License will be MIT

A pure Go port of imgui as of github.com/ocornut/imgui commit 5ee40c8d34bea3009cf462ec963225bd22067e5e 
*(the IMGUI files from this commit are included in the repository)*

I made a start on this work but I am not so familiar with C++ and 
I don't have the time to complete this. You are 
welcome to use this to build upon or to start over, it is up to you. 

It is important to me to have an Imgui package in pure Go, without any C code.
Ideally I would like the entire functionality of Imgui ported, including all
current widgets and features.

The end goal/deliverable is to have the complete IMGUI debug window running 
natively in Go. I started a on this but my attempt is largely incomplete...

## Existing attempt
In the stb directory is my attempt to port the imstb_truetype.h and imstb_rectpack.h files from the imgui source.

    stb/stbtt = imstb_truetype.h 
    stb/stbrp = imstb_rectpack.h 

the stbtt package is incomplete, the header file in the `h` subdirectory contains
the remaining code needing porting, I have deleted the code that has already been ported.

the stbrp package is complete 'in theory' but I have not tested it and it may not work (I may have made mistakes).

The example directory contains a backend for glfw/opengl3 and a compilable
program, however it won't work until more of the imgui source is ported.
The idea is that it panics on any unimplemented function.

`go get && go mod download && go build` with Go 1.16+ should build it fine.

The first thing that Imgui needs to do is load the embedded font, and it seems 
like the font handling of the library is quite complex. That's about as far as
I made it. 

**NOTE**
It may not be worth it to reuse the code from my attempt, so please start from scratch if you professionally believe that will be easier.

## Helpful tips
Whilst attempting to port the `stbtt` and `stbrp` packages, I stumbled upon a few
roadblocks, and I have included a list of things to look out for. 
You may or may not find them useful. 

1. Go has a stronger type system then C++, therefore any integer conversions need to be explicit

        var x uint16
        var y int32
        y = int32(x) //annoying but required for porting

2. A lot of the code I have seen, passes arrays by pointer, in Go
this can generally be replaced by a slice.

        STBTT_DEF int stbtt_GetNumberOfFonts(const unsigned char *data);
        //becomes
        func GetNumberOfFonts(data []byte) int {

3. Sometimes these C++ array pointers are incremented, I find this to be a very    
strange pattern but it can be replicated in Go by using a slice operation.
Keep in mind that the solution below only works for pointer increments, pointer decrements require a different approach (passing an additional index offset).

        //C++
        stbtt_uint8 *points;
        stbtt_uint8 flags = *points++; 

        //Go 
        var points []byte
        points = points[1:] //slice operation, moves the pointer forward by one byte
        var flags byte = points[0] //get the first byte of the slice
        

4. Go doesn't have struct/array constants, but you can just use a variable to 
hold the value.

        //C++
        const ImVec2 zero = ImVec2(0,0);
        //Go
        var zero = ImVec2{0,0}

5. Go doesn't have `static`, so the porting of the imgui debug window is
painful, perhaps it can be broken up into seperate functions? I don't know.

6. If a C++ function uses a *void pointer, this can be replaced with a Go
interface{} type, any value can be assigned to it.

        //C++
        void *userdata;
        //Go
        var userdata interface{}

7. Imgui appears to support callbacks of some kind, they are normally
passed as a struct, I am not too sure how this works in C++, but Go
has function types and closures that can be used instead.

        //C++
        typedef void    (*ImGuiSizeCallback)(ImGuiSizeCallbackData* data)
        //Go
        type ImGuiSizeCallback func(data *ImGuiSizeCallbackData)


8. Go doesn't have ternary operator.

        //C++
        int x = a ? b : c;
        //Go
        var x int32
        if a {
            x = b
        } else {
            x = c
        }

9. Go is garbage collected, so the ImGui/STB memory management can
be removed. Can just use builtin `new` and `make` for allocations.

10. I found it useful to move cpp files into the different packages and go through them function-by-function and delete them as I go to track porting progress. It's nice to see the lines of code in the file you are working on slowly go down as you
port them.

## Outcome
What's important is having a working port so it's probably best not to be to
too clever and to port 1:1 where possible, in my attempt, I have tried to organise
things into different folders, I may be trying to be a bit to clever here and
I believe it will be easier to do a simple 1:1 port changing as little as possible to get the library to work.

Don't hesitate to get in touch with me if there are any issues, or you have any questions. I have extensive experience in Go, so if you run into anything unusual
I can probably help.