- **Missing Database Password**

    A database user should always have a password.

    *Added password entry to the DSN.*

- **Missing Error Check in `main`**

    Unhandled errors can lead to problems, especially if they are recoverable.

    *Added error checks to fail when the DB connection is unavailable (retries could be added later).*

- **Missing HTTP Verbs in `HandleFunc`**

    Without HTTP verbs, the type of operation is unclear.

    *Added corresponding verbs for each route.*

- **Fixed `createUser` URL**

    It makes more sense to `POST` to `/users` since `createUser` adds a user to the collection.

    *Renamed `/create` to `/users`.*

- **Query Users by ID**

    Querying by ID is more precise, as names are not unique.

    *Modified the query to select by ID.*

- **Missing Error Checks in `getUsers`**

    Similar to the error handling in `main`, errors need to be checked.

    *Added error handling before `rows.Close()` to prevent nil pointer exceptions; added error checks for `Scan()`.*

- **Extended Response for `getUsers`**

    User IDs should be returned alongside names for clarity.

    *Updated the response to return users as JSON, including IDs.*

- **Removed Wait Group**

    Since a DBMS like PostgreSQL is used, concurrency control is handled by the database.

    *Removed the wait group.*

- **Removed Unnecessary Goroutines**

    As `getUsers` and `createUser` are triggered via HTTP handlers, which already run in their own goroutines, explicit goroutines are redundant.

    *Removed the explicit goroutines.*

- **Query Param → Payload**

    For `POST` requests, sending data in the payload is more appropriate than using query parameters.

    *Extracted `username` from the request body.*

- **Validate Before Long Database Operations**

    Validation is crucial to avoid unnecessary processing.

    *Added input validation.*

- **Prevent SQL Injection**

    Never allow user input to be directly included in SQL queries.

    *Passed `username` as an argument.*

- **Return HTTP Error on Failure**

    Returning `http.StatusOK` on failure is misleading.

    *Return proper HTTP error codes using `http.Error()`.*

- **Return ID for Created Entities**

    Returning the ID of newly created entities helps the client avoid extra requests.

    *Used `db.QueryRow()` as per the pq documentation to return the ID.*
