# Sort an Array according to another Array

## Sort an Array according to another Array.rb

```ruby
a1 = [4, 5, 2, 3, 1]
a2 = [{:id => 5}, {:id => 2}, {:id => 1}, {:id => 3}, {:id => 4}]

a2.sort_by{|x| a1.index(x[:id])}

=> [{:id=>4}, {:id=>5}, {:id=>2}, {:id=>3}, {:id=>1}]
```

