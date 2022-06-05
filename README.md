# forum-security
This is an optional for the forum project in another repo.
It is/was hosted on GCP.
I have used the autocert library, to get the SSL certificates from Let's Encrypt.
I have also implemented a rate limiter (not allowed to use exteranl library), to limit the no. of request to the server in a given time interval.
The rate limiter functions by first filling with 3 tokens, then it refills one by every 0.2s (or longer for demo purpose).
The passwords are hashed, and a uuid is used for the session cookies (from forum).
There are other optionals for forum, and I am responsible to do the forum-security optional.

forum-security
Objectives
You must follow the same principles as the first subject.

For this project you must take into account the security of your forum.

You should implement a Hypertext Transfer Protocol Secure (HTTPS) protocol :

Encrypted connection : for this you will have to generate an SSL certificate, you can think of this like a identity card for your website. You can create your certificates or use "Certificate Authorities"(CA's)

We recommend you to take a look into cipher suites.

The implementation of Rate Limiting must be present on this project

You should encrypt at least the clients passwords. As a Bonus you can also encrypt the database, for this you will have to create a password for your database.

Sessions and cookies were implemented in the previous project but not under-pressure (tested in an attack environment). So this time you must take this into account.

Clients session cookies should be unique. For instance, the session state is stored on the server and the session should present an unique identifier. This way the client has no direct access to it. Therefore, there is no way for attackers to read or tamper with session state.
Hints
You can take a look at the openssl manual.
For the session cookies you can take a look at the Universal Unique Identifier (UUID)
Instructions
You must handle website errors, HTTPS status.
You must handle all sort of technical errors.
The code must respect the good practices.
It is recommended to have test files for unit testing.
Allowed packages
All standard Go packages are allowed.
sqlite3
bcrypt
UUID
autocert
This project will help you learn about :

HTTPS
Cipher suites
Goroutines
Channels
Rate Limiting
Encryption
password
session/cookies
Universal Unique Identifier (UUID)
