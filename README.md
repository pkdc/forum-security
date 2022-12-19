# forum-security
This is an optional for the forum project in another repo.

It was hosted on GCP. Now hosting on heroku https://forum-secu.herokuapp.com

I have used the autocert library to get the SSL certificates from Let's Encrypt, and used the TLS cipher suites to encrypted the traffic between clients/servers.

I have also implemented a rate limiter, to limit the no. of request to the server in a given time interval.
The rate limiter functions by first filling with 3 tokens, then it refills one by every 0.2s.

The passwords are hashed, and a uuid is used for the session cookie.

There are other optionals for forum, and I am responsible to do the forum-security optional.
