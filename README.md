# slog
## this is version 0.2, there are a lot of features missing
A simple blog thus a slog, designed to be super minimal and to be run for your friends or a small organization. 
# Requirments:
    - Postgresql running on port 5432 with a DB called slog
    - the tables are set up automatically by the ORM and the existence of them is checked every startup
# Usage
compile and run; the HTTP server will be available on the port set in the SLOG_SERVER_PORT env variable.

# New things since 0.1c
    - unliking
    - bug fixes
    - functions necessary for unliking and couting likes in the ORM layer --- COUNT WHERE and DELETE WHERE
# Coming in version 0.2b
    - api docs
    - code docs
    - probaly more docs
    - docs docs docs...

# Coming 0.3
    - Query by: title, tags and more...

# Coming in 0.3b
    - Update to the admin tools and stuff

# Later development will focus on SLOC --- Slog-Client

