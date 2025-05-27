# Mori.js ClojureScript Data Structures in plain vanilla JavaScript -> http://swannodette.github.io/mori/ and https://github.com/swannodette/mori

[View original Gist on GitHub](https://gist.github.com/Integralist/e8b0adc5bb96a4162aea)

## mori.js

```javascript
var mori = require("mori");

var inc = function(n) {
  return n+1;
};

var v1 = mori.vector(1,2,3,4,5);
var v2 = mori.map(inc, v1);

console.log(v2.toString()); // => (2 3 4 5 6) 
console.log(mori.into_array(v1)); // => [1, 2, 3, 4, 5] 
console.log(mori.into_array(v2)); // => [2, 3, 4, 5, 6] 

console.log(
  mori.into_array(
    mori.map(inc, mori.vector(1,2,3,4,5))
  )
); // => [2, 3, 4, 5, 6] 
```

