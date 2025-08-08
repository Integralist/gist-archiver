# How Our Search Works

Our search is designed to be both fast and intelligent, helping you find the gists you're looking for even if you don't remember the exact wording. This is achieved through a combination of two key features: prefix matching and word stemming.

## Prefix Matching (Search-As-You-Type)

The search function matches partial words from the beginning. This allows you to see results immediately as you type.

For example, if you type `iter`, the search will instantly find gists containing words that start with `iter`, such as `iterator` or `iteration`.

## Word Stemming (Finding Related Words)

To provide more relevant results, our search also uses a technique called "stemming." Stemming reduces words to their common root or "stem." This allows the search to find gists that are conceptually related, even if they don't use the exact word you searched for.

For example, the words:
- `iterator`
- `iterators`
- `iterating`
- `iteration`

...are all reduced to the common root stem: **`iter`**.

This means that searching for any one of these words will find gists containing any of the others, making the search much more powerful.

### Example: Why does a search for "iterator" show unexpected results?

You might search for `iterator` and see gists that don't appear to contain that exact word. For instance, you might see results including:

- A gist about Go concurrency (`gists/05fb91c42021195b727be5afb28122ec.html`)
- A gist about parallel processing (`gists/1391150b69eebcaac98984627ba26b7d.html`)

This happens because of stemming:

1.  The search for `iterator` is stemmed to its root, `iter`.
2.  The Go concurrency gist contains the word **`iterating`**, which is also stemmed to `iter`.
3.  The parallel processing gist contains the word **`iteration`**, which is also stemmed to `iter`.

Because both documents contain words that share the same root as your search term, they are correctly identified as relevant results.

## Summary

By combining prefix matching for speed and word stemming for intelligence, our search helps you find the most relevant gists quickly and efficiently.
