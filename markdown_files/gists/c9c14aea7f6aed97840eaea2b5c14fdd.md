# [Short Bash if else statement] 

[View original Gist on GitHub](https://gist.github.com/Integralist/c9c14aea7f6aed97840eaea2b5c14fdd)

**Tags:** #bash #shell #if #conditional

## Short if else statement.bash

```shell
# yay
export FOO=true
[[ $FOO == "true" ]] && echo yay || echo nay

# nay
export BAR=false
[[ $BAR == "true" ]] && echo yay || echo nay
```

