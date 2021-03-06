#include "textflag.h"

// define a variable
GLOBL ·Id(SB),NOPTR,$8

// initialize the content of variable
DATA ·Id+0(SB)/1,$0x37
DATA ·Id+1(SB)/1,$0x25
DATA ·Id+2(SB)/1,$0x00
DATA ·Id+3(SB)/1,$0x00
DATA ·Id+4(SB)/1,$0x00
DATA ·Id+5(SB)/1,$0x00
DATA ·Id+6(SB)/1,$0x00
DATA ·Id+7(SB)/1,$0x00

// define a string
GLOBL ·Name(SB),NOPTR,$24
DATA ·Name+0(SB)/8, $·Name+16(SB)
DATA ·Name+8(SB)/8, $12
DATA ·Name+16(SB)/12, $"hello, world"

// define a function which call C function in assemably
// passing arguments and return result by registers, which is C stdcall rules
TEXT ·CallCFunc(SB), NOSPLIT, $32-0
    MOVQ fn+0(FP), AX // cfunc
    MOVQ a+8(FP), DI // a
    MOVQ b+16(FP), SI // b

    CALL AX
    MOVQ AX, ret+24(FP)
    RET
