# JavaScript Function Programming (scratch pad) -> Most of the code here is modified from the excellent O'Reilly book "Functional JavaScript".

## 02. ES5.js

```javascript
var csv = 'name,age,hair\nmark,32,brown\ncat,27,red';

function parseCSV(csv) {
    var data = csv.split('\n');
        data.unshift([]); // accumulator

    // ES5's `reduce` method isn't as nice to use as found in Underscore/Lo-Dash
    // This is because it uses the first item in the collection as the accumulator 
    // So really the `current` argument starts at index 1 and `previous` is index 0
    // Underscore/Lo-Dash works the same but also allows you to pass in an accumulator
    // So to make things easier/clearer I decided to inject my own accumulator into the collection
    return data.reduce(
        function(previous, current, index, collection) {
            previous.push(current.split(',').map(function(item) {
                return item.trim();
            })); // `push` returns the new length of the array rather than the array itself...

            return previous; // ...hence we need to explicitly return the array
        }
    );
}

parseCSV(csv) // => [ [ 'name', 'age', 'hair' ],
              //    [ 'mark', '32', 'brown' ],
              //    [ 'cat', '27', 'red' ] ]
```

## 03. LoDash.js

```javascript
// Note: using Lo-Dash/Underscore's `_.reduce` is a lot cleaner and easier code to use
// By that I mean we don't need to manipulate our data structure to allow for an accumulator to be injected

var _ = {
    reduce: require('lodash.reduce'), // npm install lodash.reduce
    map:    require('lodash.map')     // npm install lodash.map
};

function _parseCSV(csv) {
    return _.reduce(csv.split('\n'), function(accumulator, item) {
            accumulator.push(_.map(item.split(','), function(item) {
                return item.trim();
            }));

            return accumulator;
        }, []);
}

_parseCSV(csv) // => [ [ 'name', 'age', 'hair' ],
               //    [ 'mark', '32', 'brown' ],
               //    [ 'cat', '27', 'red' ] ]

```

## 05. construction.js

```javascript
var _ = require('lodash'); // lodash.tap and lodash.chain weren't available on npm as separate modules (update: there is a Lo-Dash npm module that allows for pulling in specific functions)

function lyricSegment(number) {
    return _.chain([])
            .push(number + ' bottles of beer on the wall')
            .push(number + ' bottles of beer')
            .push('Take one down, pass it around')
            .tap(function(lyrics) {
                if (number > 1)
                    lyrics.push((number - 1) + ' bottles of beer on the wall.')
                else
                    lyrics.push('No more bottles of beer on the wall!')
            })
            .value()
}

function song(start, end, generator) {
    return _.reduce(_.range(start, end, -1), function(accumulator, number) {
        return accumulator.concat(generator(number))
    }, [])
}

song(99, 0, lyricSegment)
```

## 09. Rules of Recursion.md

```markdown
- Know when to stop
- Decide how to take one step
- Break the problem down into a smaller problem

```js
function myLength(arr) {
    if (_.isEmpty(arr)) {
        return 0;
    } else {
        return 1 + myLength(_.rest(arr));
    }
}

myLength(_.range(10));     // 10
myLength(_.range(1000));   // 1000
myLength(_.range(10000));  // 10000
myLength(_.range(100000)); // FATAL ERROR: JS Allocation failed - process out of memory

