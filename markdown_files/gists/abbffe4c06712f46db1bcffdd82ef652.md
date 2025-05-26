# [Ruby single line script in shell] #ruby #shell #script

## main.rb

```ruby
cat tokens.txt | xargs -I % ruby -e 'require "digest"; puts Digest::SHA256.hexdigest("%")'
```

