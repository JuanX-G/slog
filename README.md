# slog
## this is version 0.2b, there are a lot of features missing
A simple blog thus a slog, designed to be super minimal and to be run for your friends or a small organization. 
# Requirments:
    - Postgresql running on port 5432 with a DB called slog
    - the tables are set up automatically by the ORM and the existence of them is checked every startup
# Usage
compile and run; the HTTP server will be available on the port set in the SLOG_SERVER_PORT env variable.

# New things since 0.1c
    - unliking
    - bug fixes
    - functions necessary for unliking and couting likes in the UMM-DAL (Ultra Minimal Modern - Data Access Layer) --- COUNT WHERE and DELETE WHERE
    - full end-point docs

# Coming in 0.3
    - Update to the admin tools and stuff

# Coming 0.3b
    - Query by: title, tags and more...

# Later development will focus on SLOC --- Slog-Client

