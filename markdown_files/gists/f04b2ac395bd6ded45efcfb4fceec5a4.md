# [AWS Amplify Tips] #aws #amplify #cognito #js #javascript

[View original Gist on GitHub](https://gist.github.com/Integralist/f04b2ac395bd6ded45efcfb4fceec5a4)

## AWS Amplify Tips.md

## Enable debug log

```js
window.LOG_LEVEL = 'DEBUG';
```

## Disable analytics 

The noise in the debug logs is insane:

```js
import { Analytics } from 'aws-amplify';

Analytics.disable();
```

## Manual Configuration

```js
import Amplify,{Auth} from 'aws-amplify';

Amplify.configure({
  Auth: {
    identityPoolId: '...',
    region: '...',
    userPoolId: '...',
    userPoolWebClientId: '...',
  }
});
```

## Sign Up

```js
const phone = this.state.phone.trim();

var attributes = {
    email: email,
    name: fullname
}

if (phone) {
    attributes['phone_number'] = phone;
}

Auth.signUp({
    username: username,
    password: password,
    attributes: attributes
})
.then(
    this.setState(() => {
        return {
            enterAuth: true
        }
    })
)
.catch( err => {
    console.log(username, password, email, phone, fullname);
    console.log(err);
})
```

