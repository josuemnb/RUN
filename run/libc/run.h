#pragma once

#include <io.h>
#include <malloc.h>
#include <stdio.h>
#include <stdlib.h>

#ifndef bool
#define bool unsigned char
#endif

#define BOOL(b) b == 0 ? "false" : "true"

#ifndef number
#define number long long
#endif

#ifndef real
#define real double
#endif

#ifndef byte
#define byte char
#endif

#ifndef string
#define string byte *
#endif

#define gcnew(T, size) ({  \
    T *_new = new T[size]; \
    _new;                  \
})
// #include "Run_string.h"

// Run_string func_string_number(int n) {
//     char v[40];
//     sprintf(v,"%lld",n);
//     Run_string s(v);
//     return s;
// }

number func_number_string(string s) {
    return atol(s);
}

void terminate(const char *msg) {
    puts(msg);
    exit(-1);
}

// #include "includes.h"
// #include "imports.h"
// #include "types.h"