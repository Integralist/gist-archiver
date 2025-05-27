# [Install Python package directly for the interpreter Vim is using] #vim #python

[View original Gist on GitHub](https://gist.github.com/Integralist/9975d87f2aef9bd1f3e6fcfdf23f75dd)

## Install Python package directly for the interpreter Vim is using.md

Sometimes Vim's Python binary can't find a package you've installed.

- in vim  
  `:py3 import sys; print(sys.path)`

- navigate to that location and find the python binary, in my case:  
  `cd /usr/local/opt/python@3.8/Frameworks/Python.framework/Versions/3.8/bin`

- install isort using this binary  
  `./python3.8 -m pip install isort`

