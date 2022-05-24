#include "file.h"

int Open(char *filename){
    return open((const char*)filename, O_RDWR|O_CREAT|O_APPEND, 0664);
}

int Write(int fd, char *msg, int len) {
    return write(fd, msg, len);
}

void Close(int fd) {
    close(fd);
}