{{define "body"}}
{{$card := "w3-card w3-margin w3-padding w3-round"}}
<div class="{{$card}}" id="license">
	<form action="/remove" method="post">
		<p>Enter either the <a href="#control">Control Code</a> or the <a href="#url">URL</a> field.</p>
		<div class="w3-half w3-padding">
			<label for="control">Control Code</label>
			<input class="w3-input w3-round" name="control" id="control" type="password">
		</div>
		<div class="w3-half w3-padding">
			<label for="url">URL</label>
			<input class="w3-input w3-round" name="url" id="url" type="url">
		</div>
		<input class="w3-button w3-border w3-round w3-padding" type="submit" value="Remove">
	</form>
</div>
{{end}}
