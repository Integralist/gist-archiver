# [Calculating BigO] #bigo #algorithms #analysis

## Calculating BigO.md

```markdown
> This is condensed information gleened from the excellent [interactivepython.org](http://interactivepython.org/runestone/static/pythonds/AlgorithmAnalysis/BigONotation.html)

## Algorithm

```py
def sumOfN(n):
   theSum = 0
   for i in range(1,n+1):
       theSum = theSum + i

   return theSum

print(sumOfN(10))  # 55
```

## Analysis Steps

You want to quantify the number of operations (or 'steps') in the algorithm.

To do this:

* Identify the basic unit of computation.
* Track any operational constants (e.g. `theSum = 0`).
* Track repeatable operations (e.g. `theSum = theSum + i`).
* Identify the most 'dominant' portion of `T(n)` (see explanation below).

## Explanation

If `T(n)` is the 'size of the problem', then we can say...

```
T(n) == 1+n steps
```

As the problem gets larger, some portion of `T(n)` tends to overpower the rest.

> Note: 'order of magnitude' describes part of `T(n)` that increases the _fastest_ as `n` increases.

We represent order of magnitude in 'Big O' syntax:

```
O(f(n))
```

Where:

```
f(n) == dominant part of T(n)
```

As `n` gets larger, using `T(n) = 1+n` as an example, the 'constant' `1` (i.e. the computation that happened once: `theSum = 0`) becomes less and less significant.

Meaning we can drop `1` from our syntax, resulting in just `O(n)`, and our approximation is just as accurate without it.

## Significant or Insignificant?

I'm going to paste verbatim the original author's comments...

> As another example, suppose that for some algorithm, the exact number of steps is `T(n) = 5n2 + 27n + 1005`. 
> 
> When `n` is small, say `1` or `2`, the constant `1005` seems to be the dominant part of the function. 
> 
> However, as `n` gets larger, the `n2` term becomes the most important. 
> 
> In fact, when `n` is really large, the other two terms become insignificant in the role that they play in determining the final result. 
> 
> Again, to approximate `T(n)` as `n` gets large, we can ignore the other terms and focus on `5n2`. 
> 
> In addition, the coefficient `5` becomes insignificant as `n` gets large. 
> 
> We would say then that the function `T(n)` has an order of magnitude `f(n) = n2`, or simply that it is `O(n2)`.

## Example Analysis

```py
a = 5
b = 6
c = 10

for i in range(n):
   for j in range(n):
      x = i * i
      y = j * j
      z = i * j

for k in range(n):
   w = a*k + 45
   v = b*b

d = 33
```

The above code can be calculated as:

```
T(n) == 3 + 3n2 + 2n + 1
```

Which can be condensed slightly, by combining the singular constants, to:

```
T(n) == 3n2 + 2n + 4
```

The constants `3` and `1` are the top level variable assignments: `a=5`, `b=6`, `c=10` and `d=33`.

The `3n2` is because there are three constant variable assignments (`x`, `y` and `z`) that are occurring within the first set of `for` statements. The top level `for` statement iterates over `n` items, and then does so _again_ hence the `n2` portion of `3n2`.

The `2n` is because there are two constant assignments (`w` and `v`) and they happen `n` times due to the last `for` statement iterating over `n` items.

With this in mind we can say the code is `O(n2)` because when we look at the exponents of each segment of the time analysis (i.e. the condensed version: `3n2 + 2n + 4`) we can see that as `n` grows, the `n2` portion is the most significant.

> Remember: although we write 'Big O' as `O(...)` the underlying principle is `O(f(...))`, where `f(...)` is the dominant part of `T(...)` and when focusing in on the dominant part of the time complexity we drop the constants -- also known as the _coefficient_ -- (e.g. `3n2` thus becomes `n2`). This is because the constants become _insignificant_ as `n` grows.
```

