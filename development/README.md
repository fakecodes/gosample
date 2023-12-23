
## Run Locally

Clone the project

```bash
$  git clone https://github.com/fakecodes/gosample.git
```

Go to the project directory

```bash
$  cd gosample
```

Install dependencies, in this case i'll use podman to provision postgresql

```bash
$  podman pod create --name postgre-sql -p 5432:5432
```

```bash
$  podman volume create postgres-data
```

```bash
$  podman run --name db --pod=postgre-sql -d \
    -e POSTGRES_USER=admin \
    -e POSTGRES_PASSWORD=password \
    -v postgres-data:/var/lib/postgresql/data \
    docker.io/library/postgres:14
```

Access the database server using pgadmin4 installed on your local and do manual table creation
```bash
    # CREATE DATABASE task;
```

```bash
    # CREATE TABLE tasks (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        due_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        priority TEXT,
        completed BOOLEAN
    );
```

Start the apps

```bash
$  go run app/main.go
```