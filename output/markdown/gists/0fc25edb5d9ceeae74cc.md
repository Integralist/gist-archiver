# Bash Strict Mode (http://redsymbol.net/articles/unofficial-bash-strict-mode/)

## Bash Strict Mode.sh

```shell
#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# -e: immediately exit if any command has a non-zero exit status
# -u: reference to any variable you haven't previously defined is an error
# -o: prevents errors in a pipeline from being masked
# IFS new value is less likely to cause confusing bugs when looping arrays or arguments (e.g. $@)

# Example: positional parameters
# Solution: use 'default value' setup (set default value to an empty string)

name=${1:-}
if [[ -z "$name" ]]; then
    echo "usage: $0 NAME"
    exit 1
fi
echo "Hello, $name"

# Example: pipeline masking issues
# Solution: short-circuit with boolean operator

count=$(grep -c some-string some-file || true)
echo "count: $count"

# Solution: disable -e temporarily (in order to get return value)

set +e
count=$(grep -c some-string some-file)
retval=$?
set -e

# grep's return code is 0 when one or more lines match;
# 1 if no lines match; and 2 on an error. This pattern
# lets us distinguish between them.

echo "return value: $retval"
echo "count: $count"

# Example: chaining doesn't stop execution if there's an error
# Solution: use a block {} to wrap code you want to stop execution if it fails

# first_task && second_task && third_task
# next_task
# ...if second_task fails then third_task doesn't run
# ...but next_task would
# ...instead use the following...

first_task && {
    second_task
    third_task
}
next_task

# Example: Avoid inline if statements (always use long form if/then)
# Notes: inline if statements are issues when they're the last line of the script
```

