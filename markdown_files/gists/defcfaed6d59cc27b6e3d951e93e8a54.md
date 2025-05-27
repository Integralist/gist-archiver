# Bold text in bash output

[View original Gist on GitHub](https://gist.github.com/Integralist/defcfaed6d59cc27b6e3d951e93e8a54)

## Bold text in bash output.sh

```shell
bold=$(tput bold)
normal=$(tput sgr0)

echo "So, ${bold}I'm bolded${normal} but I'm not bolded"
```

