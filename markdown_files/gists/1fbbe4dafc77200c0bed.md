# Ruby: override `new` constructor method using meta programming

[View original Gist on GitHub](https://gist.github.com/Integralist/1fbbe4dafc77200c0bed)

## ruby_override_new.rb

```ruby
module Bar
  module ClassMethods
    # `new` is a class method on the `Class` object
    # It then uses `send` to access `initialize` which would otherwise be a private instance method
    # So it can be overridden by extending the your class with a new `new` class method
    def new(*args, &block)
      super
      p "new constructor defined"
    end
  end

  class << self
    def included(klass)
      p "#{klass} is the class that's including #{self}"
      klass.extend(ClassMethods)
    end
  end
end

class Foo
  include Bar

  def initialize
    p "Foo constructor"
  end
end

Foo.new

# "Foo is the class that's including Bar"
# "Foo constructor"
# "new constructor defined"
```

