# vpub-plus

Simple message board software and also a fork.

## Examples

[Status Cafe Forum][example-vpub]

<img alt="An example of a forum" width="640" src="https://user-images.githubusercontent.com/44648612/176660262-f386251f-12f2-4e97-9e6e-c540640728bf.png">

[Turkey forum][example-mysh]

<img alt="An example of a forum, but in Turkey!" width="640" src="https://user-images.githubusercontent.com/44648612/178062763-0d0978f9-6604-4d0e-b55f-e040218ecbf7.png">

## Installation

To host it and install it you have to have:

* Golang
* Postgresql
* Git
* Make
* Systemd (optional)

### Compiling vpub-plus from the source

Here is how to build vpub:

1. `git clone https://github.com/hugmouse/vpub-plus.git`
2. `cd vpub-plus`
3. `make`

You should now have `vpub` in `./bin/`!

### Creating a `vpub` user

For isolation purposes, we can create a user that is going to run a `vpub` instance

* `useradd vpub`

Make sure that `vpub` group exists too! And if it does not, then:

* `groupadd vpub`

### Database setup

Make sure that you have [postgresql][postgres] installed!

* Create a new database: `createdb vpub` (or create it from `psql`)

### Set up environment variables

Now you have to set those environment variables:

* `DATABASE_URL` - [Postgresql connection URL][postgres-url-format]
* `SESSION_KEY` - 32 bytes long session key
* `CSRF_KEY` - 32 bytes longs CSRF key
* `CSRF_SECURE` - Makes CSRF cookies secure (`true`/`false`)
* `TITLE` - Your forum name
* `PORT` - What port is going to be used by a `vpub` HTTP server

You can check the example configuration in systemd config!

----

At this point you can run `vpub` just fine, other steps are optional ones

### Systemd config (optional)

Create a `/etc/systemd/system/vpub.service` file and add this example config to there:

```
[Install]
WantedBy=multi-user.target

[Unit]
Description="Message board"
Documentation="https://github.com/hugmouse/vpub-plus"

[Service]
ExecStart=/usr/local/bin/vpub
User=vpub
Group=vpub

Environment=DATABASE_URL=postgres://vpub@127.0.0.1/vpub?sslmode=disable

# Default port 8080
Environment=PORT=1337

# Those keys should be 32 bytes long
Environment=SESSION_KEY=CHANGE ME
Environment=CSRF_KEY=CHANGE ME

# If you are going to use HTTPS, then use secure cookies
Environment=CSRF_SECURE=true

Environment=TITLE="My very own forum"
```

After that you just can run it like any other systemd service: `systemctl enable --now vpub`

If something goes wrong, you can use `journalctl -eu vpub` to troubleshoot this service

### Where to go next

At this point a `vpub` service should be running on a 1337 port without HTTPS.

* You can add a reverse proxy, like `NGINX`, to handle a secure connection
* You can install a `certbot` that is going to create a let's encrypt cert for your domain

## Credentials

On the first run `vpub` will create an admin user with the password "admin".
Log in and change it by navigating to `/admin/users` route.

## Registering new users

To register you have to have a unique key! Create one on `/admin/keys` page and use it to create a new user account.

## Atom feed

You need to change the `URL` in settings to pass the atom feed check.

Specify your domain name in `/admin/settings/edit` like this: `https://example.com/`.

[postgres]: https://www.postgresql.org/download/

[postgres-url-format]: https://stackoverflow.com/q/3582552

[example-mysh]: https://turkey.mysh.dev/

[example-vpub]: https://forum.status.cafe/
