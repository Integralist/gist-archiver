# Perl: add language to markdown code block #perl #regex #codeblock #markdown

[View original Gist on GitHub](https://gist.github.com/Integralist/bc358eb37de04fa536e34fcc4a7e8cba)

## cli.md

```bash
# Ensure all code blocks without a language get 'plain' appended.
# NOTE: constraint is that there needs to be an empty line preceding the code block.
cat file.mdx | perl -pe 'BEGIN { undef $/ }; s/(?<=^\n)(```)(\n.+?```)/\1plain\2/s'
```

> **NOTE**: The Perl variable `$/` stores the line ending, used for processing files line-by-line. By calling `undef` on it, we cause Perl to slurp the entire file all in one string at once for us to process.
>
> **WARNING**: If the code block has a non-empty preceding line, and the contents of the code block has an empty line before the closing code fence, then the current pattern will accidentally match. So if that's the case, we need to ensure all code blocks don't have an unnecessary empty line at the bottom of them.

## file.mdx

```mdx
Example: https://regex101.com/r/6wcYU9/1

The following block will get `plain` appended:

```
Content here
```

And this one:

```
Content here
```

BUT for this to work we expect the code block to have an empty line preceding it.

So the following will NOT get `plain` appended because it violates the empty preceding line expectation:

a non-empty line before the code block
```
content here
```

Here is a code block that has a language already:

```http
some text here
```

One more time, let's see an expected code block match:

```
some stuff
```
```

