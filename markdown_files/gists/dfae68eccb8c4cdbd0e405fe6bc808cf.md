# [Python CPU and Memory Profiling Tools] 

[View original Gist on GitHub](https://gist.github.com/Integralist/dfae68eccb8c4cdbd0e405fe6bc808cf)

**Tags:** #tags: python, profiling, perf

## 0. README.md

It is worth considering the variation in load that you get on a normal computer. 

Many background tasks are running (e.g. Dropbox, backups) that could impact the CPU and disk resources at random.

Tools available:

* `timeit` module: understand behavior of statements and functions
* `cProfile`: understand which functions in your code take the longest to run (high-level only)
* `line_profiler`: profiles your chosen functions on a line-by-line basis (granular)
* `memory_profiler`: chart RAM usage over time on a labelled chart (function level)
* `heapy`: track all of the objects inside Python’s memory (useful for hunting down memory-leaks)
* `perf stat`: understand number of instructions executed on CPU and how efficiently the CPU’s caches are utilized
* CPython: understand Python stack based vm operates (understand why certain coding styles run more slowly than others)

> [Reference article](http://www.marinamele.com/7-tips-to-time-python-scripts-and-control-memory-and-cpu-usage)

### Locating Python Packages

```bash
find / -type d -name <package_name>

# for example, the result could look like...
#
# ./usr/local/lib/python3.6/site-packages/package_name
```

## 1. app-source.py

```python
"""This module is for profiling purposes."""

from pprint import pprint


def factorize_naive(n):
    """
    Naive factorization method.

    Take integer 'n', return list of factors.
    """
    if n < 2:
        return []
    factors = []
    p = 2

    while True:
        if n == 1:
            return factors

        r = n % p
        if r == 0:
            factors.append(p)
            n = n / p
        elif p * p >= n:
            factors.append(n)
            return factors
        elif p > 2:
            # Advance in steps of 2 over odd numbers
            p += 2
        else:
            # If p == 2, get to 3
            p += 1

    assert False, "unreachable"


def serial_factorizer(nums):
    """Process factorials lots of times."""
    return {n: factorize_naive(n) for n in nums}


if __name__ == "__main__":
    pprint(serial_factorizer([10000, 100, 4450, 8320, 500000]))

```

## 2. basic decorator timer.py

```python
"""This module is for profiling purposes."""

from time import time
from pprint import pprint
from functools import wraps


def timefn(fn):
    """Simple timer decorator."""
    @wraps(fn)
    def measure_time(*args, **kwargs):
        t1 = time()
        result = fn(*args, **kwargs)
        t2 = time()
        print(f"@timefn: {fn.__name__} took {str(t2 - t1)} seconds")
        return result
    return measure_time


@timefn
def factorize_naive(n):
    """
    Naive factorization method.

    Take integer 'n', return list of factors.
    """
    if n < 2:
        return []
    factors = []
    p = 2

    while True:
        if n == 1:
            return factors

        r = n % p
        if r == 0:
            factors.append(p)
            n = n / p
        elif p * p >= n:
            factors.append(n)
            return factors
        elif p > 2:
            # Advance in steps of 2 over odd numbers
            p += 2
        else:
            # If p == 2, get to 3
            p += 1

    assert False, "unreachable"


def serial_factorizer(nums):
    """Process factorials lots of times."""
    return {n: factorize_naive(n) for n in nums}


if __name__ == "__main__":
    pprint(serial_factorizer([10000, 100, 4450, 8320, 500000]))
```

## 3. timeit command line module.sh

```shell
# python -m timeit -n <average> -r <repetitions> -s "<import your app>" "<code to be run>"
#
# for longer-running functions it can be sensible to specify the number of loops which are averaged (-n 4) and a number of repetitions of the loops (-r 5)

python -m timeit -n 4 -r 5 -s "from pprint import pprint; import app" "pprint(app.serial_factorizer([10000, 100, 4450, 8320, 500000]))"
```

## 4. cProfile command line module.sh

```shell
python -m cProfile -s cumulative app.py
```

## 5. memory_profiler.sh

```shell
pip install memory_profiler
pip install psutil

# From VM...
mprof run python /app/video_player.py

# From Host...
ulimit -n 1024
siege <url> -c 255 -t 30S
ll | grep mprofile | tail -n 1 | awk '{print $9}' | xargs cat | tail -n 1 # locate latest mprofile_<timestamp>.dat and display its last line

# MEM 36.324219

```

