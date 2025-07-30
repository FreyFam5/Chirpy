# 🐦 Chirpy

An http server api for a very real (_Not real^TM_) site called Chirpy!

## What it does

This api aims to add a number of database based elements to allow for easy & safe saving of user data!
This includes:

- Allowing user log in to a server[^1] with an email and password
- Timeout tokens used for authentication
- Chirp posting [^2]
- Chirp finding [^3]
- Chirp deletion [^4]
- Hit counting [^5]
- ...

This was a guided project for [boot.dev](https://www.boot.dev).

[^1]: This server will automatically start on a local serer
[^2]: Sending messages to a local server
[^3]: Retrieving a list of the messages sent
[^4]: Deleting a specific message
[^5]: The amount of times someone has opened the main page

## ⚙️ Installation

Inside a Go module:

```bash
go get https://github.com/FreyFam5/Chirpy
```

## 🪧 Use

After installing, you'll need to make a separate `.env` file for certain security aspects. The items you'll need in this `.env` file include (_the values of each are surrounded by quotes_):

- `DB_URL=<url of your servers database>`
- `PLATFORM=<the platform that the api will see you on, use "dev" for admin permissions>`
- `SECRET=<a random 64 character string, used for access tokens signing>`
- `POLKA_KEY=<your own polka key here>`

After your `.env` file is done, you should be able to start your server! This server is meant to be used with a local sever of `http://localhost:8080/` on start. Make sure your terminal is in the chirpy directory then run:

```bash
go run .
```

Chirpy can be used with a few requests, simply start the url with your server's url (for me it is `http://localhost:8080/`) while the server is running.
Heres a list of all of them:

- `POST` requests:

  - `/api/chirps` : Makes a chirp message
    _Takes a json `"body"` key with a message and responds with a full `chirp`_
    Ex:

    ```json
    {
    	"body": "This is an example message!"
    }
    ```

  - `/api/users` : Makes a new user
    _Takes a json `"email"` and `"password"` key and responds with a full `user`_
    Ex:

    ```json
    {
    	"email": "example@chirpy.com",
    	"password": "your_password"
    }
    ```

  - `/api/login` : Logs into the given user
    _Takes a json `"email"` and `"password"` keys and responds with a full `user`, a `access token string`, and a `refresh token string`_
    Ex:

    ```json
    {
    	"email": "example@chirpy.com",
    	"password": "your_password"
    }
    ```

  - `/api/polka/webhooks` : Updates the users email and password, based on who's logged in
    _Takes a json `"event"` key with `"user.upgraded"` as it's value (anything else will not allow it) and `"data": {"user_id"}` key, responds with nothing_
    Ex:

    ```json
    {
    	"event": "user.upgraded",
    	"data": {
    		"user_id": "<user_id>"
    	}
    }
    ```

  - `/api/refresh` : Refreshes the current user's access token (These go out every 30 min)
    _Takes the header's `access token`, responds with a new `access token string`_
  - `/api/revoke` : Revokes a user's refresh token, forcing them to login again
    _Takes the header's `refresh token`, responds with nothing_
  - `/admin/reset` : If you have `dev` permissions, will reset all data in database
    _Takes the platform, responds with a confirmation message on success_

- `GET` requests:

  - `/api/healthz` : Checks if the server is ready
    _Takes nothing, responds with `OK` on success_
  - `/api/chirps` : Gets all the chirps in the database, there is an optional query that can be used for a specific users chirps to be listed (`/api/chirps?author_id=<user_id>`)
    _Takes an optional `author_id` query, responds with a list of chirps_
  - `/api/chirps/{chirpID}` : Gets a specific chirp, just replace the `{chirpID}` with the wanted chirp ID
    _Takes a chirp id in the url, responds with the specified chirp if it exists_
  - `/admin/metrics` : If you are an admin, will show how many time the server's `/api` url's have been visited
    _Takes nothing, responds with how many times the server's api url has been visited_

- `PUT` requests:

  - `/api/users` : Updates the current users email and password
    _Takes a json `"email"` key and `"password"` key, responds with the newly updated user_

- `DELETE` requests:

  - `/api/chirps/{chirpID}` : Deletes the current users chirp based on the chirp id given
    _Takes a chirp id in the url, responds with nothing_
