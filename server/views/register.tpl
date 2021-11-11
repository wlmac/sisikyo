{{define "title"}}
Register
{{end}}
{{define "head"}}
<meta http-equiv="refresh" content="0;url={{.url}}" >
{{end}}
{{define "body"}}
{{$card := "w3-card w3-margin w3-padding w3-round"}}
<div class="{{$card}}">
	<h2>Redirecting</h2>
	<p>You will be redirected to the authorization page in 0 second(s).</p>
	<p>If you are not, please go to <a href="{{.url}}">the authorization page</a>.</p>
	<p>Note: you must complete the registration in under 60 seconds (the state cookie is set to expire in 60 seconds)</p>
</div>
<div class="{{$card}}">
	<h2>Instructions</h2>
	<ol>
		<li>press <a href="#authorize">"Authorize"</a></li>
		<li>review the permissions that you will grant to this app</li>
		<li>press "Authorize" on the next page to register</li>
	</ol>
	<a class="w3-button w3-border {{$card}}" id="authorize" href="{{.url}}">Authorize</a>
</div>
{{end}}
