package dependencies

type framework struct {
	Name string
	URL  string
}

var Frameworks = []framework{
	{Name: "Gin", URL: "github.com/gin-gonic/gin"},
	{Name: "Echo", URL: "github.com/labstack/echo"},
	{Name: "Iris", URL: "github.com/kataras/iris"},
}
