# [Vim arg list and Search Replacement] #vim #search #replace #args #case #lower #upper

## Vim arg list and Search Replacement.md

```markdown
Imagine we have a file containing...

```python
from foo import Bar
```

We want `Bar` to be `bar`.

We could do that across all files like so:

```viml
:args
:args **/*.py
:argdo %s/\v(from foo import )B(ar)/\1b\2/e | update
```

> `e` flag tells Vim to ignore any errors

Alternatively instead of hardcoding the lowercase `b` we could have done this dynamically using either `\l` which lowercases the first character inside of a backreference or `\L` which lowercases the entire backreference ([Reference](https://vim.fandom.com/wiki/Changing_case_with_regular_expressions)).

Examples...

```viml
:argdo %s/\v(from foo import Bar)/\L\1/e | update
:argdo %s/\vfrom foo import (B)ar/\l\1/e | update
:argdo %s/\vfrom foo import (*+)/\l\1/e | update
```
```

