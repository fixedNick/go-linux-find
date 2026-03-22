



# go-linux-find

Implementation of the `find(1)` utility in the **Golang** language.

## CmdLets

1. `-name` search by file name
2. `-iname` search by file name case-insensitive
3. `-depth` `-maxdepth` `-mindepth` set depth relative to root
4. `-type` support flags:
    - **f** - file
    - **d** - directory
5. `-path` search by path (part of path)
6. `-ipath` search by path (part of path) case-insensitive
7. `-empty` search only directories with no files or Zero size files