package templates

import (
	"net/http"
	"fmt"
)

func ErrorInvalid(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
<html>
<head>
    <title>Error</title>
</head>
<body>
	<center>
	<h2>Session Invalid!</h2>

</body>
</html>`)
}
