# Functional Programming: Polling XHR

[View original Gist on GitHub](https://gist.github.com/Integralist/d84d65c909eeafe4dfca)

## fp polling xhr.js

```javascript
define(['module/bootstrap'], function(news) {
    var $ = news.$;

    function seconds(secs) {
        return secs * 1000;
    }

    function incrementDelay() {
        // don't want the delay to become ridiculous so I set a threshold I don't want it to go beyond
        if (delay < threshold) {
            delay *= 2;
        }

        return delay; // we return a value for the sake of our `dispatch` method called by `checkStatus`
    }

    function exists(value) {
        return value !== undefined && value !== null;
    }

    function dispatch() {
        var fns  = Array.prototype.slice.call(arguments, 0);
        var size = fns.length;

        return function(target) {
            var ret;

            for (var i = 0; i < size; i++) {
                var fn = fns[i];

                ret = fn.call(fn, target);

                if (exists(ret)) {
                    return ret; // break out of the loop as we've executed a function that doesn't return null
                }
            }

            return ret;
        };
    }

    var checkStatus = dispatch(
        function(status) { return status === 304 ? incrementDelay()   : null; },
        function(status) { return status === 200 ? delay = seconds(5) : null; }
    );

    function success(data, status, xhr) {
        checkStatus(xhr.status);
        poll(endpoint);
    }

    function failure(xhr, status, error) {
        incrementDelay();
        poll(endpoint);
    }

    function endpoint() {
        $.ajax('http://localhost:3000/foobar', { timeout: limit }).then(success, failure);
    }

    function poll(endpoint) {
        window.setTimeout(endpoint, delay);
    }

    var limit     = seconds(5);  // XHR timeout limit
    var delay     = seconds(5);  // timer delay
    var threshold = seconds(20); // timer delay threshold

    poll(endpoint);
});
```

