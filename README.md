# GATOR
Gator is an RSS Feed Aggregator for the command line written in Go and PostgreSQL.  

## DEPENDENCIES
In order to run gator, you'll need to have Postgres and Go installed on your machine.  
For information on how to install Postgres, visit: [PostgreSQL](https://www.postgresql.org/download/)  
For information on how to install Go, visit: [Go](https://go.dev/doc/install)  
  
You will need to create a `.gatorconfig.json` file in your home directory. The format of the `.gatorconfig.json` should be as follows:  

    ```json  
    {
        "db_url":"postgres://username:password@localhost:port/gator?sslmode=disable",
        "current_user_name":""
     }
     ```  
Where username and password are the values you assigned when setting up PostgreSQL. For more information on setting up PostgreSQL, visit: [Microsoft](https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql)  

## BUILDING AND INSTALLING GATOR
Assuming that Git and Go are installed; cloning, running, building, and installing the code can be done using the following commands:  

    ```bash
    $ git clone https://github.com/jubilant-gremlin/gator.git
    $ cd gator
    $ go build
    $ go install
    ```  

The config file expects that you have a `postgres` instance running on the provided port with a `gator` database created. If issues are encountered while installing gator, please ensure all dependencies are satisified and that the `gator` database exists.  

## USING GATOR
Once gator is installed on your machine, you can run the application with `gator` combined with one of the following commands:
- `login <username>`: sets the current user to the provided username if the user is registered.  
- `register <username>`: registers a new user in the database with the provided username, and logs in as that user.  
- `reset`: resets all of the data in the database  
- `users`: lists all users registered in the database  
- `agg <timeBetweenRequests>`: aggregates posts from the users followed feeds with the given time between requests. Feeds never fetched will be aggregated first, followed by oldest feeds first. Time must be in the format `<integer>s|m|h` e.g. 15s, 10m, 1h  
- `addfeed <name> <url>`: adds a given feed as the logged in user, then follows that feed as the logged in user.  
- `follow <url>`: follows the given feed as the logged in user.
- `following`: displays the logged in user's currently followed feeds.
- `unfollow <url>`: unfollows the given feed as the logged in user.
- `browse [limit]`: displays the logged in user's posts, showing the most recently published first. Takes an optional limit argument, which changes how many posts are displayed. If not provided, the default limit is 2.






