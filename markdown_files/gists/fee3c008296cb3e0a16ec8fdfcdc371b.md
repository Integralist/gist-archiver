# [Python Exception Handling Guidelines] #python #exceptions #design #guidelines

[View original Gist on GitHub](https://gist.github.com/Integralist/fee3c008296cb3e0a16ec8fdfcdc371b)

**Tags:** #python #exceptions #design #guidelines

## Python Exception Handling Guidelines.md

It's best to raise your exceptions at the point in the code where errors are going to occur.

When you have deeply nested code modules, then don't try to handle exceptions in-between, just raise the exception (instrument there too) and then let a _top level_ handler catch generalized exceptions to do its own instrumentation (if necessary).

This means it's a good idea for the top level handler to catch a generalized parent exception class type, and have your internal code raise subclass exceptions from that parent class.

