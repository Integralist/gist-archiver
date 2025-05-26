# Python Class Decorator (AOP)

## Python Class Decorator (AOP).py

```python
def aop(cls):
    class Wrapper:
        def __init__(self, *args):
            self.wrapped = cls(*args) # create instance of decorated class
        def __getattr__(self, name):
            print('Getting the {} of {}'.format(name, self.wrapped))
            return getattr(self.wrapped, name)
    return Wrapper

@aop
class Foo:
    def test_foo(self):
        print("foo called")
    def test_bar(self):
        print("bar called")

f = Foo()
f.test_foo()

# Getting the test_foo of <__main__.Foo object at 0x1099c36a0>
# foo called
```

