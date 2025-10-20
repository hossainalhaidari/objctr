# objctr

A minimalistic object storage with that is easy to setup and a very simple to use.

## Quick Start

Clone the project, create a YAML file with the following content and save it as `objctr.yml` in the root directory of the project:

```yaml
port: 3000 # The port your app will run on
path: /data # Path to root directory of your object storage in your local computer
users: # List of users that can access the API
  - key: default # For all unauthorized users
    read:
      - /
  - key: 5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8 # For an authorized user (SHA-256 for the text `password`)
    write:
      - /private
```

Run `go run .` to start the server locally at `http://localhost:3000` and use the following requests to use the API:

```sh
# List directory content
curl -X GET http://localhost:3000

# Create a directory
curl -X POST http://localhost:3000/public -H "Authorization: password"

# Upload a file
curl -X POST http://localhost:3000/public/file.txt -H "Authorization: password" -F "file=@/path/to/file.txt"

# Copy a file/directory to another path
curl -X PUT http://localhost:3000/public/file.txt?to=file2.txt -H "Authorization: password"

# Move a file/directory to another path
curl -X PATCH http://localhost:3000/public/file.txt?to=file3.txt -H "Authorization: password"

# Remove a file/directory
curl -X DELETE http://localhost:3000/file2.txt -H "Authorization: password"
```

You can also use the shell function in [objctr.sh](./objctr.sh) file to call your API using CLI:

```sh
# List directory content
$ objctr ls

# Create a directory
$ objctr mk public

# Upload a file
$ objctr mk public/file.txt /path/to/file.txt

# Copy a file/directory to another path
$ objctr cp public/file.txt file2.txt

# Move a file/directory to another path
$ objctr mv public/file.txt file3.txt

# Remove a file/directory
$ objctr rm file2.txt
```

## Configuration

### File location

The YAML configuration file should be saved in one of the locations below (the list is in order and the first file found will be used and the rest will be ignored):

- In any location, but provided as an argument: `go run . /path/to/config.yml`
- As `objctr.yml` in the same directory where the app is running
- In your home folder as `objctr.yml`

### Authentication and Authorization

The list of users in the YAML file are structured as follows:

```yaml
- key: <HASH> # SHA-256 of the user's password
  read: # An array of directories that this user can read from
    - <PATH1>
    - <PATH2>
    - ...
  write: # An array of directories that this user can write to
    - <PATH1>
    - <PATH2>
    - ...
```

The authorized users can access the API by providing `Authorization: PASSWORD` as a header (replace `PASSWORD` with the user's actual password). However in the config file, the hashed version of the password is stored for security reasons. To create a SHA-256 of a password you can use the helper command `go run . hash YOUR_PASSWORD` (replace `YOUR_PASSWORD` with your own).

The first user in the users list is always used as the `default` user which means all users - authorized or unauthorized - will inherit the accesses given to this user. The key is arbitrary and you can use any name for it.

```yaml
users:
  - key: this-can-be-any-text
    read:
      - / # Everyone can read all directories
    write:
      - /public # Everyone can write to /public and all its subdirectories
  - key: <HASH>
    # No need to add read here, since its inherited
    write:
      - /private # Only this user can write to /private and all its subdirectories
```

If you don't provide either `read` or `write` to a user, then that user will not have `read` or `write` access respectively (unless inherited from the `default` user). This means that if you want to fully disable the `default` user (so unauthorized users can't access your API), then you can remove both `read` and `write` params from it like so:

```yaml
users:
  - key: disabled # The default user is now disabled, and only authorized users can access the API
  - key: <HASH> # This user can still access the API using its own password
    read:
      - /
    write:
      - /
```

### TLS/SSL

For TLS/SSL configs, setup a reverse proxy using a web server like [Caddy](https://caddyserver.com) or [nginx](https://nginx.org).