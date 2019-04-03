#pragma once

#include <stdio.h>
#include <malloc.h>
#include <string.h>

class class_string
{
  public:
    char *value;
    int len;
    class_string()
    {
        value = NULL;
        len = 0;
    }

    class_string(char *s)
    {        
        len = strlen(s);
        value = (char*)malloc(len);
        strcpy(value,s);
    }

    class_string(class_string *s)
    {
        value = (char *)malloc(s->len);
        strcpy((char *)value, s->value);
        len = s->len;
    }

    class_string &operator=(class_string &other)
    {
        len = other.len;

        value = (char *)malloc(len);
        strcpy((char *)value, other.value);
        value[len] = '\0';
        return *this;
    }

    class_string &operator=(char *val)
    {
        len = strlen(val);

        value = (char *)malloc(len);
        strcpy((char *)value, val);
        return *this;
    }
    class_string set(char *val, int size)
    {
        value = (char *)malloc(size);
        for (int i = 0; i < size; i++)
            value[i] = val[i];
        value[size] = '\0';
        return *this;
    }

    bool operator==(char *val)
    {
        if (value == val)
            return true;
        int l = strlen(val);
        if (len != l)
        {
            return false;
        }
        return (strcmp(value, val) == 0);
    }
    bool operator==(class_string &val)
    {
        if (len != val.len)
        {
            return false;
        }
        return (strcmp(value, val.value) == 0);
    }
    bool operator!=(char *val)
    {
        if (value != val)
            return true;
        int l = strlen(val);
        if (len == l)
        {
            return false;
        }
        return (strcmp(value, val) != 0);
    }

    bool operator!=(class_string &val)
    {
        if (len == val.len)
        {
            return false;
        }
        return (strcmp(value, val.value) != 0);
    }

    class_string &operator+=(class_string &other)
    {
        value = (char *)realloc((void *)value, len + other.len);
        strcat((char *)value, other.value);
        len += other.len;
        return *this;
    }

    class_string &operator+=(char *other)
    {
        int l = strlen(other);
        value = (char *)realloc(value, len + l+1);
        if(value==NULL) {
            puts("Error");
            exit(-2);
        }
        strcat((char *)value, other);
        len += l;
        return *this;
    }

    class_string method_append_string_number(char *val, int l)
    {
        if (len > 0)
        {
            value = (char *)realloc(value, len + l + 1);
            if (value == NULL)
            {
                puts("Error");
                exit(-1);
            }
        }
        else
        {
            value = (char *)malloc(l + 1);
        }
        for (int i = 0; i < l; i++)
            value[len + i] = val[i];
        value[len + l] = '\0';
        len += l;
        return *this;
    }

    class_string method_substring_number_number(int b, int e)
    {
        if (len < e)
        {
            exit(-1);
        }
        char *val = (char *)malloc(e - b + 1);
        strncpy((char *)val, value + b, e - b);
        val[e - b] = '\0';
        class_string v(val);
        return v;
    }

    char method_at_number(int i)
    {
        return value[i];
    }

    int method_size()
    {
        return len;
    }
};