package dependencies

type dbms struct {
	DBMS string
	URL  string
}

var DBMSs = []dbms{
	{DBMS: "MySQL", URL: "github.com/go-sql-driver/mysql"},
	{DBMS: "PostgreSQL", URL: "github.com/lib/pq"},
	{DBMS: "SQLite", URL: "github.com/mattn/go-sqlite3"},
	{DBMS: "MongoDB", URL: "go.mongodb.org/mongo-driver/mongo"},
}
