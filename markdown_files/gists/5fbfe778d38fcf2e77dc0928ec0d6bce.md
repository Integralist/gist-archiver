# Log to SysLog on the command line (terminal)

[View original Gist on GitHub](https://gist.github.com/Integralist/5fbfe778d38fcf2e77dc0928ec0d6bce)

## Log to SysLog on the command line (terminal).sh

```shell
logger integralist

cat /var/log/system.log | grep integralist
# Nov 19 12:34:16 MarkMcDonnell-MBPr markmcdonnell[11781]: integralist
```

