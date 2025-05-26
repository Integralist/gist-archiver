# Clojure deftype, defrecord, defprotocol

## 1. class.clj

```clojure
; Interface
; Note: it doesn't make much sense to apply an interface to a defrecord
;       considering a defrecord is just a map
(defprotocol Bazzer
  "This is an interface that states a `baz` method should be implemented"
  (baz [this] [a b])) ; you can define multi arity interface, but seemingly can't implement it on a defrecord?
                      ; instead use `extend-protocol` for those situations
                      ; see following example

(extend-protocol Bazzer
  Object ; the return type determines if symbols referenced (e.g. a and b) can be resolved
         ; if not defined (or the wrong type) then errors can occur
  (baz
    ([a] 1)
    ([a b] 2)))

(prn (baz "any value and I'll return 1"))        ; 1
(prn (baz "any two values" "and I'll return 2")) ; 2

; Constructor
(defrecord Foo [a b]
  Bazzer ; enforces the interface (but the error thrown without this defined, doesn't actually clarify)
  (baz [this] "Foo Bazzer")) ; associate a function with our `defrecord` map

; Constructor
(defrecord Bar [a b] Bazzer
  (baz [this] (str "Bar Bazzer -> " a " / " b)))

; Either pass in each argument to the constructor separately
(def foo (->Foo :bar :baz)) ; user.Foo{:a :bar, :b :baz}

; Or pass a single argument with keys that align with the class' required parameters
(def bar (map->Bar {:a 1 :b 2})) ; user.Foo{:a 1, :b 2}

(prn bar)      ; user.Foo{:a 1, :b 2}
(prn (:a foo)) ; :bar
(prn (:b foo)) ; :baz

; mutate and return
(prn (assoc foo :b :qux)) ; user.Foo{:a :bar, :b :qux}

; but the source is immutable
(prn foo) ; user.Foo{:a :bar, :b :baz}

(:b (assoc foo :b :qux)) ; :qux

(def basil (Foo. "hai" "hi"))
(def baril (Bar. "boo" "yah"))
(prn (baz basil)) ; "Foo Bazzer""
(prn (baz baril)) ; "Bar Bazzer -> boo / yah""
; (prn (baz baril :a :b))
```

## 2. interface.clj

```clojure
; Interface
(defprotocol MyInterface 
  (foo [this]) ; `this` is required to let the interface refer to the class
  (bar [this] [this x] [this x y] [this x y z])) ; multi-arity method signature defined

; `deftype` dynamically generates compiled bytecode for the specified identifier (e.g. MyClass)
(deftype MyClass [a b c]
  MyInterface ; implement the specified protocol (i.e. interface)
    
    ; each function's scope is defined by 
    ; the object provided as the first argument
    ; i.e. something that is of the `MyClass` type
    (foo [this] a)
    (bar [this] b)
    (bar [this x] (+ c x))
    (bar [this x y] (+ c x y))
    (bar [this x y z] (+ c x y z)))

(def obj (MyClass. 1 2 3))

(prn (foo obj))       ; 1
(prn (bar obj))       ; 2
(prn (bar obj 1))     ; 4 (3 + )
(prn (bar obj 1 2))   ; 6 (3 + 1 + 2)
(prn (bar obj 1 2 3)) ; 9 (3 + 1 + 2 + 3)
```

## 0. description.md

```markdown
- `defprotocol`: defines an interface
- `deftype`: create a bare-bones object which implements a protocol
- `defrecord`: creates an immutable persistent map which implements a protocol

Typically you'll use `defrecord` (or even a basic `map`);  
unless you need some specific Java inter-op,  
where by you'll want to use `deftype` instead.

> Note: `defprotocol` allows you to add new abstractions in a clean way
Rather than (like OOP) having polymorphism on the class itself,
polymorphic functions are created in namespaces.
Meaning different namespaces can implement different functionality
```

