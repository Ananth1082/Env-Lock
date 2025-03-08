# Elk: ENV Lock

Is a simple CLI tool to manage all your env files.

## Working

Elk stores all your env file in ecrypted format which can only be accessed with the password. So you can store all your env files safely using it.

## CLI Commands

Elk has 5 commnands create, update, list, get, delete.

### Create

Creates a lock for a given env file.

```
elk create -f <env_filename> -n <name_of_lock> -d <description>

```

### Update

Updates the value in the lock

```
elk update -id <lock_id> -n <optional: new_name> -d <optional: new description> -f <optional: new_env_file>

```

### Get

Fetches the env file in the lock

```
elk get -id <lock_id>
```

### List

Lists all the locks present

```
elk list
```

### Delete

Delete the given lock

```
elk delete -id <lock_id>
```
