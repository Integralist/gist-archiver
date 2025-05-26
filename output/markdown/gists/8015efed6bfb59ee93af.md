# Ruby: Check Balanced Params

## Check Balanced Params.rb

```ruby
def check_balanced(str)
  str.chars.reduce(0) { |open, char|
    if char == ')' && open == 0
      return false
    elsif char == '('
      open + 1
    elsif char == ')'
      open - 1
    else
      open
    end
  } == 0
end

check_balanced "(()())"
```

