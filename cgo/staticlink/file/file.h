#include <fcntl.h>
#include <errno.h>
#include <unistd.h>

int Open(char *filename);

int Write(int fd, char *msg, int len);

void Close(int fd);