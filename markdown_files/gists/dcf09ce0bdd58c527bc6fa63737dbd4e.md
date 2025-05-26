# Python access the stack trace failed functions

## Python access the stack trace failed functions.py

```python
import sys
import traceback
import logging

def i_will_fail():
    raise Exception

def i_call_the_failing_function():
    i_will_fail()

try:
    i_call_the_failing_function()
except Exception:
    _, _, tb = sys.exc_info()
    for frame in traceback.extract_tb(tb):
        logging.error('Module name: %s; scope name: %s', frame[0], frame[2])
        
"""
ERROR:root:Module name: exception-traceback.py; scope name: <module>                                                                                                        
ERROR:root:Module name: exception-traceback.py; scope name: i_call_the_failing_function                                                                                     
ERROR:root:Module name: exception-traceback.py; scope name: i_will_fail  
"""
```

