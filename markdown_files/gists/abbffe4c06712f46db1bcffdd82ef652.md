# Shell: single line Ruby script in shell 

[View original Gist on GitHub](https://gist.github.com/Integralist/abbffe4c06712f46db1bcffdd82ef652)

**Tags:** #ruby #shell #script

## main.rb

```ruby
cat tokens.txt | xargs -I % ruby -e 'require "digest"; puts Digest::SHA256.hexdigest("%")'
```

