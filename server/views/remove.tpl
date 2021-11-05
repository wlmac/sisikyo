{{define "title"}}
About
{{end}}
{{define "body"}}
{{$card := "w3-card w3-margin w3-padding w3-round"}}
<div class="{{$card}}" id="license">
	<form action="/remove" method="post">
		<label for="control">Control Code</label>
		<input name="control" id="control" type="password" required>
		<input type="submit" value="Remove">
	</form>
</div>
{{end}}
