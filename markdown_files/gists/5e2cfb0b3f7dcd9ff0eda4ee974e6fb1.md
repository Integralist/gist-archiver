# [Python Security Hashing] 

[View original Gist on GitHub](https://gist.github.com/Integralist/5e2cfb0b3f7dcd9ff0eda4ee974e6fb1)

**Tags:** #python #security #hashing #hash

## Python Security Hashing.py

```python
# see also: https://docs.python.org/3.7/library/crypto.html

from hashlib import blake2b

h = blake2b(key=b'pseudorandom key', digest_size=16)
h.update(b'message data')
h.hexdigest()  # '3d363ff7401e02026f4a4687d4863ced'
```

