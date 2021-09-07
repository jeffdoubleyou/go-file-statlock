# go-file-statlock
File locking based on age of lock file.  Basic file lock for NFS filesystem.

# Why

I am using this to handle services that run on multiple systems that need to know if another process is running on a different machine without any other access to the other system outside of a shared NFS filesystem.

Normal file locking isn't supported on NFS and I could not find any suitable library to handle this.

# Usage

tbd
