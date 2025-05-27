# [Brew switch to custom versions of software] #homebrew #brew #install #versions #switch

[View original Gist on GitHub](https://gist.github.com/Integralist/cbbbb95b571bd08bb5aa)

## 1. [CURRENT] Brew switch.md

If you wish to switch your `python` command to use a different Python interpreter (and it's a Python version you previously had installed using Homebrew):

```bash
brew info python           # To see what you have previously installed
brew switch python 3.x.x_x # Ex. 3.6.5_1
```

> NOTE: this might not be wise to do as you might have other software that relies on the Python interpreter you have currently (i.e. before the switch).

Otherwise to install multiple Python versions using Homebrew (instead of something like `pyenv`)...

Template:

```bash
brew install https://address/to/the/formula/FORMULA_NAME.rb
```

Example Python 3 template:

```bash
brew install https://raw.githubusercontent.com/Homebrew/homebrew-core/COMMIT_IDENTIFIER/Formula/python.rb
```

Now you can't view the commits via the GitHub UI because there are so many, so you'll have to clone the homebrew-core repository and do it via the command line:

```bash
git clone git@github.com:Homebrew/homebrew-core.git
git log master -- Formula/python.rb
```

From there you can search for the commits which update the Python version by searching for the message `python: update <version> bottle.`

Example: we want to install Python version `3.7.5`

```bash
git log --oneline --grep '^python: update .\+ bottle'

c14db427e7 python: update 3.7.7 bottle.
f02346bd48 python: update 3.7.6_1 bottle.
7ef074a882 python: update 3.7.6 bottle.
2efdfe5519 python: update 3.7.5 bottle.    # << HERE IT IS!
6589f0f6f5 python: update 3.7.4_1 bottle.
e9004bd764 python: update 3.7.4_1 bottle.
1c2239bfdf python: update 3.7.4_1 bottle.
48aba7218e python: update 3.7.4 bottle.
c24d6bcd47 python: update 3.7.3 bottle.
22c80fc362 python: update 3.7.3 bottle.

...etc
```

Now we know the commit hash for the Python version `3.7.5` (i.e. `2efdfe5519`) we want we can install that version like so:

```bash
brew unlink python # ONLY if you have installed (with brew) another version of python 3
brew install --ignore-dependencies https://raw.githubusercontent.com/Homebrew/homebrew-core/2efdfe5519/Formula/python.rb
```

## 2. [OLD] Brew switch.md

- `brew versions {package}`
- `brew switch {package} {version}`

OR

- `brew uninstall {package}`
- `brew versions {package}`
- `git checkout {commit} Library/Formula/{package}.rb`
- `brew install {package}`

