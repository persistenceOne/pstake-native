### Abstract

This document specifies how to use the `git subtree` command.

### git subtree

Git subtree is the most common replacement for Git submodule. A Git subtree is a replica of a Git repository that has
been dragged into the main repository. A Git submodule is a reference to a particular commit in a different repository.
Git subtrees, which were first introduced in Git 1.7.11, help you make a copy of any repo into a subdirectory of
another.

Example :
Adding a subtree into a Primary project (ProjA) which has modules, from a second project (ProjB). In order to accomplish
this, maintain a copy ProjB.

Split ProjB:

```
cd ProjB
git checkout -b split-maint
git subtree split --prefix=important/dir --branch=module-for-A
```

Add the subtree to ProjA:

``` 
cd ProjA
git remote add ProjB_remote /path/to/ProjB
git fetch ProjB_remote
git subtree add --prefix=modules/projB_mod ProjB_remote/module-for-A --message="commit message"
```

Steps to update the module

```
cd ProjA
git remote add ProjB_remote /path/to/ProjB
git fetch ProjB_remote
git subtree merge --prefix=modules/projB_mod ProjB_remote/module-for-A --message="commit message"
``` 

Now resolve the conflicts and push the changes! 