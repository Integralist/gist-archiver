# Ruby OOP vs FP (examples are from ThoughtBot's Weekly Iteration -> you should subscribe!)

## 1. API.rb

```ruby
# Examples taken from ThoughtBot's Weekly Iteration
# We need to implement a solution that allows us to cleanly use two separate APIs

# Our fake APIs
PayPal.charge!(auth_code, 25)
Stripe::CreditCard.new(credit_card_token).charge(25)
```

## 2. OOP.rb

```ruby
# Both of the below classes are utilising the principle of polymorphism

# Strategy pattern (passing in object that determines the strategy to be used)
class Payment < Struct.new(:price, :client)
  def charge
    client.charge(price)
  end
end

# Adapter pattern (creates a consistent internal API when using an external API)
class PayPalAdapter < Struct.new(:auth_code)
  def charge(amount)
    PayPal.charge!(auth_code, amount)
  end
end

# Consistent API thanks to Adapter pattern
Payment.new(25, Stripe::CreditCard.new(credit_card_token)).charge
Payment.new(25, PayPalAdapter.new(auth_code)).charge
```

## 3. FP.rb

```ruby
# The .call() method allows us to execute the provided method
# and pass through an argument to that method.
# In the following code examples we call Payment.new and pass in a method which can be called
class Payment < Struct.new(:price, :client)
  def charge
    client.call(price)
  end
end

# Use .method() to return a Proc which the above charge method can trigger a .call() on
# http://www.ruby-doc.org/core-2.1.1/Method.html
Payment.new(
  25, 
  Stripe::CreditCard.new(credit_card_token).method(:charge)
).charge

# Again use .method() to return a Proc
# This means we can then use Currying to return another Proc 
# That returned Proc has its first argument pre-filled with :auth_code
# The above charge method can then trigger .call() on the curried Proc (passing in the remaining argument)
Payment.new(
  25, 
  PayPal.method(:charge!).curry(:auth_code)
).charge
```

