# Python Singleton

[View original Gist on GitHub](https://gist.github.com/Integralist/59d6a84a7a083a60ebdeceacf6f63cd9)

## Python Singleton.py

```python
# singleton.py

instance = None


class Foo():
    def __new__(cls, arg):
        global instance
        if instance is not None:
            return instance
        instance = object.__new__(cls)
        instance.foo = arg
        return instance

    def __init__(self, arg):
        print(self.foo)  # __init__ called each time
                         # but foo is set only once

# consumer code within a separate module...

f = Foo('a')
f = Foo('b')
f = Foo('c')

```

