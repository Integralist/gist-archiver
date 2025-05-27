# 01_compojure.md

[View original Gist on GitHub](https://gist.github.com/Integralist/acb341d6b0af7bbe4be0)

## 00_destructuring.md

Clojure Destructuring Tutorial and Cheat Sheet
==============================================

([Related blog post](http://john2x.com/blog/clojure-destructuring))

Simply put, destructuring in Clojure is a way extract values from a datastructure and bind them to symbols, without having to explicitly traverse the datstructure. It allows for elegant and concise Clojure code.

Vectors
-------

**Syntax:** `[symbol another-symbol] ["value" "another-value"]`

```clojure
(def my-vector [:a :b :c :d])
(def my-nested-vector [:a :b :c :d [:x :y :z]])

(let [[a b c d] my-vector]
  (println a b c d))
;; => :a :b :c :d

(let [[a _ _ d [x y z]] my-nested-vector]
  (println a d x y z))
;; => :a :d :x :y :z
```

You don't have to match the full vector.

```clojure
(let [[a b c] my-vector]
  (println a b c))
;; => :a :b :c
```

You can use `& the-rest` to bind the remaining part of the vector to `the-rest`. 

```clojure
(let [[a b & the-rest] my-vector]
  (println a b the-rest))
;; => :a :b (:c :d)
```

When a destructuring form "exceeds" a vector (i.e. there not enough items in the vector to bind to), the excess symbols will be bound to `nil`.

```clojure
(let [[a b c d e f g] my-vector]
  (println a b c d e f g))
;; => :a :b :c :d nil nil nil
```

You can use `:as some-symbol` as the *last two items* in the destructuring form to bind the whole vector to `some-symbol`

```clojure
(let [[:as all] my-vector]
  (println all))
;; => [:a :b :c :d]

(let [[a :as all] my-vector]
  (println a all))
;; => :a [:a :b :c :d]

(let [[a _ _ _ [x y z :as nested] :as all] my-nested-vector]
  (println a x y z nested all))
;; => :a :x :y :z [:x :y :z] [:a :b :c :d [:x :y :z]]
```

You can use both `& the-rest` and `:as some-symbol`.

```clojure
(let [[a b & the-rest :as all] my-vector]
  (println a b the-rest all))
;; => :a :b (:c :d) [:a :b :c :d]
```

### Optional arguments for functions

With destructuring and the `& the-rest` form, you can specify optional arguments to functions.

```clojure
(defn foo [a b & more-args]
  (println a b more-args))
(foo :a :b) ;; => :a :b nil
(foo :a :b :x) ;; => :a :b (:x)
(foo :a :b :x :y :z) ;; => :a :b (:x :y :z)

(defn foo [a b & [x y z]]
  (println a b x y z))
(foo :a :b) ;; => :a :b nil nil nil
(foo :a :b :x) ;; => :a :b :x nil nil
(foo :a :b :x :y :z) ;; => :a :b :x :y :z
```

Maps
----

**Syntax:** `{symbol :key, another-symbol :another-key} {:key "value" :another-key "another-value"}`

```clojure
(def my-hashmap {:a "A" :b "B" :c "C" :d "D"})
(def my-nested-hashmap {:a "A" :b "B" :c "C" :d "D" :q {:x "X" :y "Y" :z "Z"}})

(let [{a :a d :d} my-hashmap]
  (println a d))
;; => A D

(let [{a :a, b :b, {x :x, y :y} :q} my-nested-hashmap]
  (println a b x y))
;; => A B X Y
```

Similar to vectors, if a key is not found in the map, the symbol will be bound to `nil`.

```clojure
(let [{a :a, not-found :not-found, b :b} my-hashmap]
  (println a not-found b))
;; => A nil B
```

You can provide an optional default value for these missing keys with the `:or` keyword and a map of default values.

```clojure
(let [{a :a, not-found :not-found, b :b, :or {not-found ":)"}} my-hashmap]
  (println a not-found b))
;; => A :) B
```

The `:as some-symbol` form is also available for maps, but unlike vectors it can be specified anywhere (but still preferred to be the last two pairs).

```clojure
(let [{a :a, b :b, :as all} my-hashmap]
  (println a b all))
;; => A B {:a A :b B :c C :d D}
```

And combining `:as` and `:or` keywords (again, `:as` preferred to be the last).

```clojure
(let [{a :a, b :b, not-found :not-found, :or {not-found ":)"}, :as all} my-hashmap]
  (println a b not-found all))
;; => A B :) {:a A :b B :c C :d D}
```

There is no `& the-rest` for maps.

### Shortcuts

Having to specify `{symbol :symbol}` for each key is repetitive and verbose (it's almost always going to be the symbol equivalent of the key), so shortcuts are provided so you only have to type the symbol once.

Here are all the previous examples using the `:keys` keyword followed by a vector of symbols:

```clojure
(let [{:keys [a d]} my-hashmap]
  (println a d))
;; => A D

(let [{:keys [a b], {:keys [x y]} :q} my-nested-hashmap]
  (println a b x y))
;; => A B X Y

(let [{:keys [a not-found b]} my-hashmap]
  (println a not-found b))
;; => A nil B

(let [{:keys [a not-found b], :or {not-found ":)"}} my-hashmap]
  (println a not-found b))
;; => A :) B

(let [{:keys [a b], :as all} my-hashmap]
  (println a b all))
;; => A B {:a A :b B :c C :d D}

(let [{:keys [a b not-found], :or {not-found ":)"}, :as all} my-hashmap]
  (println a b not-found all))
;; => A B :) {:a A :b B :c C :d D}
```

There are also `:strs` and `:syms` alternatives, for when your map has strings or symbols for keys (instead of keywords), respectively.

```clojure
(let [{:strs [a d]} {"a" "A", "b" "B", "c" "C", "d" "D"}]
  (println a d))
;; => A D

(let [{:syms [a d]} {'a "A", 'b "B", 'c "C", 'd "D"}]
  (println a d))
;; => A D
```

### Keyword arguments for function

Map destructuring also works with lists (but not vectors).

```clojure
(let [{:keys [a b]} '("X", "Y", :a "A", :b "B")]
(println a b))
;; => A B
```

This allows your functions to have optional keyword arguments.

```clojure
(defn foo [a b & {:keys [x y]}]
  (println a b x y))
(foo "A" "B")  ;; => A B nil nil
(foo "A" "B" :x "X")  ;; => A B X nil
(foo "A" "B" :x "X" :y "Y")  ;; => A B X Y
```

More examples
=============

Here be dragons.

**TODO**



## 01_compojure.md

Compojure Destructuring
=======================

[Official docs][compojure]

The basic syntax of a Compojure route is as follows:

```
(REQUEST-METHOD "/path" request-data (handler-fn request-data))
```

The simplest form is to pass the entire [Ring request map][ring-request-map] to your handler function.

Here, the request map will be bound to the `request` var and when `my-handler` is called it is passed as an argument. You can then extract/manipulate whatever data you want in the `request` map:


```clojure
(defn my-handler [request] ...)
;;                   ___________________
;;                  |                   V
(GET "/some/path" request (my-handler request))
```

Since the request is just a regular Clojure map, Compojure allows you to destructure that map using the same map destructuring methods above. This is useful if you only want specific values from the map for the handler function to work with:

```clojure
(defn my-handler [uri query-string]
  (println uri)
  ;; note that query-string is raw. e.g. "foo=bar&fizz=buzz"
  (println query-string))

(GET "/some/path" {uri :uri, query-string :query-string} (my-handler uri query-string))
;; or with the shortcuts
(GET "/some/path" {:keys [uri query-string]} (my-handler uri query-string))
```

If you want to pass the entire request map as well:

```clojure
(defn my-handler [uri query-string request]
  (println request)
  (println uri)
  (println query-string))

(GET "/some/path" {:keys [uri query-string] :as request} (my-handler uri query-string request))
```

Compojure-specific Destructuring
--------------------------------

Since regular destructuring can be quite verbose, Compojure offers a more specialised form of destructuring. If you supply a vector, Compojure will use this custom destructuring syntax [1][compojure].

**Note:** This only works for routes wrapped with the [`wrap-params` middleware][wrap-params]. This is because by default, a [Ring request map][ring-request-map] doesn't include the `:params` entry (which contains the parsed value of `:query-string`).
The `wrap-params` middleware parses the value of `:query-string` and inserts it into the request map as `:params`, which Compojure's special destructuring syntax targets. (TODO: Ring wiki is old. Is this still true?)

From here, assume that all handlers are wrapped with `wrap-params`.

(I got sleepy)

[compojure]: https://github.com/weavejester/compojure/wiki/Destructuring-Syntax
[ring-request-map]: https://github.com/ring-clojure/ring/wiki/Concepts#requests
[wrap-params]: https://github.com/ring-clojure/ring/wiki/Parameters

