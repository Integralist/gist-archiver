# Bold text in bash output

## Bold text in bash output.sh

```shell
bold=$(tput bold)
normal=$(tput sgr0)

echo "So, ${bold}I'm bolded${normal} but I'm not bolded"
```

