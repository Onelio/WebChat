package templates

import (
	"net/http"
	"fmt"
)

func Logout(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, `
<html>
<head>
    <title>WebChat Logout Page</title>
</head>
<body>
	<center>
	<h2>Successfully logged out!</h2>

</body>
</html>`)

}