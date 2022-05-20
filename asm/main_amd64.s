#include "textflag.h"

// 16: frame size of this function, it's the size of local variables and args of sayHello call (in this function is a string, we define it in go source file)
// 0: args size for this function (we don't have args so it's 0)
TEXT ·Print(SB), $16-0
    MOVQ ·name+0(SB), AX; MOVQ AX, 0(SP) // data part of StringHeader
    MOVQ ·name+8(SB), BX; MOVQ BX, 8(SP) // len part of StringHeader
    CALL ·sayHello(SB)
    RET

TEXT ·SyscallDarwin(SB), NOSPLIT, $0
    MOVQ $(0x2000000+4), AX // syscall NO. (4: write syscall in darwin)
    MOVQ fd+0(FP), DI // 1st args, fd
    MOVQ msg_data+8(FP), SI // 2end args, data part of StringHeader
    MOVQ msg_len+16(FP), DX // 3rd args, len part of StringHeader
    SYSCALL
    MOVQ AX, ret+0(FP) // return result
    RET

TEXT ·getg(SB), NOSPLIT, $32-16
    // get runtime.g
    MOVQ (TLS), AX
    // get runtime.g type from global variable 'type·runtime·g'
    MOVQ $type·runtime·g(SB), BX

    // 在 go1.12 之后下面的代码不再能编译成功，因为引入了 ABIInternal 调用规范
    // convert (*g) to interface{}
    MOVQ AX, 8(SP)
    MOVQ BX, 0(SP)
    CALL runtime·convT2E(SB)
    MOVQ 16(SP), AX // return value
    MOVQ 24(SP), BX

    // return interface{}
    MOVQ AX, ret+0(FP)
    MOVQ BX, ret+8(FP)
    RET
