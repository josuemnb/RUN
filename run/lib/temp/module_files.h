#include "../../libc/run.h"
;
;

    #include <stdio.h>
    #include <stdlib.h>
    #include <unistd.h>
    #include <fcntl.h>
    #include <sys/stat.h>
 ;
;

class files_class_FileInfo {
private:
class_string field___name;
private:
number field___size;
private:
bool field___exist;
public:
number method_size() {
return field___size;
}
public:
bool method_exist() {
return field___exist;
}
public:
files_class_FileInfo(class_string param_n) {
field___name=param_n;

            int fd = open(field___name.value,O_RDONLY);
            if(fd>-1) {
                field___exist = fd>-1;
                struct stat s;
                fstat(fd,&s);
                field___size = s.st_size;
            }
         ;
}
};
;

class files_class_FileReader {
private:
number field___fp;
private:
class_string field___name;
private:
bool field___has;
public:
void method_goto_number(number param_pos) {

            lseek(field___fp,param_pos,SEEK_SET);
         ;
}
public:
bool method_close() {
bool var_ret;

            var_ret = close(field___fp)>-1;
         ;
return var_ret;
}
public:
void method_open_string(class_string param_name) {
field___name=param_name;

            field___fp  = open(field___name.value,O_RDONLY);
         ;
field___has=true;
}
public:
number method_readNumber() {
number var_n;

            read(field___fp,&var_n, sizeof(var_n));              
         ;
return var_n;
}
public:
class_string method_read_number(number param_size) {
class_string var_s;
if(field___has==false) {
return var_s;
}

            char buf[param_size];
            size_t nread;
            nread = read(field___fp, buf, param_size);   
            buf[nread] = '\0';     
            var_s = buf;
         ;
return var_s;
}
public:
class_string method_readLine() {
class_string var_s;

            char buf[128];
            char c;
            int i = 0,nread;
            while((nread =read(field___fp, &c, 1))>0) {
                if(c=='\r' || c=='\n') {
                    buf[i] = '\0';                   
                    var_s = buf;
                    break;
                }
                buf[i] = c;
                i++;
            }                       
            if(nread<1)
                field___has=false;                
         ;
return var_s;
}
public:
bool method_has() {
return field___has;
}
public:
number method_position() {
number var_ret;

            var_ret = lseek(field___fp,0,SEEK_CUR);
         ;
return var_ret;
}
};
;

class files_class_FileAppender {
private:
number field___fp;
private:
class_string field___name;
public:
void method_append_string(class_string param_s) {

            write(field___fp,param_s.value,param_s.len);
         ;
}
public:
bool method_close() {
bool var_ret;

            var_ret = close(field___fp)>-1;
         ;
return var_ret;
}
public:
void method_open_string(class_string param_name) {
field___name=param_name;

            field___fp  = open(field___name.value,O_APPEND|O_CREAT|O_WRONLY, 0644);
            if(field___fp<1) {
                fprintf(stderr, "New File Error: %d\n", errno);
            }
         ;
}
};
;

class files_class_FileWriter {
private:
number field___fp;
private:
class_string field___name;
public:
bool method_close() {
bool var_ret;

            var_ret = close(field___fp)>-1;
         ;
return var_ret;
}
public:
void method_open_string(class_string param_name) {
field___name=param_name;

            field___fp  = creat(field___name.value, 0644);
            if(field___fp<1) {
                fprintf(stderr, "New File Error: %d\n", errno);
            }
         ;
}
public:
void method_write_number(number param_n) {

           write(field___fp,&param_n,sizeof(long));        
         ;
}
public:
void method_write_string(class_string param_s) {

            write(field___fp,param_s.value,param_s.len);
         ;
}
public:
number method_position() {
number var_ret;

            var_ret = lseek(field___fp,0,SEEK_CUR);
         ;
return var_ret;
}
public:
void method_goto_number(number param_pos) {

            lseek(field___fp,param_pos,SEEK_SET);
         ;
}
};
;
files_class_FileReader func_files_open_string(class_string param_n) {
files_class_FileReader var_f;
var_f.method_open_string(param_n);
return var_f;
}
;
files_class_FileWriter func_files_new_string(class_string param_n) {
files_class_FileWriter var_f;
var_f.method_open_string(param_n);
return var_f;
}
;
files_class_FileAppender func_files_append_string(class_string param_n) {
files_class_FileAppender var_f;
var_f.method_open_string(param_n);
return var_f;
}
;
files_class_FileInfo func_files_info_string(class_string param_n) {
return (param_n);
}
;
