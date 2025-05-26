# [Sed Replace Pattern with File Contents] #sed #bash #replace

## Sed Replace Pattern with File Contents.bash

```shell
cd -- "$(mktemp -d)" 
printf '%s\n' 'nuts' 'bolts' > first_file.txt
printf '%s\n' 'foo' 'bar' 'baz' > second_file.txt
sed -e '/bar/r ./first_file.txt' second_file.txt

### output...
foo
bar
nuts
bolts
baz
### ...notice 'bar' in second_file.txt has the contents from first_file.txt appended after it
```

