# [Python coroutine comparison functions] 

[View original Gist on GitHub](https://gist.github.com/Integralist/1efc8dcfc0b1e9e8e8b89a4b2019f3af)

**Tags:** #python #python3 #asyncio #async #coroutine #concurrency

## Python coroutine comparison functions.py

```python
"""the results aren't what you might expect."""

import asyncio

@asyncio.coroutine
def old_style_coroutine():
    yield from asyncio.sleep(1)

async def main():
    await old_style_coroutine()

asyncio.iscoroutine(old_style_coroutine)
False

asyncio.iscoroutine(main)
False

asyncio.iscoroutinefunction(old_style_coroutine)
True

asyncio.iscoroutinefunction(main)
True

import inspect

inspect.iscoroutine(old_style_coroutine)
False

inspect.iscoroutine(main)
False

inspect.iscoroutinefunction(old_style_coroutine)
False

inspect.iscoroutinefunction(main)
True

inspect.isawaitable(old_style_coroutine)
False

inspect.isawaitable(main)
False

async def foobar():
    await asyncio.sleep(1)

inspect.isawaitable(foobar)
False

inspect.isawaitable(asyncio.sleep)
False
```

