#include <stdio.h>
#include <math.h>

#define FLOAT 1

int main() {
    printf("%s", "This is from K&R 2ed, p. 214\n");
    union {
        struct {
            int type;
        } n;
        struct {
            int type;
            int intnode;
        } ni;
        struct {
            int type;
            float floatnode;
        } nf;
    } u;

    u.nf.type = FLOAT;
    u.nf.floatnode = 3.14;
    if (u.n.type == FLOAT) {
        printf("What does sin do? %f\n", sin(u.nf.floatnode));
        printf("Value is a float: %f\n", u.nf.floatnode);
    } else {
        printf("Value is not a float.\n");
    }

    return 0;
}
