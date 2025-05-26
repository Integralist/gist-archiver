# Sift example that demonstrates how to ignore a directory and also display the line numbers whilst using a regex pattern with a word boundary

## sift.sh

```shell
# https://sift-tool.org/
# https://github.com/svent/sift
# brew install sift

sift --exclude-dirs 'Godeps' --line-number '\bint\b'
```

