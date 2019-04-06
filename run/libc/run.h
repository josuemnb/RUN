#pragma once

#include <stdio.h>
#include <stdlib.h>
#include <malloc.h>
#include <io.h>

#ifndef RUN_bool
#define RUN_bool unsigned char
#endif

#define BOOL(b) b == 0 ? "false" : "true"

#ifndef RUN_number
#define RUN_number long long
#endif

#ifndef RUN_real
#define RUN_real double
#endif

#include "class_string.h"

class_string func_string_number(int n) {
    char v[40];
    sprintf(v,"%lld",n);
    class_string s(v);
    return s;
}

RUN_number func_number_string(class_string s) {
    return atol(s.value);
}

#include "types.h"
#include "collections.h"