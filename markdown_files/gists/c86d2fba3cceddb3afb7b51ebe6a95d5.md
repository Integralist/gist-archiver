# [Single Sign-On SSO] 

[View original Gist on GitHub](https://gist.github.com/Integralist/c86d2fba3cceddb3afb7b51ebe6a95d5)

**Tags:** #sso #singlesignon #auth #multipledomains #aws

## Single Sign-On SSO].md

> These are all suggestions from AWS Technical Support

## Initial suggestion

- You have a "Master" domain, let's say "auth.ronan.com". 
- You have 2 websites, a.com and b.com. Of course, your end user. Let's call him "Bob". 
- When Bob makes a request to a.com, we will check if they have sent a token. If not, we will redirect them to auth.ronan.com to handle the auth, then return a token along with a cookie for said domain in the form of a redirect. 
- Next, Bob sends a request to b.com without a token, it again would need to redirect to the "auth" website to obtain the token and set the cookie for b.com too. 
- Now, your 2 sites would have a token which can be used to make API calls. 

## Follow-up Discussion

I've continued to work with our internal Cognito experts here, where I have the following clarification: 

- AWS User Pools does not limit the usage of One User Pool per domain, however the way cookie/local storage works prevents cross site access, which in turn really doesn't work in our favour in this situation.  This isn't something we can work around from an AWS perspective. 

In this case, the handling of this logic needs to be done on the client side. If you build your applications so they do not rely on the cookie and the localStorage, then it's possible to use one user pool to authenticate users for multiple domains applications. 

Our internal Cognito SME noted that there are several possible options he could can think of. For example, store the tokens in memory, passing the id_tokens in the http request's header when sending requests to your own applications. 

Another option would be using the User Pool as an Oauth SSO.  Let's consider an ALB as "the application".  When we go to the ALB URL, we authenticate with a User Pool Hosted UI. In this case, you can have multiple ALBs using one User Pool as their authoriser. This also requires some application side implementation. 

The auth flow would look something like this, and I hope it improves on my previous illustration:

1. User accesses application A
2. Application A checks user login status against its own storage/cache. 
3. Redirect user to hosted UI so they can login if step 2 found the user to be in the status of not logged in.
4. User authenticates in Hosted UI and gets redirected back to application A.

5. User accesses application B
6. Application B checks user login status against its own storage/cache (like step 2). 
7. Redirect user to hosted UI so they can login if step 5 found the user to be in the status of not logged in.
   7a. user already signed-in to the User Pool (step3.), directly redirect back to application B with tokens 

The outlined flow provided by our local SME is one using IdP SSO for multiple applications. If you can consider to use the Cognito Hosted UI, I think this would also remove some of the heavy lifting. 

