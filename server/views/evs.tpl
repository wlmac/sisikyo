{{define "body"}}
{{$fmt := "Mon, 02 Jan 2006 15:04:05 MST"}}
{{$ec := .ec}}
{{range $i, $ev := .evs}}
<article class="w3-card w3-round w3-row w3-margin w3-padding vevent" id="event-{{$ev.Id}}">
	<h2><span class="summary">{{$ev.Name}}</span></h2>
	<div class="w3-quarter">
		<h3>Timeframe</h3>
		from
		<time class="dt-start" datetime="{{$ev.Start}}">
			{{$ev.Start.Format $fmt}}
		</time>
		<br/>
		to
		<time class="dt-end" datetime="{{$ev.End}}">
			{{$ev.Start.Format $fmt}}
		</time>
	</div>
	<div class="w3-quarter">
		<h3>Tags</h3>
		<ul>
		{{range $j, $tag := $ev.Tags}}
			<li id="tag-{{$tag.ID}}">
				{{$tag.Name}}
			</li>
		{{end}}
		</ul>
	</div>
	<div class="w3-quarter">
		<h3>Desc</h3>
		<p>
			{{$ev.Desc}}
		</p>
	</div>
	<div class="w3-quarter">
		<h3>Other</h3>
		ID: {{$ev.Id}}
		<br/>
		{{with $_ := $ev.Org.Icon}}
		Org: <a href="{{($ev.Org.URL $ec.API).String}}">
			<img class="w3-circle" height="15rem" src="{{($ev.Org.IconURL $ec.API).String}}"/>
			{{$ev.Org.Name}}
		</a>
		{{else}}
		Org: Unknown
		{{end}}
		<br/>
		Public: {{if $ev.Public}}Yes{{else}}No{{end}}
		<br/>
		Term ID: {{$ev.Term}}
	</div>
</article>
{{end}}
{{end}}
