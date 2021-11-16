<meta charset="UTF-8">
{{with .keywords}}<meta name="keywords" content="{{.|join " "}}">{{end}}
{{with .desc}}<meta name="description" content="{{.}}">{{end}}
<meta name="author" content="Ken Shibata">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.title}}</title>
<link rel="stylesheet" href="/static/w3.css">
<link rel="stylesheet" href="/static/base.css">
