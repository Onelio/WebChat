package templates

import (
	"net/http"
	"fmt"
)

func ErrorTaken(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
<html>
<head>
    <title>Error</title>
</head>
<body>
	<center>
	<h2>Username exist!</h2>

</body>
</html>`)
}

