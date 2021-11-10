<!DOCTYPE html>
<html>
	<head>
	        {{include "layouts/head"}}
		{{block "head" .}}{{end}}
	</head>
	<body>
        	{{include "layouts/header"}}
		{{block "body" .}}{{end}}
        	{{include "layouts/footer"}}
	</body>
</html>
