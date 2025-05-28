# [Memory Sharing] 

[View original Gist on GitHub](https://gist.github.com/Integralist/6d0ee04eb68ecfdaab8509b1eccadc98)

**Tags:** #python #rss #resident #virtual #memory

## Memory Sharing.md

ECS Task has memory allocation of 500mb.

It's a multi-process tornado service with 4 processes.

Each process is reporting 200+ mb of RSS (resident) memory usage, which would mean total memory should be larger than 500mb.

But as processes share memory, the total RSS is likely to be the sum of the memory shared across all child processes + the "unique" memory used by each child process as they diverge from one another.

