#include <fcntl.h>
#include <errno.h>
#include <unistd.h>
#include <string.h>

void SayHello(_GoString_ s);

static int Open(char *filename){
    return open((const char*)filename, O_RDWR|O_CREAT|O_APPEND, 0664);
}

static int Write(int fd, char *msg, int len) {
    return write(fd, msg, len);
}

static void Close(int fd) {
    close(fd);
}