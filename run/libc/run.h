#pragma once

#include <stdio.h>
#include <stdlib.h>
#include <malloc.h>

#ifndef bool
#define bool unsigned char
#endif

#define BOOL(b) b==0?"false":"true"

#ifndef number
#define number int
#endif

#ifndef real
#define real double
#endif

#include "class_string.h"
#include "types.h"
#include "collections.h"