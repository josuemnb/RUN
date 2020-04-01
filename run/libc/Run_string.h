#pragma once

#include <malloc.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

class Run_string {
   public:
    char *value;
    int len;
    Run_string() {
        value = NULL;
        len = 0;
    }

    Run_string(char *s) {
        len = strlen(s);
        value = (char *)malloc(len);
        strcpy(value, s);
    }

    Run_string(const char *s) {
        Run_string((char*)s);
    }

    Run_string(Run_string *s) {
        value = (char *)malloc(s->len);
        strcpy((char *)value, s->value);
        len = s->len;
    }

    Run_string &operator=(Run_string &other) {
        len = other.len;

        value = (char *)malloc(len);
        strcpy((char *)value, other.value);
        value[len] = '\0';
        return *this;
    }

    Run_string &operator=(char *val) {
        len = strlen(val);

        value = (char *)malloc(len);
        strcpy((char *)value, val);
        return *this;
    }

    Run_string &operator=(const char *val) {
        len = strlen(val);

        value = (char *)malloc(len);
        strcpy((char *)value, val);
        return *this;
    }

    Run_string set(char *val, int size) {
        value = (char *)malloc(size);
        for (int i = 0; i < size; i++)
            value[i] = val[i];
        value[size] = '\0';
        return *this;
    }

    bool method_starts_class_string(Run_string other) {
        if (other.len > len)
            return false;
        for (int i = 0; i < other.len; i++) {
            if (value[i] != other.value[i])
                return false;
        }
        return true;
    }

    bool method_starts_string(char *val) {
        int l = strlen(val);
        if (l > len)
            return false;
        for (int i = 0; i < l; i++) {
            if (value[i] != val[i])
                return false;
        }
        return true;
    }

    bool method_starts_string(const char *val) {
        int l = strlen(val);
        if (l > len)
            return false;
        for (int i = 0; i < l; i++) {
            if (value[i] != val[i])
                return false;
        }
        return true;
    }

    bool operator==(char *val) {
        if (value == val)
            return true;
        int l = strlen(val);
        if (len != l) {
            return false;
        }
        return (strcmp(value, val) == 0);
    }

    bool operator==(const char *val) {
        if (value == val)
            return true;
        int l = strlen(val);
        if (len != l) {
            return false;
        }
        return (strcmp(value, val) == 0);
    }

    bool operator==(Run_string &val) {
        if (len != val.len) {
            return false;
        }
        return (strcmp(value, val.value) == 0);
    }
    bool operator!=(char *val) {
        if (value != val)
            return true;
        int l = strlen(val);
        if (len == l) {
            return false;
        }
        return (strcmp(value, val) != 0);
    }

    bool operator!=(Run_string &val) {
        if (len == val.len) {
            return false;
        }
        return (strcmp(value, val.value) != 0);
    }

    Run_string &operator+=(Run_string &other) {
        value = (char *)realloc((void *)value, len + other.len);
        strcat((char *)value, other.value);
        len += other.len;
        return *this;
    }

    Run_string &operator+=(char *other) {
        int l = strlen(other);
        value = (char *)realloc(value, len + l + 1);
        if (value == NULL) {
            puts("Error");
            exit(-2);
        }
        strcat((char *)value, other);
        len += l;
        return *this;
    }

        Run_string &operator+=(const char *other) {
        int l = strlen(other);
        value = (char *)realloc(value, len + l + 1);
        if (value == NULL) {
            puts("Error");
            exit(-2);
        }
        strcat((char *)value, other);
        len += l;
        return *this;
    }

    Run_string method_append_string_number(char *val, int l) {
        if (len > 0) {
            value = (char *)realloc(value, len + l + 1);
            if (value == NULL) {
                puts("Error");
                exit(-1);
            }
        } else {
            value = (char *)malloc(l + 1);
        }
        for (int i = 0; i < l; i++)
            value[len + i] = val[i];
        value[len + l] = '\0';
        len += l;
        return *this;
    }

    int method_indexOf_string(char *val) {
        int l = strlen(val);
        if (l > len) {
            return -1;
        }
        for (int i = 0; i < len - l + 1; i++) {
            if (value[i] == val[0]) {
                if (l == 1)
                    return i;
                bool found = true;
                for (int j = 1; j < l; j++) {
                    if (value[i + j] != val[j]) {
                        found = false;
                        break;
                    }
                }
                if (found) {
                    return i;
                }
            }
        }
        return -1;
    }

    int method_indexOf_string_number(char *val, int pos) {
        int l = strlen(val);
        if (l + pos > len) {
            return -1;
        }
        for (int i = pos; i < len - l + 1; i++) {
            if (value[i] == val[0]) {
                if (l == 1)
                    return i;
                bool found = true;
                for (int j = 1; j < l; j++) {
                    if (value[i + j] != val[j]) {
                        found = false;
                        break;
                    }
                }
                if (found) {
                    return i;
                }
            }
        }
        return -1;
    }

    Run_string method_substring_number_number(int b, int e) {
        if (len < e) {
            exit(-1);
        }
        char *val = (char *)malloc(e - b + 1);
        strncpy((char *)val, value + b, e - b);
        val[e - b] = '\0';
        Run_string v(val);
        return v;
    }

    char method_at_number(int i) {
        return value[i];
    }

    int method_size() {
        return len;
    }

    ~Run_string() {
        free(value);
    }
};