// See http://www.integralist.co.uk/posts/understanding-recursion-in-functional-javascript-programming/
// to understand tail-call optimisation
```

- Know when to stop (`if (_.isEmpty(arr)) {`)
- Decide how to take one step (`1 + myLength(...)`)
- Break the problem down into a smaller problem (`myLength(_.rest(arr)`)
```

## 10. Recursive Checking and Validating.js

```javascript
var _ = require('lodash'); // lodash.toArray and lodash.isEmpty weren't available on npm as separate modules (update: there is a Lo-Dash npm module that allows for pulling in specific functions)

function isEven(item) {
    return item % 2 === 0 ? true : false;
}

function isOdd(item) {
    return item % 2 === 1 ? true : false;
}

function isZero(item) {
    return item === 0 ? true : false;
}

function valid() {
    var predicates = _.toArray(arguments);

    return function() {
        var args = _.toArray(arguments);

        // "truth" is whatever we want to return,
        // if all the arguments provided to the above return function
        // pass each predicate's requirements
        var everything = function(preds, truth) {
            // If there are no preds
            // (or no more since the recursive calls to this function)
            // then we've safely validated our "args" so we can return our "truth"
            if (_.isEmpty(preds)) {
                return truth;
            } else {
                // recursiveness happens here...
                // if all the args pass the first pred's requirements
                // then call "everything" again, but now reduce the problem
                // down to the remaining preds the args have yet to be tested against.
                return _.every(args, _.first(preds)) && everything(_.rest(preds), truth);

                // Note: we have short-circuited the recursiveness
                // by first making sure the arguments pass the first set of requirements
                // e.g. whatever the first predicate function that was provided
            }
        };

        return everything(predicates, true);
    };
}

var evenNumbers = valid(_.isNumber, isEven);

console.log(
    evenNumbers(1, 2)
); // => false

console.log(
    evenNumbers(2, 4, 6)
); // => true

// We can change the functionality to allow some items to pass by
// setting our truth to be "false" while using _.some instead of _.every
// and also chaning our condition from && to ||
// See below example...

function possiblyValid() {
    var predicates = _.toArray(arguments);

    return function() {
        var args = _.toArray(arguments);

        var something = function(preds, truth) {
            if (_.isEmpty(preds)) {
                return truth;
            } else {
                return _.some(args, _.first(preds)) || something(_.rest(preds), truth);
            }
        };

        return something(predicates, false);
    };
}

var zeroOrOdd = possiblyValid(isOdd, isZero);

console.log(
    zeroOrOdd()
); // => false

console.log(
    zeroOrOdd(0, 2, 4, 6)
); // => true

console.log(
    zeroOrOdd(2, 4, 6)
); // => false

console.log(
    zeroOrOdd(2, 4, 6, 7)
); // => true
```

## 01. composability.js

```javascript
// Higher order function (takes in a function and returns a function)
function comparator(predicate) {
    return function(x, y) {
        if (predicate(x, y)) {
            return -1;
        } else if (predicate(y, x)) {
            return 1;
        } else {
            return 0;
        }
    };
}

function lessOrEqual(x, y) {
    return x <= y;
}

// Composability
[2, 3, -1, -6, 0, -108, 42, 10].sort(comparator(lessOrEqual)) // => [-108, -6, -1, 0, 2, 3, 10, 42]
```

## 07. Null object pattern.js

```javascript
var _ = {
    reduce: require('lodash.reduce'), // npm install lodash.reduce
    rest:   require('lodash.rest'),   // npm install lodash.rest
    map:    require('lodash.map')     // npm install lodash.map
};

function exists(value) {
    return value !== undefined && value !== null
}

function fnull(fn) {
    var defaults = _.rest(arguments); // default values to use in case a null value is passed in

    return function() {
        /*
            Modify arguments so null values are replaced with default values
            e.g. args passed in our example below:
                 1, 2 (1 * 2 == 2)
                 2, 3 (2 * 3 == 6)
                 6, null (null replaced with 1 and so this becomes 6 * 1 == 6)
                 6, 5 (6 * 5 == 30)
        */
        var args = _.map(arguments, function(item, index) {
            return exists(item) ? item : defaults[index];
        });

        return fn.apply(null, args);
    };
}

var data = [1, 2, 3, null, 5];
var multiplier   = function(total, n) { return total * n }; 
var safeMultiply = fnull(multiplier, 1, 1);

_.reduce(data, multiplier);   // => 0
_.reduce(data, safeMultiply); // => 30
```

## 14.Pipeline.js

```javascript
var _ = require('lodash');

function rev(arr) {
    return _.chain(arr)
            .reverse()
            .value();
}

function pipeline(seed) {
    return _.reduce(_.rest(arguments), function(l, r) {
        return r(l);
    }, seed);
}

var data = [2, 3, null, 1, 42, false];

console.log(
    pipeline(data,
            _.compact,
            _.initial,
            _.rest,
            rev)
); // => [1, 3]

// Evaluates to...

/*
rev(
    _.rest(
        _.initial(
            _.compact(data))))
*/
```

## 12. Composition can help avoid mutation.js

```javascript
var _ = require('lodash');

// Make an Array out of the provided arguments
function construct() {
    return _.reduce(_.toArray(arguments), function(accumulator, value) {
        return _.chain(accumulator)
                .push(value)
                .value();
    }, []);
}

function merge() {
    return _.extend.apply(null,
                          construct.apply(
                              construct, _.flatten([{}, arguments])
                          )
                         );
}

var person = { name: "Mark" };

// We return an object that is a combination of the provided objects
console.log(
    merge(person, { age: 32 }, { location: "London" })
); // => { name: "Mark", age: 32, location: "London" }

// But we don't mutate the original object
console.log(
    person
); // => { name: "Mark" }
```

## 13. If statements in Functional form.js

```javascript
var _ = require('lodash');

function exists(value) {
    return value !== undefined && value !== null;
}

function dispatch() {
    var fns  = _.toArray(arguments);
    var size = fns.length;

    return function(target) {
        var ret;

        for (var i = 0; i < size; i++) {
            var fn = fns[i];

            ret = fn.call(fn, target);

            if (exists(ret)) return ret;
        }

        return ret;
    };
}

// The `dispatch` function executes each provided function
// until a non-undefined value is returned.
// Effectively action like a batch of if statements
var convertToString = dispatch(
    function(s) { return _.isString(s) ? s : undefined; },
    function(s) { return _.isObject(s) ? JSON.stringify(s) : undefined; },
    function(s) { return s.toString(); }
);

console.log(
    convertToString({ foo: 'bar' })
);

console.log(
    convertToString(123)
);

console.log(
    convertToString('abc')
);
```

## 04. conditionals.js

```javascript
var _ = require('lodash'); // lodash.tap wasn't available on npm as separate modules (update: there is a Lo-Dash npm module that allows for pulling in specific functions)

function exists(value) {
    return value !== undefined && value !== null
}

function truthy(value) {
    return value !== false && exists(value)
}

// The following code abstracts away the ugliness of using conditionals to say:
// "do this thing, only if this other thing exists"
// The use of `if (!!condition)` would be more concise than having two separate functions
// like we use above (`exists` and `truthy`) but we're abiding by FP principle of:
// "replace values with functions"

function doWhen(condition, action) {
    if (truthy(condition))
        return action()
    else
        return
}

function executeIfFieldExists(target, field) {
    return doWhen(target[field], function() {
        return _(target).tap(function(target) {
            console.log('The result is', _.result(target, field))
        })
    })
}

executeIfFieldExists({ foo: 'bar' }, 'foo');   // => 'bar'
executeIfFieldExists([1, 2, 3, 4], 'reverse'); // => [4, 3, 2, 1]
```

## 08. Checking and Validating data.js

```javascript
var _ = require('lodash'); // lodash.tap and lodash.chain weren't available on npm as separate modules (update: there is a Lo-Dash npm module that allows for pulling in specific functions)

// Variation of a no-op that
// always returns the value specified
function always(value) {
    return function() {
        return value;
    };
}

// This function requires that all provided "validators"
// return an object that has a "message" property.
// This is so the "errs" accumulator Array will hold relevant error messages
function checker() {
    var validators = _.toArray(arguments);

    return function(obj) {
        return _.reduce(validators, function(errs, check) {
            if (check(obj)) {
                return errs;
            } else {
                // We use the "chain" method much like "tap" in Ruby.
                // The "chain" method needs "value" to be called
                // otherwise it'll always return a chain obj until
                // the final value is requested
                return _.chain(errs).push(check.message).value();
            }
        }, []);
    };
}

var alwaysPasses = checker(always(true), always(true)); // we can have as many "validators" as we want

console.log(
    alwaysPasses({})
); // => [] no errors

// We'll make sure our next example fails,
// so we need a "message" property on the "validator" we provide
var fails = always(false);
    fails.message = "This thing failed";

var alwaysFails = checker(fails);

console.log(
    alwaysFails({})
); // => ["This thing failed"]

// We want to be careful adding a "message" property to objects we don't own
// The following function automates this for us
function validator(message, fn) {
    var f = function() {
        // Call the original function
        // such as _.isObject and pass itself as the "this" value
        // and pass through the original argument(s) to that function
        // e.g. var test = validator("Error!", _.isObject) will mean when we call
        // test([]) then test.message == "Error!"
        return fn.apply(fn, arguments);
    };

    f.message = message;

    return f;
}

var objectCheck = checker(validator("object validation fail", _.isObject));
var arrayCheck = checker(validator("array validation fail", _.isArray));

console.log(
    objectCheck(123)
); // => ["object validation fail"]

console.log(
    arrayCheck(123)
); // => ["array validation fail"]

console.log(
    arrayCheck([])
); // => []
```

## 11. Definition of Purity.md

```markdown
A pure function adheres to the following properties:

- Its result is calculated only from the values of its arguments
- It cannot rely on data that changes external to its control
- It cannot change the state of something external to its body
```

## 0. JavaScript Function Programming TOC.md

```markdown
This code is modified from the excellent O'Reilly book "Functional JavaScript". You should buy it, I highly recommend it! Don't kid yourself into thinking this gist even remotely covers the great content from a 200+ page technical book on the subject; it doesn't. Buy the book and get the in-depth knowledge for yourself. It's worth it.

- [Composability through higher order functions](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce#file-01-composability-js)
- [Parsing CSV (ES5 version)](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce/#file-02-es5-js)
- [Parsing CSV (LoDash version)](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce/#file-03-lodash-js)
- [Abstract away conditionals](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce/#file-04-conditionals-js)
- [Data construction](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce/#file-05-construction-js)
- [Currying and Partial Application](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce/#file-06-currying-and-partial-application-js)
- [Null Object pattern](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce/#file-07-null-object-pattern-js)
- [Checking and Validating data](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce#file-08-checking-and-validating-data-js)
- [Rules of Recursion](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce#file-09-rules-of-recursion-md)
- [Recursive Checking and Validating](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce#file-10-recursive-checking-and-validating-js)
- [Definition of Purity](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce#file-11-definition-of-purity-md)
- [Composition can help avoid mutation](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce#file-12-composition-can-help-avoid-mutation-js)
- [If statements in Functional form](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce#file-13-if-statements-in-functional-form-js)
- [Pipeline](https://gist.github.com/Integralist/b7ed2e337a0b5cbbecce#file-14-pipeline-js)
```

## 06. Currying and Partial Application.js

```javascript
// Code copied from http://eliperelman.com/fn.js/

var fn = {};

fn.toArray = function (collection) {
	return [].slice.call(collection);
};

fn.cloneArray = fn.toArray;

fn.apply = function (handler, args) {
	return handler.apply(null, args);
};

fn.concat = function () {
	var args = fn.toArray(arguments);
	var first = args[ 0 ];

	if (!fn.is('array', first) && !fn.is('string', first)) {
		first = args.length ? [ first ] : [ ];
	}

	return first.concat.apply(first, args.slice(1));
};

fn.type = function (value) {
	// If the value is null or undefined, return the stringified name,
	// otherwise get the [[Class]] and compare to the relevant part of the value
	return value == null ?
		'' + value :
		{ }.toString.call(value).slice(8, -1).toLowerCase();
};

fn.is = function (type, value) {
	return type === fn.type(value);
};

fn.partial = function () {
	var args = fn.toArray(arguments);
	var handler = args[0];
	var partialArgs = args.slice(1);

	return function () {
		return fn.apply(handler, fn.concat(partialArgs, fn.toArray(arguments)) );
	};
};

fn.identity = function (arg) {
	return arg;
};

fn.reverse = function (collection) {
	return fn.cloneArray(collection).reverse();
};

var currier = function makeCurry(rightward) {
	return function (handler, arity) {
		if (handler.curried) {
			return handler;
		}

		arity = arity || handler.length;

		var curry = function curry() {
			var args = fn.toArray(arguments);

			if (args.length >= arity) {
				var transform = rightward ? 'reverse' : 'identity';
				return fn.apply(handler, fn[ transform ](args));
			}

			var inner = function () {
				return fn.apply(curry, args.concat(fn.toArray(arguments)));
			};

			inner.curried = true;

			return inner;
		};

		curry.curried = true;

		return curry;
	};
};

fn.curry = currier(false);

fn.curryRight = currier(true);

// Partial Application

var fullName = function (firstName, lastName) {
	return firstName + ' ' + lastName;
};

var billName = fn.partial(fullName, 'Bill');

billName('Smith'); // "Bill Smith"
billName('Clinton'); // "Bill Clinton"

// Curry

var fullName2 = fn.curry(function (firstName, middleName, lastName) {
	return firstName + ' ' + middleName + ' ' + lastName;
});

var billName2 = fullName2('Bill');

billName2('Damon')('Smith'); // "Bill Damon Smith"
billName2('Jefferson', 'Clinton'); // "Bill Jefferson Clinton"
fullName2('Jenn', 'Anne', 'Cochran'); // "Jenn Anne Cochran"
fullName2('Jenn', 'Anne')('Cochran'); // "Jenn Anne Cochran"

// Curry right

var fullName3 = fn.curryRight(function (firstName, middleName, lastName) {
	return firstName + ' ' + middleName + ' ' + lastName;
});

var smithName = fullName3('Smith');

smithName('Damon')('Bill'); // "Bill Damon Smith"
smithName('Jefferson', 'Bill'); // "Bill Jefferson Clinton"
fullName3('Cochran', 'Anne', 'Jenn'); // "Jenn Anne Cochran"
fullName3('Cochran', 'Anne')('Jenn'); // "Jenn Anne Cochran"
```

