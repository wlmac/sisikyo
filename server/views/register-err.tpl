{{define "body"}}
{{$card := "w3-card w3-margin w3-padding w3-round"}}
<div class="{{$card}} w3-half">
	<p>{{.desc}}</p>
	<code>{{.err}}</code>
</div>
{{template "links" (dict "control" .control "class" "w3-half")}}
{{end}}
