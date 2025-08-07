# vpub-plus

Simple message board software and also a fork.

## Examples

<table>
  <tr>
    <td>
      <a href="https://github.com/user-attachments/assets/4ded261a-f2c3-4f8f-b474-be268aa61ff7">
        <img alt="Status Cafe Forum - vpub instance" width="640" src="https://github.com/user-attachments/assets/4ded261a-f2c3-4f8f-b474-be268aa61ff7">
      </a>
      <p><a href="https://forum.status.cafe">Status Cafe Forum - vpub instance</a></p>
    </td>
    <td>
      <a href="https://github.com/user-attachments/assets/6e65c795-6ef6-40f3-b222-18f4c3f48548">
        <img alt="Vpub Plus Forum - vpub-plus instance" width="640" src="https://github.com/user-attachments/assets/6e65c795-6ef6-40f3-b222-18f4c3f48548">
      </a>
      <p><a href="https://vpub.mysh.dev">Vpub Plus Forum - vpub-plus instance</a></p>
    </td>
  </tr>
</table>

## Installation

### Using `Docker` and `Docker Compose`

The easiest way to get started is with `Docker` and `Docker Compose`.

1. **Clone the repository**
   ```bash
   git clone https://github.com/hugmouse/vpub-plus.git
   cd vpub-plus
   ```

2. **Configure your environment**
   Copy the example `.env.example` file to `.env` and edit it to your liking.
   ```bash
   cp .env.example .env
   ```
   The default values in `.env.example` are a good starting point for local testing.

3. **Run the application**
   ```bash
   docker-compose up -d
   ```

And then you can navigate to `localhost:1337` (or whatever `HOST_PORT` you set in `.env`) to see your very own forum!

### Compiling vpub-plus from the source

To host it and install it you have to have:

* Golang
* Postgresql
* Git
* Make
* Systemd (optional)

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

`vpub-plus` is configured using environment variables. If you are using Docker, you can set these in the `.env` file. If you are running from source, you can set them in your shell or use a tool like `direnv`.

Here are the available variables:

| Variable                          | Description                                                                                                | Default                               |
| --------------------------------- | ---------------------------------------------------------------------------------------------------------- | ------------------------------------- |
| `PORT`                            | The port the application will listen on.                                                                   | `8080`                                |
| `HOST_PORT`                       | The port on the host machine to map to the application port (for Docker).                                  | `1337`                                |
| `DATABASE_URL`                    | The full connection URL for your PostgreSQL database.                                                      | `postgres://vpub:yourpassword@db:5432/vpub?sslmode=disable` |
| `POSTGRES_USER`                   | The PostgreSQL user.                                                                                       | `vpub`                                |
| `POSTGRES_PASSWORD`               | The PostgreSQL password.                                                                                   | `yourpassword`                        |
| `POSTGRES_DB`                     | The PostgreSQL database name.                                                                              | `vpub`                                |
| `SESSION_KEY`                     | A 32-byte random string for session authentication. **Change this!**                                       | `your32byteslongsessionkeyhere`       |
| `CSRF_KEY`                        | A 32-byte random string for CSRF protection. **Change this!**                                              | `your32byteslongcsrfkeyhere`          |
| `CSRF_SECURE`                     | Set to `true` if you are using HTTPS to make CSRF cookies secure.                                            | `true`                                |
| `TITLE`                           | The title of your forum.                                                                                   | `My vpub-plus forum`                  |
| `PROXYING_ENABLED`                | Set to `true` if you are running behind a reverse proxy.                                                   | `true`                                |
| `POSTGRES_MAX_OPEN_CONNECTIONS`   | The maximum number of open connections to the database.                                                    | `0` (unlimited)                       |
| `POSTGRES_MAX_IDLE_CONNECTIONS`   | The maximum number of idle connections to the database.                                                    | `0` (unlimited)                       |
| `POSTGRES_MAX_LIFETIME`           | The maximum amount of time a connection may be reused.                                                     | `5m`                                  |

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
Environment=PORT=1337

# IMPORTANT: Those keys should be 32 bytes long and random
Environment=SESSION_KEY=CHANGE ME
Environment=CSRF_KEY=CHANGE ME

# If you are going to use HTTPS, then use secure cookies
Environment=CSRF_SECURE=true
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

[example-vpub]: https://forum.status.cafe/
[example-vpub-plus]: https://vpub.mysh.dev/
