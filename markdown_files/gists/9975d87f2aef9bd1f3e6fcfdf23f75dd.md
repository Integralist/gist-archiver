# [Install Python package directly for the interpreter Vim is using] #vim #python

## Install Python package directly for the interpreter Vim is using.md

```markdown
Sometimes Vim's Python binary can't find a package you've installed.

- in vim  
  `:py3 import sys; print(sys.path)`

- navigate to that location and find the python binary, in my case:  
  `cd /usr/local/opt/python@3.8/Frameworks/Python.framework/Versions/3.8/bin`

- install isort using this binary  
  `./python3.8 -m pip install isort`
```

